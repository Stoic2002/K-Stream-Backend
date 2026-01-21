package season

import (
	"context"
	"errors"
)

type Service interface {
	GetByDramaID(ctx context.Context, dramaID string) ([]Season, error)
	GetByID(ctx context.Context, id string) (*Season, error)
	Create(ctx context.Context, req CreateSeasonRequest) (*Season, error)
	Update(ctx context.Context, id string, req UpdateSeasonRequest) (*Season, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByDramaID(ctx context.Context, dramaID string) ([]Season, error) {
	return s.repo.FindByDramaID(ctx, dramaID)
}

func (s *service) GetByID(ctx context.Context, id string) (*Season, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, req CreateSeasonRequest) (*Season, error) {
	season := &Season{
		DramaID:      req.DramaID,
		SeasonNumber: req.SeasonNumber,
		Title:        req.Title,
	}

	if err := s.repo.Create(ctx, season); err != nil {
		return nil, err
	}

	return season, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateSeasonRequest) (*Season, error) {
	season, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if season == nil {
		return nil, errors.New("season not found")
	}

	season.SeasonNumber = req.SeasonNumber
	season.Title = req.Title

	if err := s.repo.Update(ctx, season); err != nil {
		return nil, err
	}

	return season, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
