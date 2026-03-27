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
		"username":   p.Username,
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
	limit := int32(20)
	cursor := r.URL.Query().Get("cursor")

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = int32(v)
		}
	}

	posts, nextCursor, err := h.service.ListLatest(r.Context(), limit, cursor)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	userID := middlewares.GetUserID(r.Context())

	resp, err := h.service.EnrichPosts(r.Context(), userID, posts)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	// format created_at properly
	for _, p := range resp {
		if t, ok := p["created_at"].(time.Time); ok {
			p["created_at"] = t.Format(time.RFC3339)
		}
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":        resp,
		"next_cursor": nextCursor,
		"has_more":    nextCursor != "",
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

	// prevent duplicate likes
	alreadyLiked, err := h.service.HasUserLiked(r.Context(), userID, postID)
	if err == nil && alreadyLiked {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.service.LikePost(r.Context(), userID, postID); err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	count, _ := h.service.GetLikesCount(r.Context(), postID)

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"liked": true,
		"likes": count,
	})
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

	count, _ := h.service.GetLikesCount(r.Context(), postID)

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"liked": false,
		"likes": count,
	})
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
	postID := chi.URLParam(r, "id")

	limit := int32(20)
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	if order != "desc" {
		order = "asc"
	}

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = int32(v)
		}
	}

	posts, nextCursor, err := h.service.ListRepliesPaginated(r.Context(), postID, limit, cursor, order)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	resp := make([]post.ReplyResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, post.ReplyResponse{
			ID:        p.ID,
			UserID:    p.UserID,
			Username:  p.Username,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":        resp,
		"next_cursor": nextCursor,
		"has_more":    nextCursor != "",
	})
}

func (h *PostHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// get post first
	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		util.NotFound(w, r)
		return
	}

	// resolve root
	rootID := p.ID
	if p.RootPostID != nil {
		rootID = *p.RootPostID
	}

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
			Username:     p.Username,
			Content:      p.Content,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
			ParentPostID: p.ParentPostID,
		})
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":      resp,
		"target_id": id,
	})
}
