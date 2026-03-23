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
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
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

	util.WriteJSON(w, http.StatusOK, post.PublicPostResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	})
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
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

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = int32(v)
		}
	}

	posts, err := h.service.ListLatest(r.Context(), limit)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	resp := make([]post.PublicPostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, post.PublicPostResponse{
			ID:        p.ID,
			UserID:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *PostHandler) Reply(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
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

	posts, err := h.service.ListReplies(r.Context(), postID)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}

	resp := make([]post.ReplyResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, post.ReplyResponse{
			ID:        p.ID,
			UserID:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	util.WriteJSON(w, http.StatusOK, resp)
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
			ID:        p.ID,
			UserID:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
