package routes

import (
	"github.com/akshzyx/gorum/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func SettingsRoutes(h *handlers.SettingsHandler) chi.Router {
	r := chi.NewRouter()

	r.Patch("/profile", h.UpdateProfile)
	r.Patch("/email", h.UpdateEmail)
	r.Patch("/password", h.UpdatePassword)

	return r
}
