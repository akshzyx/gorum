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
		switch err {
		case auth.ErrEmailExists:
			util.WriteJSONError(w, http.StatusConflict, err.Error())
		default:
			util.InternalServerError(w, r, err)
		}
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
		switch err {
		case auth.ErrInvalidCredentials:
			util.Unauthorized(w, r, err)
		case auth.ErrEmailNotVerified:
			util.WriteJSONError(w, http.StatusForbidden, err.Error())
		default:
			util.InternalServerError(w, r, err)
		}
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

	err := h.service.Activate(r.Context(), req)
	if err != nil {
		switch err {
		case auth.ErrTokenNotFound:
			util.WriteJSONError(w, http.StatusNotFound, err.Error())
		case auth.ErrTokenExpired:
			util.WriteJSONError(w, http.StatusBadRequest, err.Error())
		case auth.ErrTokenUsed:
			util.WriteJSONError(w, http.StatusBadRequest, err.Error())
		default:
			util.InternalServerError(w, r, err)
		}
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

	err := h.service.ResendActivation(r.Context(), req)
	if err != nil {
		switch err {
		case auth.ErrUserNotFound:
			util.WriteJSONError(w, http.StatusNotFound, err.Error())
		case auth.ErrAlreadyVerified:
			util.WriteJSONError(w, http.StatusBadRequest, err.Error())
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "activation email resent",
	})
}
