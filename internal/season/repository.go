package season

import (
	"context"
	"drakor-backend/pkg/database"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	FindByDramaID(ctx context.Context, dramaID string) ([]Season, error)
	FindByID(ctx context.Context, id string) (*Season, error)
	Create(ctx context.Context, season *Season) error
	Update(ctx context.Context, season *Season) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindByDramaID(ctx context.Context, dramaID string) ([]Season, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, drama_id, season_number, title, created_at FROM seasons WHERE drama_id = $1 ORDER BY season_number ASC`
	rows, err := db.Query(ctx, query, dramaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []Season
	for rows.Next() {
		var s Season
		if err := rows.Scan(&s.ID, &s.DramaID, &s.SeasonNumber, &s.Title, &s.CreatedAt); err != nil {
			return nil, err
		}
		seasons = append(seasons, s)
	}

	return seasons, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Season, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, drama_id, season_number, title, created_at FROM seasons WHERE id = $1`
	var s Season
	err := db.QueryRow(ctx, query, id).Scan(&s.ID, &s.DramaID, &s.SeasonNumber, &s.Title, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) Create(ctx context.Context, season *Season) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `INSERT INTO seasons (drama_id, season_number, title, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return db.QueryRow(ctx, query, season.DramaID, season.SeasonNumber, season.Title, time.Now()).Scan(&season.ID)
}

func (r *repository) Update(ctx context.Context, season *Season) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `UPDATE seasons SET season_number = $1, title = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, season.SeasonNumber, season.Title, season.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `DELETE FROM seasons WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}
