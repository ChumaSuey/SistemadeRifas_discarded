package handler

import (
	"io"
	"log/slog"
	"net/http"

	"raffle-saas/internal/utils"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) Stripe(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read stripe webhook body", "error", err)
		utils.WriteError(w, utils.NewAPIError(400, "invalid request body", ""))
		return
	}

	slog.Info("stripe webhook received", "body_length", len(body))

	// TODO: Implement full Stripe webhook handling (see PENDING.md)
	// 1. Verify webhook signature using STRIPE_WEBHOOK_SECRET
	// 2. Parse event type
	// 3. Handle payment_intent.succeeded → update ticket status to verified
	// 4. Handle payment_intent.payment_failed → cancel ticket

	w.WriteHeader(http.StatusOK)
}
