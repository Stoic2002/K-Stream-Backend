package genre

import (
	"context"
	"drakor-backend/pkg/validator"
	"errors"
)

type Service interface {
	GetAll(ctx context.Context) ([]Genre, error)
	Create(ctx context.Context, req CreateGenreRequest) (*Genre, error)
	Update(ctx context.Context, id string, req UpdateGenreRequest) (*Genre, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]Genre, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Create(ctx context.Context, req CreateGenreRequest) (*Genre, error) {
	// Generate slug if empty
	if req.Slug == "" {
		req.Slug = validator.GenerateSlug(req.Name)
	} else {
		req.Slug = validator.GenerateSlug(req.Slug)
	}

	// Check if slug exists
	existing, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("slug already exists")
	}

	genre := &Genre{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.repo.Create(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateGenreRequest) (*Genre, error) {
	genre, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if genre == nil {
		return nil, errors.New("genre not found")
	}

	// Update fields
	genre.Name = req.Name
	newSlug := validator.GenerateSlug(req.Slug)

	// Check if new slug conflicts with OTHER genre
	if newSlug != genre.Slug {
		existing, err := s.repo.FindBySlug(ctx, newSlug)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("slug already exists")
		}
		genre.Slug = newSlug
	}

	if err := s.repo.Update(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	// Check/validasi referensi drama_genres nanti bisa ditambahkan disini
	return s.repo.Delete(ctx, id)
}
