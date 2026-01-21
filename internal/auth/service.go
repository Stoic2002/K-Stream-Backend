package auth

import (
	"context"
	"drakor-backend/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*User, error)
	// Admin
	GetAllUsers(ctx context.Context, page, limit int) ([]User, int64, error)
	UpdateUserRole(ctx context.Context, userID, role string) error
	DeleteUser(ctx context.Context, userID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if email exists
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // Default role
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find user
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userID string) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.AvatarURL = req.AvatarURL

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetAllUsers(ctx context.Context, page, limit int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) UpdateUserRole(ctx context.Context, userID, role string) error {
	if role != "admin" && role != "user" {
		return errors.New("invalid role")
	}
	return s.repo.UpdateRole(ctx, userID, role)
}

func (s *service) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}
