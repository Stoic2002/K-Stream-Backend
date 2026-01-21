package episode

import (
	"context"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	FindBySeasonID(ctx context.Context, seasonID string) ([]Episode, error)
	FindByID(ctx context.Context, id string) (*Episode, error)
	Create(ctx context.Context, episode *Episode) error
	Update(ctx context.Context, episode *Episode) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindBySeasonID(ctx context.Context, seasonID string) ([]Episode, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `
		SELECT id, season_id, episode_number, title, video_url, duration, thumbnail_url, view_count, source_url, created_at 
		FROM episodes 
		WHERE season_id = $1 
		ORDER BY episode_number ASC
	`
	rows, err := db.Query(ctx, query, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []Episode
	for rows.Next() {
		var e Episode
		var thumbnail, source *string
		if err := rows.Scan(&e.ID, &e.SeasonID, &e.EpisodeNumber, &e.Title, &e.VideoURL, &e.Duration, &thumbnail, &e.ViewCount, &source, &e.CreatedAt); err != nil {
			return nil, err
		}
		if thumbnail != nil {
			e.ThumbnailURL = *thumbnail
		}
		if source != nil {
			e.SourceURL = *source
		}
		episodes = append(episodes, e)
	}

	return episodes, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Episode, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `
		SELECT id, season_id, episode_number, title, video_url, duration, thumbnail_url, view_count, source_url, added_by, created_at
		FROM episodes WHERE id = $1
	`
	var e Episode
	var thumbnail, source, addedBy *string
	err := db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.SeasonID, &e.EpisodeNumber, &e.Title, &e.VideoURL, &e.Duration,
		&thumbnail, &e.ViewCount, &source, &addedBy, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if thumbnail != nil {
		e.ThumbnailURL = *thumbnail
	}
	if source != nil {
		e.SourceURL = *source
	}
	if addedBy != nil {
		e.AddedBy = *addedBy
	}

	return &e, nil
}

func (r *repository) Create(ctx context.Context, episode *Episode) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `
		INSERT INTO episodes (season_id, episode_number, title, video_url, duration, thumbnail_url, source_url, added_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return db.QueryRow(ctx, query,
		episode.SeasonID, episode.EpisodeNumber, episode.Title, episode.VideoURL,
		episode.Duration, episode.ThumbnailURL, episode.SourceURL, episode.AddedBy, time.Now(),
	).Scan(&episode.ID)
}

func (r *repository) Update(ctx context.Context, episode *Episode) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `
		UPDATE episodes 
		SET episode_number=$1, title=$2, video_url=$3, duration=$4, thumbnail_url=$5, source_url=$6
		WHERE id=$7
	`
	_, err := db.Exec(ctx, query,
		episode.EpisodeNumber, episode.Title, episode.VideoURL,
		episode.Duration, episode.ThumbnailURL, episode.SourceURL, episode.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	_, err := db.Exec(ctx, "DELETE FROM episodes WHERE id = $1", id)
	return err
}
