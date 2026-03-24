package container

import (
	"context"

	"github.com/akshzyx/gorum/internal/api"
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/config"
	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/auth"
	"github.com/akshzyx/gorum/internal/domain/follow"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/akshzyx/gorum/internal/infra/email"
	"github.com/akshzyx/gorum/internal/repository/postgres"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Router *chi.Mux
	DB     *pgxpool.Pool
}

func NewContainer(cfg *config.Config) (*Container, error) {
	ctx := context.Background()

	// DB
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)

	// JWT
	util.InitJWT(cfg.JWTSecret, cfg.JWTExpiresIn)

	// EMAIL
	emailSender := email.NewEmailSender(
		cfg.ResendAPIKey,
		cfg.ResendFrom,
		cfg.AppBaseURL,
	)

	// REPOSITORIES
	userRepo := postgres.NewUserRepository(queries)
	authRepo := postgres.NewAuthRepo(queries, pool)
	postRepo := postgres.NewPostRepository(queries)
	followRepo := postgres.NewFollowRepository(queries)

	// SERVICES
	userService := user.NewService(userRepo)
	authService := auth.NewService(authRepo, emailSender)
	postService := post.NewService(postRepo)
	followService := follow.NewService(followRepo)

	// HANDLERS
	userHandler := handlers.NewUserHandler(userService)
	settingsHandler := handlers.NewSettingsHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)
	followHandler := handlers.NewFollowHandler(followService)

	// ROUTER
	router := api.NewRouter(
		authHandler,
		userHandler,
		settingsHandler,
		postHandler,
		followHandler,
	)

	return &Container{
		Router: router,
		DB:     pool,
	}, nil
}
