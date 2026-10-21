package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"raffle-saas/internal/dto"
	"raffle-saas/internal/middleware"
	"raffle-saas/internal/service"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RaffleHandler struct {
	svc *service.RaffleService
}

func NewRaffleHandler(svc *service.RaffleService) *RaffleHandler {
	return &RaffleHandler{svc: svc}
}

func (h *RaffleHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	var req dto.CreateRaffleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	resp, err := h.svc.Create(r.Context(), rafflerID, &req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *RaffleHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)
	resp, err := h.svc.ListPublic(r.Context(), page, perPage)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	raffleID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	resp, err := h.svc.GetByID(r.Context(), raffleID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	page, perPage := parsePagination(r)
	resp, err := h.svc.ListMy(r.Context(), rafflerID, page, perPage)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	var req dto.UpdateRaffleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	resp, err := h.svc.Update(r.Context(), raffleID, rafflerID, &req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) Publish(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	resp, err := h.svc.Publish(r.Context(), raffleID, rafflerID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	resp, err := h.svc.Cancel(r.Context(), raffleID, rafflerID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) Sort(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	resp, err := h.svc.ExecuteSort(r.Context(), raffleID, rafflerID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) CommitSort(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	resp, err := h.svc.CommitSort(r.Context(), raffleID, rafflerID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *RaffleHandler) ReSort(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	rafflerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	var body struct {
		Reason string `json:"disqualify_reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	resp, err := h.svc.ReSort(r.Context(), raffleID, rafflerID, body.Reason)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func parsePagination(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage
}
