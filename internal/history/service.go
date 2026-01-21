package history

import "context"

type Service interface {
	RecordProgress(ctx context.Context, userID, episodeID string, progress int, isFinished bool) error
	GetMyHistory(ctx context.Context, userID string, page, limit int) ([]WatchHistory, int64, error)
	GetEpisodeProgress(ctx context.Context, userID, episodeID string) (*WatchHistory, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RecordProgress(ctx context.Context, userID, episodeID string, progress int, isFinished bool) error {
	return s.repo.Upsert(ctx, userID, episodeID, progress, isFinished)
}

func (s *service) GetMyHistory(ctx context.Context, userID string, page, limit int) ([]WatchHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetByUser(ctx, userID, limit, offset)
}

func (s *service) GetEpisodeProgress(ctx context.Context, userID, episodeID string) (*WatchHistory, error) {
	return s.repo.GetByEpisode(ctx, userID, episodeID)
}
