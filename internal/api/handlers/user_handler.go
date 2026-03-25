package handlers

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/akshzyx/gorum/internal/domain/user"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service     *user.Service
	postService *post.Service
}

func NewUserHandler(service *user.Service, postService *post.Service) *UserHandler {
	return &UserHandler{
		service:     service,
		postService: postService,
	}
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

func (h *UserHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	u, err := h.service.GetPublicProfile(r.Context(), username)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	posts, err := h.postService.GetUserPosts(r.Context(), u.ID, 20)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	userID := middlewares.GetUserID(r.Context())

	resp, err := h.postService.EnrichPosts(r.Context(), userID, posts)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserReplies(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	u, err := h.service.GetPublicProfile(r.Context(), username)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	posts, err := h.postService.GetUserReplies(r.Context(), u.ID, 20)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	userID := middlewares.GetUserID(r.Context())

	resp, err := h.postService.EnrichPosts(r.Context(), userID, posts)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
