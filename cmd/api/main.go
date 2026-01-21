package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"drakor-backend/internal/actor"
	"drakor-backend/internal/analytics"
	"drakor-backend/internal/auth"
	"drakor-backend/internal/comment"
	"drakor-backend/internal/drama"
	"drakor-backend/internal/episode"
	"drakor-backend/internal/genre"
	"drakor-backend/internal/history"
	"drakor-backend/internal/review"
	"drakor-backend/internal/season"
	"drakor-backend/internal/watchlist"
	"drakor-backend/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Printf("⚠️  Warning: Database connection failed: %v", err)
		log.Println("   Server will start without database connection")
	} else {
		defer database.Close()
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "not_configured"
		if db := database.GetDB(); db != nil {
			if err := db.Ping(c.Request.Context()); err != nil {
				dbStatus = "disconnected"
			} else {
				dbStatus = "connected"
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "Drakor API is running",
			"database": dbStatus,
		})
	})

	// Initialize Authn dependencies
	authRepo := auth.NewRepository()
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Initialize Genre dependencies
	genreRepo := genre.NewRepository()
	genreService := genre.NewService(genreRepo)
	genreHandler := genre.NewHandler(genreService)

	// Initialize Actor dependencies
	actorRepo := actor.NewRepository()
	actorService := actor.NewService(actorRepo)
	actorHandler := actor.NewHandler(actorService)

	// Initialize Drama dependencies
	dramaRepo := drama.NewRepository()
	dramaService := drama.NewService(dramaRepo)
	dramaHandler := drama.NewHandler(dramaService)

	// Initialize Season dependencies
	seasonRepo := season.NewRepository()
	seasonService := season.NewService(seasonRepo)
	seasonHandler := season.NewHandler(seasonService)

	// Initialize Episode dependencies
	episodeRepo := episode.NewRepository()
	episodeService := episode.NewService(episodeRepo)
	episodeHandler := episode.NewHandler(episodeService)

	// Initialize Watchlist dependencies
	watchlistRepo := watchlist.NewRepository()
	watchlistService := watchlist.NewService(watchlistRepo)
	watchlistHandler := watchlist.NewHandler(watchlistService)

	// Initialize History dependencies
	historyRepo := history.NewRepository()
	historyService := history.NewService(historyRepo)
	historyHandler := history.NewHandler(historyService)

	// Initialize Review dependencies
	reviewRepo := review.NewRepository()
	reviewService := review.NewService(reviewRepo)
	reviewHandler := review.NewHandler(reviewService)

	// Initialize Comment dependencies
	commentRepo := comment.NewRepository()
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService)

	// Initialize Analytics dependencies
	analyticsRepo := analytics.NewRepository()
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)

	// API routes group
	api := r.Group("/api")
	{
		// Public routes
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// --- GENRE Routes ---
		// Public
		api.GET("/genres", genreHandler.GetAll)

		// Admin
		genreGroup := api.Group("/genres")
		genreGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			genreGroup.POST("", genreHandler.Create)
			genreGroup.PUT("/:id", genreHandler.Update)
			genreGroup.DELETE("/:id", genreHandler.Delete)
		}

		// --- ACTOR Routes ---
		// Public
		api.GET("/actors", actorHandler.GetAll)
		api.GET("/actors/:id", actorHandler.GetByID)

		// Admin
		actorGroup := api.Group("/actors")
		actorGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			actorGroup.POST("", actorHandler.Create)
			actorGroup.PUT("/:id", actorHandler.Update)
			actorGroup.DELETE("/:id", actorHandler.Delete)
		}

		// --- DRAMA Routes ---
		// Public
		api.GET("/dramas", dramaHandler.GetAll)
		api.GET("/dramas/:id", dramaHandler.GetByID)

		// Admin
		dramaGroup := api.Group("/dramas")
		dramaGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			dramaGroup.POST("", dramaHandler.Create)
			dramaGroup.PUT("/:id", dramaHandler.Update)
			dramaGroup.DELETE("/:id", dramaHandler.Delete)
		}

		// --- SEASON Routes ---
		// Public
		api.GET("/dramas/:id/seasons", seasonHandler.GetByDramaID)
		api.GET("/seasons/:id", seasonHandler.GetByID)

		// Admin
		seasonGroup := api.Group("/seasons")
		seasonGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			seasonGroup.POST("", seasonHandler.Create)
			seasonGroup.PUT("/:id", seasonHandler.Update)
			seasonGroup.DELETE("/:id", seasonHandler.Delete)
		}

		// --- EPISODE Routes ---
		// Public
		api.GET("/seasons/:id/episodes", episodeHandler.GetBySeasonID)
		api.GET("/episodes/:id", episodeHandler.GetByID)

		// Admin
		episodeGroup := api.Group("/episodes")
		episodeGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			episodeGroup.POST("", episodeHandler.Create)
			episodeGroup.PUT("/:id", episodeHandler.Update)
			episodeGroup.DELETE("/:id", episodeHandler.Delete)
		}

		// --- WATCHLIST Routes ---
		// Protected (User)
		watchlistGroup := api.Group("/watchlist")
		watchlistGroup.Use(auth.Middleware())
		{
			watchlistGroup.GET("", watchlistHandler.GetMine)
			watchlistGroup.POST("", watchlistHandler.Add)
			watchlistGroup.DELETE("/:dramaID", watchlistHandler.Remove)
			watchlistGroup.GET("/:dramaID/check", watchlistHandler.Check)
		}

		// --- HISTORY Routes ---
		// Protected (User)
		historyGroup := api.Group("/history")
		historyGroup.Use(auth.Middleware())
		{
			historyGroup.GET("", historyHandler.GetMine)
			historyGroup.POST("", historyHandler.Record)
			historyGroup.GET("/:episodeID", historyHandler.GetProgress)
		}

		// --- REVIEW Routes ---
		// Public
		api.GET("/dramas/:id/reviews", reviewHandler.GetByDrama)

		// Protected (User)
		reviewGroup := api.Group("/reviews")
		reviewGroup.Use(auth.Middleware())
		{
			reviewGroup.POST("", reviewHandler.Create)
			reviewGroup.PUT("/:id", reviewHandler.Update)
			reviewGroup.DELETE("/:id", reviewHandler.Delete)
		}

		// --- COMMENT Routes ---
		// Public
		api.GET("/episodes/:id/comments", commentHandler.GetByEpisode)

		// Protected (User)
		commentGroup := api.Group("/comments")
		commentGroup.Use(auth.Middleware())
		{
			commentGroup.POST("", commentHandler.Create)
			commentGroup.PUT("/:id", commentHandler.Update)
			commentGroup.DELETE("/:id", commentHandler.Delete)
		}

		// --- ANALYTICS Routes ---
		// Admin
		analyticsGroup := api.Group("/analytics")
		analyticsGroup.Use(auth.Middleware(), auth.AdminMiddleware())
		{
			analyticsGroup.GET("/dashboard", analyticsHandler.GetDashboard)

			// User Management
			analyticsGroup.GET("/users", authHandler.GetAllUsers)
			analyticsGroup.PATCH("/users/:id/role", authHandler.UpdateUserRole)
			analyticsGroup.DELETE("/users/:id", authHandler.DeleteUser)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)

			// Protected routes
			protected := authGroup.Use(auth.Middleware())
			{
				protected.GET("/me", authHandler.GetProfile)
				protected.PUT("/profile", authHandler.UpdateProfile)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Close database connection
	database.Close()

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
