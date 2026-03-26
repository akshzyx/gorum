package handlers

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/akshzyx/gorum/internal/util"
)

type SettingsHandler struct {
	service *user.Service
}

func NewSettingsHandler(service *user.Service) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var req user.UpdateEmailRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	err := h.service.UpdateEmail(r.Context(), userID, req)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SettingsHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req user.UpdatePasswordRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	err := h.service.UpdatePassword(r.Context(), userID, req)
	if err != nil {
		switch err {
		case user.ErrInvalidPassword:
			util.Unauthorized(w, r, err)
		case user.ErrNotFound:
			util.NotFound(w, r)
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req user.UpdateProfileRequest

	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
