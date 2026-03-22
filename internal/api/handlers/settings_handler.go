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

// func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
// 	var req user.UpdateProfileRequest
// 	if err := util.ReadJSON(r, &req); err != nil {
// 		util.BadRequest(w, r, err)
// 		return
// 	}

// 	ctxUserID := r.Context().Value(middlewares.UserIDKey)
// 	if ctxUserID == nil {
// 		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
// 		return
// 	}

// 	if err := h.service.UpdateProfile(r.Context(), ctxUserID.(string), req); err != nil {
// 		util.InternalServerError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

func (h *SettingsHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var req user.UpdateEmailRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	ctxUserID := r.Context().Value(middlewares.UserIDKey)
	if ctxUserID == nil {
		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.UpdateEmail(r.Context(), ctxUserID.(string), req); err != nil {
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

	ctxUserID := r.Context().Value(middlewares.UserIDKey)
	if ctxUserID == nil {
		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.UpdatePassword(r.Context(), ctxUserID.(string), req); err != nil {
		util.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
