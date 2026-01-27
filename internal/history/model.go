package history

import (
	"drakor-backend/internal/episode"
	"time"
)

type WatchHistory struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	EpisodeID       string          `json:"episode_id"`
	ProgressSeconds int             `json:"progress_seconds"`
	IsFinished      bool            `json:"completed"`
	LastWatchedAt   time.Time       `json:"last_watched_at"`
	Episode         episode.Episode `json:"episode,omitempty"`
}

type RecordHistoryRequest struct {
	EpisodeID       string `json:"episode_id" validate:"required,uuid"`
	ProgressSeconds int    `json:"progress_seconds" validate:"min=0"`
	IsFinished      bool   `json:"completed"`
}
