package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(h *handlers.UserHandler) chi.Router {
	r := chi.NewRouter()

	// Public profile (read-only)
	r.Get("/{username}", h.GetPublicProfile)

	return r
}
