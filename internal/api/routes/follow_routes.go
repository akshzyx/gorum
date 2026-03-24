package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
)

func FollowRoutes(h *handlers.FollowHandler) chi.Router {
	r := chi.NewRouter()

	// Protected actions
	r.With(middlewares.RequireAuth).Post("/{userID}/follow", h.FollowUser)
	r.With(middlewares.RequireAuth).Delete("/{userID}/follow", h.UnfollowUser)

	// Public
	r.Get("/{userID}/followers", h.GetFollowers)
	r.Get("/{userID}/following", h.GetFollowing)

	return r
}
