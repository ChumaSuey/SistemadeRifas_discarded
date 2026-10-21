package handler

import (
	"encoding/json"
	"net/http"

	"raffle-saas/internal/repository"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)
	users, total, err := h.userRepo.List(r.Context(), page, perPage)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	type userResp struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}

	data := make([]userResp, len(users))
	for i, u := range users {
		data[i] = userResp{
			ID:     u.ID.Hex(),
			Name:   u.Name,
			Email:  u.Email,
			Role:   u.Role,
			Status: u.Status,
		}
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":        data,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", ""))
		return
	}

	if req.Role == "" {
		utils.WriteError(w, utils.NewAPIError(400, "role is required", ""))
		return
	}

	if err := h.userRepo.UpdateRole(r.Context(), userID, req.Role); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

func (h *AdminHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", ""))
		return
	}

	if req.Status == "" {
		utils.WriteError(w, utils.NewAPIError(400, "status is required", ""))
		return
	}

	if err := h.userRepo.UpdateStatus(r.Context(), userID, req.Status); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}
