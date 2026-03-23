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

	// basic health check endpoint
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// auth routes (signup, login, activation)
		r.Mount("/auth", routes.AuthRoutes(authHandler))

		// public user profile routes
		r.Mount("/user", routes.UserRoutes(userHandler))

		// post routes (currently public, can restrict later)
		r.Mount("/post", routes.PostRoutes(postHandler))

		// routes that require a logged-in user
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth) // enforce auth before reaching handlers
			r.Mount("/settings", routes.SettingsRoutes(settingsHandler))
		})
	})

	return r
}
