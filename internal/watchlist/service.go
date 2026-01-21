package watchlist

import "context"

type Service interface {
	AddToWatchlist(ctx context.Context, userID, dramaID string) error
	RemoveFromWatchlist(ctx context.Context, userID, dramaID string) error
	GetMyWatchlist(ctx context.Context, userID string, page, limit int) ([]WatchlistItem, int64, error)
	CheckIsWatchlisted(ctx context.Context, userID, dramaID string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddToWatchlist(ctx context.Context, userID, dramaID string) error {
	return s.repo.Add(ctx, userID, dramaID)
}

func (s *service) RemoveFromWatchlist(ctx context.Context, userID, dramaID string) error {
	return s.repo.Remove(ctx, userID, dramaID)
}

func (s *service) GetMyWatchlist(ctx context.Context, userID string, page, limit int) ([]WatchlistItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetByUser(ctx, userID, limit, offset)
}

func (s *service) CheckIsWatchlisted(ctx context.Context, userID, dramaID string) (bool, error) {
	return s.repo.Exists(ctx, userID, dramaID)
}
