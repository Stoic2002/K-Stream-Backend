package drama

import (
	"drakor-backend/internal/actor"
	"drakor-backend/internal/genre"
	"time"
)

type Drama struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Synopsis     string        `json:"synopsis"`
	PosterURL    string        `json:"poster_url"`
	Year         int           `json:"year"`
	Rating       float64       `json:"rating"`
	TotalSeasons int           `json:"total_seasons"`
	Status       string        `json:"status"` // 'ongoing', 'completed'
	ViewCount    int           `json:"view_count"`
	SourceURL    string        `json:"source_url"` // Trailer or internal source
	AddedBy      string        `json:"added_by,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Genres       []genre.Genre `json:"genres,omitempty"`
	Actors       []DramaActor  `json:"actors,omitempty"`
}

type DramaActor struct {
	Actor actor.Actor `json:"actor"`
	Role  string      `json:"role"` // 'main', 'support'
}

type CreateDramaRequest struct {
	Title        string          `json:"title" validate:"required,min=2,max=255"`
	Synopsis     string          `json:"synopsis"`
	PosterURL    string          `json:"poster_url" validate:"omitempty,url"`
	Year         int             `json:"year" validate:"required,min=1900,max=2100"`
	TotalSeasons int             `json:"total_seasons" validate:"min=1"`
	Status       string          `json:"status" validate:"required,oneof=ongoing completed"`
	SourceURL    string          `json:"source_url" validate:"omitempty,url"`
	GenreIDs     []string        `json:"genre_ids" validate:"required,min=1"`
	Actors       []DramaActorReq `json:"actors" validate:"omitempty,dive"`
}

type DramaActorReq struct {
	ActorID string `json:"actor_id" validate:"required,uuid"`
	Role    string `json:"role" validate:"required,oneof=main support"`
}

type UpdateDramaRequest struct {
	Title        string          `json:"title" validate:"required,min=2,max=255"`
	Synopsis     string          `json:"synopsis"`
	PosterURL    string          `json:"poster_url" validate:"omitempty,url"`
	Year         int             `json:"year" validate:"required,min=1900,max=2100"`
	TotalSeasons int             `json:"total_seasons" validate:"min=1"`
	Status       string          `json:"status" validate:"required,oneof=ongoing completed"`
	SourceURL    string          `json:"source_url" validate:"omitempty,url"`
	GenreIDs     []string        `json:"genre_ids" validate:"required,min=1"`
	Actors       []DramaActorReq `json:"actors" validate:"omitempty,dive"`
}
