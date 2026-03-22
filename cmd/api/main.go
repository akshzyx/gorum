package main

import (
	"context"
	"log"
	"net/http"

	"github.com/akshzyx/gorum/internal/api"
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/config"
	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/auth"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/akshzyx/gorum/internal/infra/email"
	"github.com/akshzyx/gorum/internal/repository/postgres"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// LOAD CONFIG
	cfg := config.LoadConfig()

	// CONNECT DATABASE
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	// SQLC QUERIES
	queries := db.New(pool)

	// INIT JWT
	util.InitJWT(cfg.JWTSecret, cfg.JWTExpiresIn)

	// EMAIL SENDER
	emailSender := email.NewEmailSender(
		cfg.ResendAPIKey,
		cfg.ResendFrom,
		cfg.AppBaseURL,
	)

	// USER DOMAIN
	userRepo := postgres.NewUserRepository(queries)
	userService := user.NewService(userRepo)

	userHandler := handlers.NewUserHandler(userService)
	settingsHandler := handlers.NewSettingsHandler(userService)

	// AUTH DOMAIN
	authRepo := postgres.NewAuthRepo(queries, pool)
	authService := auth.NewService(authRepo, emailSender)
	authHandler := handlers.NewAuthHandler(authService)

	// POST DOMAIN
	postRepo := postgres.NewPostRepository(queries)
	postService := post.NewService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	// ROUTER
	router := api.NewRouter(
		authHandler,
		userHandler,
		settingsHandler,
		postHandler,
	)

	// START SERVER
	log.Printf("🚀 Server running on :%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
