package handler

import (
	"encoding/json"
	"net/http"

	"raffle-saas/internal/dto"
	"raffle-saas/internal/middleware"
	"raffle-saas/internal/service"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.WriteError(w, utils.NewAPIError(400, "name, email and password are required", ""))
		return
	}

	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.WriteError(w, utils.NewAPIError(400, "email and password are required", ""))
		return
	}

	resp, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) OturnhTwitch(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Twitch OAuth flow (see PENDING.md)
	utils.WriteError(w, utils.NewAPIError(501, "not implemented yet - Twitch OAuth pending", ""))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	if req.RefreshToken == "" {
		utils.WriteError(w, utils.NewAPIError(400, "refresh_token is required", ""))
		return
	}

	resp, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	if req.RefreshToken == "" {
		utils.WriteError(w, utils.NewAPIError(400, "refresh_token is required", ""))
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	resp, err := h.svc.GetMe(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
