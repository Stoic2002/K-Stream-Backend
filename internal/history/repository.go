package history

import (
	"context"
	"drakor-backend/internal/episode"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Upsert(ctx context.Context, userID, episodeID string, progress int, isFinished bool) error
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]WatchHistory, int64, error)
	GetByEpisode(ctx context.Context, userID, episodeID string) (*WatchHistory, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Upsert(ctx context.Context, userID, episodeID string, progress int, isFinished bool) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	// Logic: Try update, if no rows updated, then insert.
	// Or use ON CONFLICT (user_id, episode_id) DO UPDATE
	query := `
		INSERT INTO watch_history (user_id, episode_id, progress_seconds, completed, last_watched_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, episode_id) 
		DO UPDATE SET 
			progress_seconds = EXCLUDED.progress_seconds,
			completed = EXCLUDED.completed,
			last_watched_at = EXCLUDED.last_watched_at
	`
	_, err := db.Exec(ctx, query, userID, episodeID, progress, isFinished, time.Now())
	return err
}

func (r *repository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]WatchHistory, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	var total int64
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM watch_history WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch with episode detail
	query := `
		SELECT wh.id, wh.user_id, wh.episode_id, wh.progress_seconds, wh.completed, wh.last_watched_at,
		       e.id, e.season_id, e.episode_number, e.title, e.thumbnail_url, e.duration
		FROM watch_history wh
		JOIN episodes e ON wh.episode_id = e.id
		WHERE wh.user_id = $1
		ORDER BY wh.last_watched_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var histories []WatchHistory
	for rows.Next() {
		var h WatchHistory
		var e episode.Episode
		var thumbnail *string

		if err := rows.Scan(
			&h.ID, &h.UserID, &h.EpisodeID, &h.ProgressSeconds, &h.IsFinished, &h.LastWatchedAt,
			&e.ID, &e.SeasonID, &e.EpisodeNumber, &e.Title, &thumbnail, &e.Duration,
		); err != nil {
			return nil, 0, err
		}
		if thumbnail != nil {
			e.ThumbnailURL = *thumbnail
		}
		h.Episode = e
		histories = append(histories, h)
	}

	return histories, total, nil
}

func (r *repository) GetByEpisode(ctx context.Context, userID, episodeID string) (*WatchHistory, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `
		SELECT id, user_id, episode_id, progress_seconds, completed, last_watched_at
		FROM watch_history
		WHERE user_id = $1 AND episode_id = $2
	`
	var h WatchHistory
	err := db.QueryRow(ctx, query, userID, episodeID).Scan(
		&h.ID, &h.UserID, &h.EpisodeID, &h.ProgressSeconds, &h.IsFinished, &h.LastWatchedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is strictly nil, not error
		}
		return nil, err
	}
	return &h, nil
}
