package drama

import (
	"context"
	"drakor-backend/internal/genre"
	"drakor-backend/pkg/database"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	FindAll(ctx context.Context, page, limit int, query, genreID, status, sort string) ([]Drama, int64, error)
	FindByID(ctx context.Context, id string) (*Drama, error)
	Create(ctx context.Context, drama *Drama, genreIDs []string, actors []DramaActorReq, startTime time.Time) error
	Update(ctx context.Context, drama *Drama, genreIDs []string, actors []DramaActorReq) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindAll(ctx context.Context, page, limit int, queryStr, genreID, status, sort string) ([]Drama, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// Base query
	sql := `SELECT id, title, poster_url, year, rating, status, view_count, created_at FROM dramas WHERE 1=1`
	countSql := `SELECT COUNT(*) FROM dramas WHERE 1=1`
	args := []interface{}{}
	argId := 1

	// Dynamic filters
	if queryStr != "" {
		filter := fmt.Sprintf(" AND (title ILIKE $%d)", argId)
		sql += filter
		countSql += filter
		args = append(args, "%"+queryStr+"%")
		argId++
	}

	if genreID != "" {
		// Subquery to check if drama has this genre
		filter := fmt.Sprintf(" AND id IN (SELECT drama_id FROM drama_genres WHERE genre_id = $%d)", argId)
		sql += filter
		countSql += filter
		args = append(args, genreID)
		argId++
	}

	if status != "" {
		filter := fmt.Sprintf(" AND status = $%d", argId)
		sql += filter
		countSql += filter
		args = append(args, status)
		argId++
	}

	// Counting total
	var total int64
	err := db.QueryRow(ctx, countSql, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Sorting
	switch sort {
	case "popular":
		sql += " ORDER BY view_count DESC"
	case "rating":
		sql += " ORDER BY rating DESC"
	case "oldest":
		sql += " ORDER BY created_at ASC"
	default: // "latest"
		sql += " ORDER BY created_at DESC"
	}

	// Pagination
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var dramas []Drama
	for rows.Next() {
		var d Drama
		var poster *string
		if err := rows.Scan(&d.ID, &d.Title, &poster, &d.Year, &d.Rating, &d.Status, &d.ViewCount, &d.CreatedAt); err != nil {
			return nil, 0, err
		}
		if poster != nil {
			d.PosterURL = *poster
		}
		dramas = append(dramas, d)
	}

	return dramas, total, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Drama, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	// 1. Fetch Drama Details
	query := `
		SELECT id, title, synopsis, poster_url, year, rating, total_seasons, status, view_count, source_url, added_by, created_at, updated_at
		FROM dramas WHERE id = $1
	`
	var d Drama
	var synopsis, poster, source, addedBy *string

	err := db.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.Title, &synopsis, &poster, &d.Year, &d.Rating, &d.TotalSeasons,
		&d.Status, &d.ViewCount, &source, &addedBy, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if synopsis != nil {
		d.Synopsis = *synopsis
	}
	if poster != nil {
		d.PosterURL = *poster
	}
	if source != nil {
		d.SourceURL = *source
	}
	if addedBy != nil {
		d.AddedBy = *addedBy
	}

	// 2. Fetch Genres
	genreQuery := `
		SELECT g.id, g.name, g.slug 
		FROM genres g
		JOIN drama_genres dg ON g.id = dg.genre_id
		WHERE dg.drama_id = $1
	`
	gRows, err := db.Query(ctx, genreQuery, id)
	if err == nil {
		defer gRows.Close()
		for gRows.Next() {
			var g genre.Genre
			if err := gRows.Scan(&g.ID, &g.Name, &g.Slug); err == nil {
				d.Genres = append(d.Genres, g)
			}
		}
	}

	// 3. Fetch Actors
	actorQuery := `
		SELECT a.id, a.name, a.photo_url, da.role
		FROM actors a
		JOIN drama_actors da ON a.id = da.actor_id
		WHERE da.drama_id = $1
	`
	aRows, err := db.Query(ctx, actorQuery, id)
	if err == nil {
		defer aRows.Close()
		for aRows.Next() {
			var da DramaActor
			var photo *string
			if err := aRows.Scan(&da.Actor.ID, &da.Actor.Name, &photo, &da.Role); err == nil {
				if photo != nil {
					da.Actor.PhotoURL = *photo
				}
				d.Actors = append(d.Actors, da)
			}
		}
	}

	return &d, nil
}

func (r *repository) Create(ctx context.Context, drama *Drama, genreIDs []string, actors []DramaActorReq, startTime time.Time) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Drama
	query := `
		INSERT INTO dramas (title, synopsis, poster_url, year, total_seasons, status, source_url, added_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		drama.Title, drama.Synopsis, drama.PosterURL, drama.Year, drama.TotalSeasons,
		drama.Status, drama.SourceURL, drama.AddedBy, startTime, startTime,
	).Scan(&drama.ID)
	if err != nil {
		return err
	}

	// 2. Insert Genres
	if len(genreIDs) > 0 {
		for _, gid := range genreIDs {
			_, err := tx.Exec(ctx, "INSERT INTO drama_genres (drama_id, genre_id) VALUES ($1, $2)", drama.ID, gid)
			if err != nil {
				// Continue or Fail? Standard is Fail.
				// Likely foreign key violation if genre doesn't exist
				return fmt.Errorf("failed to add genre %s: %v", gid, err)
			}
		}
	}

	// 3. Insert Actors
	if len(actors) > 0 {
		for _, act := range actors {
			_, err := tx.Exec(ctx, "INSERT INTO drama_actors (drama_id, actor_id, role) VALUES ($1, $2, $3)", drama.ID, act.ActorID, act.Role)
			if err != nil {
				return fmt.Errorf("failed to add actor %s: %v", act.ActorID, err)
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) Update(ctx context.Context, drama *Drama, genreIDs []string, actors []DramaActorReq) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Update Drama Fields
	query := `
		UPDATE dramas
		SET title=$1, synopsis=$2, poster_url=$3, year=$4, total_seasons=$5, status=$6, source_url=$7, updated_at=$8
		WHERE id=$9
	`
	_, err = tx.Exec(ctx, query,
		drama.Title, drama.Synopsis, drama.PosterURL, drama.Year, drama.TotalSeasons,
		drama.Status, drama.SourceURL, time.Now(), drama.ID,
	)
	if err != nil {
		return err
	}

	// 2. Update Genres (Strategy: Delete All & Re-insert) helps avoid diff logic complexity
	_, err = tx.Exec(ctx, "DELETE FROM drama_genres WHERE drama_id = $1", drama.ID)
	if err != nil {
		return err
	}
	if len(genreIDs) > 0 {
		for _, gid := range genreIDs {
			_, err := tx.Exec(ctx, "INSERT INTO drama_genres (drama_id, genre_id) VALUES ($1, $2)", drama.ID, gid)
			if err != nil {
				return err
			}
		}
	}

	// 3. Update Actors (Strategy: Delete All & Re-insert)
	_, err = tx.Exec(ctx, "DELETE FROM drama_actors WHERE drama_id = $1", drama.ID)
	if err != nil {
		return err
	}
	if len(actors) > 0 {
		for _, act := range actors {
			_, err := tx.Exec(ctx, "INSERT INTO drama_actors (drama_id, actor_id, role) VALUES ($1, $2, $3)", drama.ID, act.ActorID, act.Role)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	// Cascade delete handles relations
	_, err := db.Exec(ctx, "DELETE FROM dramas WHERE id = $1", id)
	return err
}
