package api

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/api/routes"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	settingsHandler *handlers.SettingsHandler,
	postHandler *handlers.PostHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CORS)

	// basic health check endpoint
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		util.WriteJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "gorum",
		})
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
			r.Use(middlewares.RequireAuth)

			r.Get("/me", userHandler.GetMe)

			r.Mount("/settings", routes.SettingsRoutes(settingsHandler))
		})
	})

	return r
}
