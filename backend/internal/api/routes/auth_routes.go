package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(h *handlers.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", h.Signup)
	r.Post("/login", h.Login)
	r.Post("/activate", h.Activate)
	r.Post("/resend-activation", h.ResendActivation)
	r.Post("/logout", h.Logout)

	return r
}
