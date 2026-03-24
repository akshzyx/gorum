package api

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/api/routes"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	settingsHandler *handlers.SettingsHandler,
	postHandler *handlers.PostHandler,
	followHandler *handlers.FollowHandler,
) *chi.Mux {
	r := chi.NewRouter()

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
		// public user profile routes
		r.Route("/user", func(r chi.Router) {
			r.Mount("/", routes.UserRoutes(userHandler))
			r.Mount("/", routes.FollowRoutes(followHandler))
		})
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
