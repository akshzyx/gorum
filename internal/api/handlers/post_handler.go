package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akshzyx/gorum/internal/api/middlewares"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	service *post.Service
}

func NewPostHandler(service *post.Service) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		util.Unauthorized(w, r, nil)
		return
	}

	var req post.CreatePostRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	id, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		switch err {
		case post.ErrInvalidContent:
			util.BadRequest(w, r, err)
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, post.CreatePostResponse{ID: id})
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	userID := middlewares.GetUserID(r.Context())

	liked := false
	if userID != "" {
		if l, err := h.service.HasUserLiked(r.Context(), userID, p.ID); err == nil {
			liked = l
		}
	}

	count := int64(0)
	if c, err := h.service.GetLikesCount(r.Context(), p.ID); err == nil {
		count = c
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":         p.ID,
		"user_id":    p.UserID,
		"content":    p.Content,
		"created_at": p.CreatedAt.Format(time.RFC3339),
		"likes":      count,
		"liked":      liked,
	})
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		util.Unauthorized(w, r, nil)
		return
	}

	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id, userID)
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			util.NotFound(w, r)
		case post.ErrForbidden:
			util.WriteJSONError(w, http.StatusForbidden, err.Error())
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) ListLatest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		t, err := time.Parse(time.RFC3339, cursorStr)
		if err == nil {
			cursor = &t
		}
	}

	userID := middlewares.GetUserID(ctx)

	result, err := h.service.ListLatest(ctx, userID, cursor, limit)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	// format created_at inside data
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

// Like handlers
func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		util.Unauthorized(w, r, nil)
		return
	}

	postID := chi.URLParam(r, "id")

	if err := h.service.LikePost(r.Context(), userID, postID); err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		util.Unauthorized(w, r, nil)
		return
	}

	postID := chi.URLParam(r, "id")

	if err := h.service.UnlikePost(r.Context(), userID, postID); err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Replies
func (h *PostHandler) Reply(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
	if userID == "" {
		util.Unauthorized(w, r, nil)
		return
	}

	parentID := chi.URLParam(r, "id")

	var req post.CreateReplyRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	id, err := h.service.Reply(r.Context(), userID, parentID, &req)
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			util.NotFound(w, r)
		case post.ErrInvalidContent:
			util.BadRequest(w, r, err)
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, map[string]string{
		"id": id,
	})
}

func (h *PostHandler) ListReplies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := chi.URLParam(r, "id")

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

	userID := middlewares.GetUserID(ctx)

	result, err := h.service.ListReplies(ctx, userID, postID, cursor, limit)
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

	// format cursor
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

func (h *PostHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	rootID := chi.URLParam(r, "id")

	posts, err := h.service.GetThread(r.Context(), rootID)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	resp := make([]post.ReplyResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, post.ReplyResponse{
			ID:           p.ID,
			UserID:       p.UserID,
			Content:      p.Content,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			ParentPostID: p.ParentPostID,
		})
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
