package comment

import (
	"drakor-backend/internal/auth"
	"time"
)

type Comment struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	EpisodeID   string    `json:"episode_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        auth.User `json:"user,omitempty"`
}

type CreateCommentRequest struct {
	EpisodeID   string `json:"episode_id" validate:"required,uuid"`
	CommentText string `json:"comment_text" validate:"required,max=1000"`
}

type UpdateCommentRequest struct {
	CommentText string `json:"comment_text" validate:"required,max=1000"`
}
