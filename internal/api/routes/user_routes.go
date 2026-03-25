package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(h *handlers.UserHandler) chi.Router {
	r := chi.NewRouter()

	// Public profile (read-only)
	r.Get("/{username}", h.GetPublicProfile)

	// user posts
	r.Get("/{username}/posts", h.GetUserPosts)

	// user replies)
	r.Get("/{username}/replies", h.GetUserReplies)

	return r
}
