package drama

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	GetAll(ctx context.Context, page, limit int, query, genreID, status, sort string) ([]Drama, int64, error)
	GetByID(ctx context.Context, id string) (*Drama, error)
	Create(ctx context.Context, userID string, req CreateDramaRequest) (*Drama, error)
	Update(ctx context.Context, id string, req UpdateDramaRequest) (*Drama, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context, page, limit int, query, genreID, status, sort string) ([]Drama, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.repo.FindAll(ctx, page, limit, query, genreID, status, sort)
}

func (s *service) GetByID(ctx context.Context, id string) (*Drama, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, userID string, req CreateDramaRequest) (*Drama, error) {
	drama := &Drama{
		Title:        req.Title,
		Synopsis:     req.Synopsis,
		PosterURL:    req.PosterURL,
		Year:         req.Year,
		TotalSeasons: req.TotalSeasons,
		Status:       req.Status,
		SourceURL:    req.SourceURL,
		AddedBy:      userID,
	}

	if err := s.repo.Create(ctx, drama, req.GenreIDs, req.Actors, time.Now()); err != nil {
		return nil, err
	}

	// Refetch full object to return complete structure (or construct it manually)
	// For performance, constructing manually is better, but fetching guarantees data integrity.
	// Let's refetch for simplicity and correctness of relations.
	return s.repo.FindByID(ctx, drama.ID)
}

func (s *service) Update(ctx context.Context, id string, req UpdateDramaRequest) (*Drama, error) {
	drama, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if drama == nil {
		return nil, errors.New("drama not found")
	}

	// Update fields
	drama.Title = req.Title
	drama.Synopsis = req.Synopsis
	drama.PosterURL = req.PosterURL
	drama.Year = req.Year
	drama.TotalSeasons = req.TotalSeasons
	drama.Status = req.Status
	drama.SourceURL = req.SourceURL

	if err := s.repo.Update(ctx, drama, req.GenreIDs, req.Actors); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
