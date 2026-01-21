package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"drakor-backend/pkg/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load env
	if err := godotenv.Load("../.env"); err != nil {
		// Try root .env if not found in parent (running from root)
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	// Connect DB
	database.Connect()
	defer database.Close()
	db := database.GetDB()

	ctx := context.Background()

	fmt.Println("🌱 Seeding database...")

	// 1. Seed Genres
	seedGenres(ctx, db)

	// 2. Seed Users (Admin & User)
	adminID := seedUsers(ctx, db)

	// 3. Seed Dramas (with Genres & Actors)
	seedDramas(ctx, db, adminID)

	fmt.Println("✅ Seeding complete!")
}

func seedGenres(ctx context.Context, db *pgxpool.Pool) {
	genres := []struct {
		Name string
		Slug string
	}{
		{"Romance", "romance"},
		{"Action", "action"},
		{"Comedy", "comedy"},
		{"Fantasy", "fantasy"},
		{"Thriller", "thriller"},
		{"Historical", "historical"},
	}

	fmt.Printf("... Seeding %d genres\n", len(genres))
	for _, g := range genres {
		_, err := db.Exec(ctx, "INSERT INTO genres (name, slug) VALUES ($1, $2) ON CONFLICT (slug) DO NOTHING", g.Name, g.Slug)
		if err != nil {
			log.Printf("Error seeding genre %s: %v\n", g.Name, err)
		}
	}
}

func seedUsers(ctx context.Context, db *pgxpool.Pool) string {
	// Create Admin
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	var adminID string
	err := db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO UPDATE SET role = 'admin'
		RETURNING id
	`, "admin@drakor.com", string(hashedPwd), "Admin Drakor", "admin", time.Now(), time.Now()).Scan(&adminID)

	if err != nil {
		log.Printf("Error seeding admin: %v\n", err)
	} else {
		fmt.Println("... Seeded Admin User (admin@drakor.com / password123)")
	}

	// Create User
	_, err = db.Exec(ctx, `
		INSERT INTO users (email, password_hash, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO NOTHING
	`, "user@drakor.com", string(hashedPwd), "Regular User", "user", time.Now(), time.Now())
	if err != nil {
		log.Printf("Error seeding user: %v\n", err)
	} else {
		fmt.Println("... Seeded Regular User (user@drakor.com / password123)")
	}

	return adminID
}

func seedDramas(ctx context.Context, db *pgxpool.Pool, adminID string) {
	// Only seed if empty to avoid dups logic for now or simple check
	var count int
	db.QueryRow(ctx, "SELECT COUNT(*) FROM dramas").Scan(&count)
	if count > 0 {
		fmt.Println("... Dramas already exist, skipping drama/season/episode seeding")
		return
	}

	fmt.Println("... Seeding Dramas...")

	// Drama 1
	var dramaID string
	err := db.QueryRow(ctx, `
		INSERT INTO dramas (title, synopsis, year, rating, total_seasons, status, added_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, "Goblin: The Lonely and Great God", "In ancient times, an invincible general is betrayed...", 2016, 9.5, 1, "completed", adminID, time.Now(), time.Now()).Scan(&dramaID)

	if err == nil {
		// Attach Genre (Fantasy)
		var genreID string
		db.QueryRow(ctx, "SELECT id FROM genres WHERE slug='fantasy'").Scan(&genreID)
		if genreID != "" {
			db.Exec(ctx, "INSERT INTO drama_genres (drama_id, genre_id) VALUES ($1, $2)", dramaID, genreID)
		}

		// Add Season 1
		var seasonID string
		db.QueryRow(ctx, `
			INSERT INTO seasons (drama_id, season_number, title, created_at)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, dramaID, 1, "Season 1", time.Now()).Scan(&seasonID)

		// Add Episode 1
		db.Exec(ctx, `
			INSERT INTO episodes (season_id, episode_number, title, video_url, duration, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, seasonID, 1, "Episode 1", "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_1mb.mp4", 3600, time.Now())
	}

	// Drama 2
	db.Exec(ctx, `
		INSERT INTO dramas (title, synopsis, year, rating, total_seasons, status, added_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, "Queen of Tears", "The queen of department stores and the prince of supermarkets...", 2024, 8.8, 1, "ongoing", adminID, time.Now(), time.Now())
}
