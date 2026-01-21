package review

import (
	"context"
	"drakor-backend/internal/auth"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, review *Review) error
	GetByDramaID(ctx context.Context, dramaID string, limit, offset int) ([]Review, int64, error)
	GetByUserAndDrama(ctx context.Context, userID, dramaID string) (*Review, error)
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Review, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, review *Review) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `
		INSERT INTO reviews (user_id, drama_id, rating, review_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return db.QueryRow(ctx, query,
		review.UserID, review.DramaID, review.Rating, review.ReviewText,
		time.Now(), time.Now(),
	).Scan(&review.ID)
}

func (r *repository) GetByDramaID(ctx context.Context, dramaID string, limit, offset int) ([]Review, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// Count total
	var total int64
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM reviews WHERE drama_id = $1", dramaID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT r.id, r.user_id, r.drama_id, r.rating, r.review_text, r.created_at, r.updated_at,
		       u.id, u.name, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.drama_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, dramaID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var rev Review
		var u auth.User
		var avatar *string
		if err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.DramaID, &rev.Rating, &rev.ReviewText, &rev.CreatedAt, &rev.UpdatedAt,
			&u.ID, &u.Name, &avatar,
		); err != nil {
			return nil, 0, err
		}
		if avatar != nil {
			u.AvatarURL = *avatar
		}
		rev.User = u
		reviews = append(reviews, rev)
	}

	return reviews, total, nil
}

func (r *repository) GetByUserAndDrama(ctx context.Context, userID, dramaID string) (*Review, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `
		SELECT id, user_id, drama_id, rating, review_text, created_at, updated_at
		FROM reviews
		WHERE user_id = $1 AND drama_id = $2
	`
	var rev Review
	err := db.QueryRow(ctx, query, userID, dramaID).Scan(
		&rev.ID, &rev.UserID, &rev.DramaID, &rev.Rating, &rev.ReviewText, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rev, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Review, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, user_id, drama_id, rating, review_text, created_at, updated_at FROM reviews WHERE id = $1`
	var rev Review
	err := db.QueryRow(ctx, query, id).Scan(
		&rev.ID, &rev.UserID, &rev.DramaID, &rev.Rating, &rev.ReviewText, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rev, nil
}

func (r *repository) Update(ctx context.Context, review *Review) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `UPDATE reviews SET rating = $1, review_text = $2, updated_at = $3 WHERE id = $4`
	_, err := db.Exec(ctx, query, review.Rating, review.ReviewText, time.Now(), review.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	_, err := db.Exec(ctx, "DELETE FROM reviews WHERE id = $1", id)
	return err
}
