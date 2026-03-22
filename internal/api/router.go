package api

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/api/routes"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	settingsHandler *handlers.SettingsHandler,
	postHandler *handlers.PostHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// HEALTH CHECK
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication
		r.Mount("/auth", routes.AuthRoutes(authHandler))

		// Public user profiles
		r.Mount("/user", routes.UserRoutes(userHandler))

		// Posts (tweets + replies)
		r.Mount("/post", routes.PostRoutes(postHandler))

		// User settings (JWT required)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTAuth)
			r.Mount("/settings", routes.SettingsRoutes(settingsHandler))
		})
	})

	return r
}
