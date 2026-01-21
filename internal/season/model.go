package season

import "time"

type Season struct {
	ID           string    `json:"id"`
	DramaID      string    `json:"drama_id"`
	SeasonNumber int       `json:"season_number"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateSeasonRequest struct {
	DramaID      string `json:"drama_id" validate:"required,uuid"`
	SeasonNumber int    `json:"season_number" validate:"required,min=1"`
	Title        string `json:"title" validate:"required"`
}

type UpdateSeasonRequest struct {
	SeasonNumber int    `json:"season_number" validate:"required,min=1"`
	Title        string `json:"title" validate:"required"`
}
