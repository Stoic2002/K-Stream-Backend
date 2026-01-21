package comment

import (
	"context"
	"drakor-backend/internal/auth"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByEpisodeID(ctx context.Context, episodeID string, limit, offset int) ([]Comment, int64, error)
	GetByID(ctx context.Context, id string) (*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, comment *Comment) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `
		INSERT INTO comments (user_id, episode_id, comment_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return db.QueryRow(ctx, query,
		comment.UserID, comment.EpisodeID, comment.CommentText, time.Now(), time.Now(),
	).Scan(&comment.ID)
}

func (r *repository) GetByEpisodeID(ctx context.Context, episodeID string, limit, offset int) ([]Comment, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// Count
	var total int64
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM comments WHERE episode_id = $1", episodeID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT c.id, c.user_id, c.episode_id, c.comment_text, c.created_at, c.updated_at,
		       u.id, u.name, u.avatar_url
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.episode_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, episodeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var u auth.User
		var avatar *string
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.EpisodeID, &c.CommentText, &c.CreatedAt, &c.UpdatedAt,
			&u.ID, &u.Name, &avatar,
		); err != nil {
			return nil, 0, err
		}
		if avatar != nil {
			u.AvatarURL = *avatar
		}
		c.User = u
		comments = append(comments, c)
	}

	return comments, total, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*Comment, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, user_id, episode_id, comment_text, created_at, updated_at FROM comments WHERE id = $1`
	var c Comment
	err := db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.UserID, &c.EpisodeID, &c.CommentText, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Update(ctx context.Context, comment *Comment) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	query := `UPDATE comments SET comment_text = $1, updated_at = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, comment.CommentText, time.Now(), comment.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	_, err := db.Exec(ctx, "DELETE FROM comments WHERE id = $1", id)
	return err
}
