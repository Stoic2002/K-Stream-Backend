package auth

import (
	"context"
	"errors"
	"time"

	"drakor-backend/pkg/database"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	// Admin
	FindAll(ctx context.Context, limit, offset int) ([]User, int64, error)
	UpdateRole(ctx context.Context, userID, role string) error
	Delete(ctx context.Context, userID string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	query := `
		INSERT INTO users (name, email, password_hash, role, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := db.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.AvatarURL,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)

	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}
	query := `SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}
	query := `SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE id = $1`

	var user User
	err := db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	query := `
		UPDATE users 
		SET name = $1, avatar_url = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at
	`
	err := db.QueryRow(ctx, query,
		user.Name,
		user.AvatarURL,
		time.Now(),
		user.ID,
	).Scan(&user.UpdatedAt)

	return err
}

func (r *repository) FindAll(ctx context.Context, limit, offset int) ([]User, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	var total int64
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, email, role, avatar_url, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *repository) UpdateRole(ctx context.Context, userID, role string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	_, err := db.Exec(ctx, "UPDATE users SET role = $1, updated_at = $2 WHERE id = $3", role, time.Now(), userID)
	return err
}

func (r *repository) Delete(ctx context.Context, userID string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}
	_, err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}
