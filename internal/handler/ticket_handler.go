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

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

func (h *TicketHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	var req dto.PurchaseTicketsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	tickets, err := h.svc.Purchase(r.Context(), raffleID, userID, &req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, tickets)
}

func (h *TicketHandler) ListByRaffle(w http.ResponseWriter, r *http.Request) {
	raffleID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid raffle id", ""))
		return
	}

	page, perPage := parsePagination(r)
	resp, err := h.svc.ListByRaffle(r.Context(), raffleID, page, perPage)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *TicketHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	page, perPage := parsePagination(r)
	resp, err := h.svc.ListMy(r.Context(), userID, page, perPage)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *TicketHandler) UploadProof(w http.ResponseWriter, r *http.Request) {
	userID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	ticketID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid ticket id", ""))
		return
	}

	var req dto.UploadProofRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	resp, err := h.svc.UploadProof(r.Context(), ticketID, userID, req.PaymentProofURL)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *TicketHandler) Verify(w http.ResponseWriter, r *http.Request) {
	verifierID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	ticketID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid ticket id", ""))
		return
	}

	resp, err := h.svc.Verify(r.Context(), ticketID, verifierID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *TicketHandler) Reject(w http.ResponseWriter, r *http.Request) {
	verifierID, err := bson.ObjectIDFromHex(middleware.GetUserID(r.Context()))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid user id", ""))
		return
	}

	ticketID, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid ticket id", ""))
		return
	}

	var req dto.RejectTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", err.Error()))
		return
	}

	resp, err := h.svc.Reject(r.Context(), ticketID, verifierID, req.Reason)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
