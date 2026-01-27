package actor

import (
	"context"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	FindAll(ctx context.Context, limit, offset int, search string) ([]Actor, int64, error)
	FindByID(ctx context.Context, id string) (*Actor, error)
	Create(ctx context.Context, actor *Actor) error
	Update(ctx context.Context, actor *Actor) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindAll(ctx context.Context, limit, offset int, search string) ([]Actor, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// Count total
	var total int64
	var countQuery string
	var args []interface{}

	if search != "" {
		countQuery = "SELECT COUNT(*) FROM actors WHERE name ILIKE $1"
		args = append(args, "%"+search+"%")
	} else {
		countQuery = "SELECT COUNT(*) FROM actors"
	}

	err := db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	var query string
	var queryArgs []interface{}

	if search != "" {
		query = `SELECT id, name, photo_url, created_at FROM actors WHERE name ILIKE $1 ORDER BY name ASC LIMIT $2 OFFSET $3`
		queryArgs = append(queryArgs, "%"+search+"%", limit, offset)
	} else {
		query = `SELECT id, name, photo_url, created_at FROM actors ORDER BY name ASC LIMIT $1 OFFSET $2`
		queryArgs = append(queryArgs, limit, offset)
	}

	rows, err := db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var actors []Actor
	for rows.Next() {
		var a Actor
		// photo_url may be null in DB, handle appropriately if needed but for now assuming string match
		var photoURL *string
		if err := rows.Scan(&a.ID, &a.Name, &photoURL, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		if photoURL != nil {
			a.PhotoURL = *photoURL
		}
		actors = append(actors, a)
	}

	return actors, total, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Actor, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, name, photo_url, created_at FROM actors WHERE id = $1`
	var a Actor
	var photoURL *string
	err := db.QueryRow(ctx, query, id).Scan(&a.ID, &a.Name, &photoURL, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if photoURL != nil {
		a.PhotoURL = *photoURL
	}
	return &a, nil
}

func (r *repository) Create(ctx context.Context, actor *Actor) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `INSERT INTO actors (name, photo_url, created_at) VALUES ($1, $2, $3) RETURNING id`
	return db.QueryRow(ctx, query, actor.Name, actor.PhotoURL, time.Now()).Scan(&actor.ID)
}

func (r *repository) Update(ctx context.Context, actor *Actor) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `UPDATE actors SET name = $1, photo_url = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, actor.Name, actor.PhotoURL, actor.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `DELETE FROM actors WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}
