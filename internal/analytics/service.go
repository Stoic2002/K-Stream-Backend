package analytics

import "context"

type Service interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx)
}
