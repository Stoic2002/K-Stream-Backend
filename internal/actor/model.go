package actor

import (
	"time"
)

type Actor struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateActorRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}

type UpdateActorRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}
