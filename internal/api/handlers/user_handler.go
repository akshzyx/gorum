package handlers

import (
	"net/http"
	"strconv"
	"time"

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
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	// get user (you already have this logic)
	user, err := h.service.GetPublicProfile(ctx, username)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	limit := int32(20)

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = int32(v)
		}
	}

	// parse cursor
	var cursor *time.Time
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr != "" {
		if t, err := time.Parse(time.RFC3339, cursorStr); err == nil {
			cursor = &t
		}
	}

	viewerID := middlewares.GetUserID(ctx)

	result, err := h.postService.GetUserPosts(ctx, user.ID, cursor, limit, viewerID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	// format created_at
	for _, p := range result.Data {
		if t, ok := p["created_at"].(time.Time); ok {
			p["created_at"] = t.Format(time.RFC3339)
		}
	}

	// format next_cursor
	var nextCursor *string
	if result.NextCursor != nil {
		s := result.NextCursor.Format(time.RFC3339)
		nextCursor = &s
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":        result.Data,
		"next_cursor": nextCursor,
		"has_more":    result.HasMore,
	})
}

func (h *UserHandler) GetUserReplies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	user, err := h.service.GetPublicProfile(ctx, username)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	limit := int32(20)

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = int32(v)
		}
	}

	// parse cursor
	var cursor *time.Time
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr != "" {
		if t, err := time.Parse(time.RFC3339, cursorStr); err == nil {
			cursor = &t
		}
	}

	viewerID := middlewares.GetUserID(ctx)

	result, err := h.postService.GetUserReplies(ctx, user.ID, cursor, limit, viewerID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	// format created_at
	for _, p := range result.Data {
		if t, ok := p["created_at"].(time.Time); ok {
			p["created_at"] = t.Format(time.RFC3339)
		}
	}

	// format next_cursor
	var nextCursor *string
	if result.NextCursor != nil {
		s := result.NextCursor.Format(time.RFC3339)
		nextCursor = &s
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":        result.Data,
		"next_cursor": nextCursor,
		"has_more":    result.HasMore,
	})
}
