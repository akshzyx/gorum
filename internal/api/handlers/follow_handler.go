package handlers

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/domain/follow"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
)

type FollowHandler struct {
	service *follow.Service
}

func NewFollowHandler(service *follow.Service) *FollowHandler {
	return &FollowHandler{service: service}
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followingID := chi.URLParam(r, "userID")

	followerID, ok := r.Context().Value("userID").(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	err := h.service.FollowUser(r.Context(), followerID, followingID)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "followed successfully",
	})
}

func (h *FollowHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followingID := chi.URLParam(r, "userID")

	followerID, ok := r.Context().Value("userID").(string)
	if !ok {
		util.Unauthorized(w, r, nil)
		return
	}

	err := h.service.UnfollowUser(r.Context(), followerID, followingID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "unfollowed successfully",
	})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	users, err := h.service.GetFollowers(r.Context(), userID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, users)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	users, err := h.service.GetFollowing(r.Context(), userID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, users)
}
