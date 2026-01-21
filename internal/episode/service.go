package episode

import (
	"context"
	"errors"
)

type Service interface {
	GetBySeasonID(ctx context.Context, seasonID string) ([]Episode, error)
	GetByID(ctx context.Context, id string) (*Episode, error)
	Create(ctx context.Context, userID string, req CreateEpisodeRequest) (*Episode, error)
	Update(ctx context.Context, id string, req UpdateEpisodeRequest) (*Episode, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetBySeasonID(ctx context.Context, seasonID string) ([]Episode, error) {
	return s.repo.FindBySeasonID(ctx, seasonID)
}

func (s *service) GetByID(ctx context.Context, id string) (*Episode, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, userID string, req CreateEpisodeRequest) (*Episode, error) {
	episode := &Episode{
		SeasonID:      req.SeasonID,
		EpisodeNumber: req.EpisodeNumber,
		Title:         req.Title,
		VideoURL:      req.VideoURL,
		Duration:      req.Duration,
		ThumbnailURL:  req.ThumbnailURL,
		SourceURL:     req.SourceURL,
		AddedBy:       userID,
	}

	if err := s.repo.Create(ctx, episode); err != nil {
		return nil, err
	}

	return episode, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateEpisodeRequest) (*Episode, error) {
	episode, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if episode == nil {
		return nil, errors.New("episode not found")
	}

	episode.EpisodeNumber = req.EpisodeNumber
	episode.Title = req.Title
	episode.VideoURL = req.VideoURL
	episode.Duration = req.Duration
	episode.ThumbnailURL = req.ThumbnailURL
	episode.SourceURL = req.SourceURL

	if err := s.repo.Update(ctx, episode); err != nil {
		return nil, err
	}

	return episode, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
