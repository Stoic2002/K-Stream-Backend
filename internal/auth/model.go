package auth

import (
	"time"
)

// User represents the user model
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never result password hash
	Name         string    `json:"name"`
	AvatarURL    string    `json:"avatar_url"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest is the payload for registration
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the payload for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest is the payload for updating profile
type UpdateProfileRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

// AuthResponse is the response payload giving token
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
