package handlers

import (
	"net/http"

	"github.com/akshzyx/gorum/internal/domain/auth"
	"github.com/akshzyx/gorum/internal/util"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req auth.SignupRequest

	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	resp, err := h.service.Signup(r.Context(), req)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest

	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var req auth.ActivateRequest

	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := h.service.Activate(r.Context(), req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "account activated",
	})
}

func (h *AuthHandler) ResendActivation(w http.ResponseWriter, r *http.Request) {
	var req auth.ResendActivationRequest

	if err := util.ReadJSON(r, &req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	// TODO: implement fully later
	util.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "email resent (not implemented yet)",
	})
}
