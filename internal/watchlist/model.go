package watchlist

import (
	"drakor-backend/internal/drama"
	"time"
)

type WatchlistItem struct {
	UserID    string      `json:"user_id"`
	DramaID   string      `json:"drama_id"`
	CreatedAt time.Time   `json:"created_at"`
	Drama     drama.Drama `json:"drama,omitempty"`
}

type AddWatchlistRequest struct {
	DramaID string `json:"drama_id" validate:"required,uuid"`
}
