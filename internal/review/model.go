package review

import (
	"drakor-backend/internal/auth"
	"time"
)

type Review struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	DramaID    string    `json:"drama_id"`
	Rating     int       `json:"rating"`
	ReviewText string    `json:"review_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       auth.User `json:"user,omitempty"`
}

type CreateReviewRequest struct {
	DramaID    string `json:"drama_id" validate:"required,uuid"`
	Rating     int    `json:"rating" validate:"required,min=1,max=10"`
	ReviewText string `json:"review_text" validate:"required,max=1000"`
}

type UpdateReviewRequest struct {
	Rating     int    `json:"rating" validate:"required,min=1,max=10"`
	ReviewText string `json:"review_text" validate:"required,max=1000"`
}
