package episode

import "time"

type Episode struct {
	ID            string    `json:"id"`
	SeasonID      string    `json:"season_id"`
	EpisodeNumber int       `json:"episode_number"`
	Title         string    `json:"title"`
	VideoURL      string    `json:"video_url"`
	Duration      int       `json:"duration"` // in seconds
	ThumbnailURL  string    `json:"thumbnail_url"`
	ViewCount     int       `json:"view_count"`
	SourceURL     string    `json:"source_url"`
	AddedBy       string    `json:"added_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateEpisodeRequest struct {
	SeasonID      string `json:"season_id" validate:"required,uuid"`
	EpisodeNumber int    `json:"episode_number" validate:"required,min=1"`
	Title         string `json:"title" validate:"required"`
	VideoURL      string `json:"video_url" validate:"required,url"`
	Duration      int    `json:"duration" validate:"min=1"`
	ThumbnailURL  string `json:"thumbnail_url" validate:"omitempty,url"`
	SourceURL     string `json:"source_url" validate:"omitempty,url"`
}

type UpdateEpisodeRequest struct {
	EpisodeNumber int    `json:"episode_number" validate:"required,min=1"`
	Title         string `json:"title" validate:"required"`
	VideoURL      string `json:"video_url" validate:"required,url"`
	Duration      int    `json:"duration" validate:"min=1"`
	ThumbnailURL  string `json:"thumbnail_url" validate:"omitempty,url"`
	SourceURL     string `json:"source_url" validate:"omitempty,url"`
}
