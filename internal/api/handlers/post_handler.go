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
	UserID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	var req post.CreatePostRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}

	if err := util.ValidateStruct(req); err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, err.Error())
	}

	id, err := h.service.Create(r.Context(), UserID, &req)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, post.CreatePostResponse{ID: id})
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		util.WriteJSONError(w, http.StatusNotFound, "post not found")
		return
	}

	util.WriteJSON(w, http.StatusOK, post.PublicPostResponse{
		ID:        p.ID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	})
}

func (h *PostHandler) DeleteByOwner(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorised")
	}

	id := chi.URLParam(r, "id")

	if err := h.service.DeleteByOwner(r.Context(), id, UserID); err != nil {
		util.WriteJSONError(w, http.StatusUnauthorized, "not allowed")
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
		util.WriteJSONError(w, http.StatusInternalServerError, err.Error())
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
		util.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	parentID := chi.URLParam(r, "id")

	var req post.CreateReplyRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Reply(r.Context(), userID, parentID, &req)
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, err.Error())
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
		util.WriteJSONError(w, http.StatusInternalServerError, err.Error())
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
		util.WriteJSONError(w, http.StatusNotFound, "thread not found")
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
