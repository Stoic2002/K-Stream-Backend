package analytics

import (
	"context"
	"drakor-backend/pkg/database"
	"errors"
)

type Repository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	stats := &DashboardStats{}

	// Parallel or sequential queries
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers); err != nil {
		return nil, err
	}
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM dramas").Scan(&stats.TotalDramas); err != nil {
		return nil, err
	}
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM episodes").Scan(&stats.TotalEpisodes); err != nil {
		return nil, err
	}
	// Total views could be sum of view_count from dramas or episodes
	if err := db.QueryRow(ctx, "SELECT COALESCE(SUM(view_count), 0) FROM dramas").Scan(&stats.TotalViews); err != nil {
		return nil, err
	}

	return stats, nil
}
