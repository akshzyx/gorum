package handlers

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetPublicProfile(
	w http.ResponseWriter,
	r *http.Request,
) {
	username := chi.URLParam(r, "username")

	profile, err := h.service.GetPublicProfile(r.Context(), username)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	util.WriteJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	resp, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
