package genre

import (
	"context"
	"drakor-backend/pkg/database"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Genre, error)
	FindByID(ctx context.Context, id string) (*Genre, error)
	FindBySlug(ctx context.Context, slug string) (*Genre, error)
	Create(ctx context.Context, genre *Genre) error
	Update(ctx context.Context, genre *Genre) error
	Delete(ctx context.Context, id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindAll(ctx context.Context) ([]Genre, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, name, slug FROM genres ORDER BY name ASC`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var g Genre
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	return genres, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Genre, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, name, slug FROM genres WHERE id = $1`
	var g Genre
	err := db.QueryRow(ctx, query, id).Scan(&g.ID, &g.Name, &g.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*Genre, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	query := `SELECT id, name, slug FROM genres WHERE slug = $1`
	var g Genre
	err := db.QueryRow(ctx, query, slug).Scan(&g.ID, &g.Name, &g.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *repository) Create(ctx context.Context, genre *Genre) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `INSERT INTO genres (name, slug) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(ctx, query, genre.Name, genre.Slug).Scan(&genre.ID)
}

func (r *repository) Update(ctx context.Context, genre *Genre) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `UPDATE genres SET name = $1, slug = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, genre.Name, genre.Slug, genre.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	query := `DELETE FROM genres WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}
