package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
)

func PostRoutes(h *handlers.PostHandler) chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Get("/", h.ListLatest)
	r.Get("/{id}", h.GetByID)

	// Replies (public read)
	r.Get("/{id}/replies", h.ListReplies)
	r.Get("/{id}/thread", h.GetThread)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)

		r.Post("/", h.Create)
		r.Delete("/{id}", h.Delete)

		// Reply to a post
		r.Post("/{id}/reply", h.Reply)

		// Like system
		r.Post("/{id}/like", h.LikePost)
		r.Delete("/{id}/like", h.UnlikePost)
	})

	return r
}
