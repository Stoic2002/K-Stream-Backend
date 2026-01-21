package watchlist

import (
	"context"
	"drakor-backend/internal/drama"
	"drakor-backend/pkg/database"
	"errors"
	"time"
)

type Repository interface {
	Add(ctx context.Context, userID, dramaID string) error
	Remove(ctx context.Context, userID, dramaID string) error
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]WatchlistItem, int64, error)
	Exists(ctx context.Context, userID, dramaID string) (bool, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Add(ctx context.Context, userID, dramaID string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `INSERT INTO watchlist (user_id, drama_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := db.Exec(ctx, query, userID, dramaID, time.Now())
	return err
}

func (r *repository) Remove(ctx context.Context, userID, dramaID string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `DELETE FROM watchlist WHERE user_id = $1 AND drama_id = $2`
	_, err := db.Exec(ctx, query, userID, dramaID)
	return err
}

func (r *repository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]WatchlistItem, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// Count total
	var total int64
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM watchlist WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch items with drama details
	query := `
		SELECT w.user_id, w.drama_id, w.created_at,
		       d.id, d.title, d.poster_url, d.year, d.rating, d.status
		FROM watchlist w
		JOIN dramas d ON w.drama_id = d.id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []WatchlistItem
	for rows.Next() {
		var w WatchlistItem
		var d drama.Drama
		var poster *string

		if err := rows.Scan(
			&w.UserID, &w.DramaID, &w.CreatedAt,
			&d.ID, &d.Title, &poster, &d.Year, &d.Rating, &d.Status,
		); err != nil {
			return nil, 0, err
		}
		if poster != nil {
			d.PosterURL = *poster
		}
		w.Drama = d
		items = append(items, w)
	}

	return items, total, nil
}

func (r *repository) Exists(ctx context.Context, userID, dramaID string) (bool, error) {
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database not connected")
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM watchlist WHERE user_id = $1 AND drama_id = $2)`
	err := db.QueryRow(ctx, query, userID, dramaID).Scan(&exists)
	return exists, err
}
