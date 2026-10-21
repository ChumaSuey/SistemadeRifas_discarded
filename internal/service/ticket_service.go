package service

import (
	"context"
	"time"

	"raffle-saas/internal/domain"
	"raffle-saas/internal/dto"
	"raffle-saas/internal/repository"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
	raffleRepo *repository.RaffleRepository
	auditRepo  *repository.AuditRepository
	client     *mongo.Client
}

func NewTicketService(ticketRepo *repository.TicketRepository, raffleRepo *repository.RaffleRepository, auditRepo *repository.AuditRepository, client *mongo.Client) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		raffleRepo: raffleRepo,
		auditRepo:  auditRepo,
		client:     client,
	}
}

func (s *TicketService) Purchase(ctx context.Context, raffleID, participantID bson.ObjectID, req *dto.PurchaseTicketsRequest) ([]*dto.TicketResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.Status != domain.RaffleStatusActive {
		return nil, utils.NewAPIError(400, "raffle is not active", "")
	}

	if raffle.RafflerID == participantID {
		return nil, utils.NewAPIError(409, "raffle creator cannot participate in their own raffle", "")
	}

	if req.Quantity <= 0 {
		return nil, utils.NewAPIError(400, "quantity must be at least 1", "")
	}

	available := raffle.TicketsTotal - raffle.TicketsSold
	if req.Quantity > available {
		return nil, utils.NewAPIError(400, "not enough tickets available", "")
	}

	if req.PaymentMethod != domain.PaymentMethodManual && req.PaymentMethod != domain.PaymentMethodStripe {
		return nil, utils.NewAPIError(400, "invalid payment method", "")
	}

	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	var tickets []*domain.Ticket

	for i := 0; i < req.Quantity; i++ {
		ticket := &domain.Ticket{
			RaffleID:      raffleID,
			ParticipantID: participantID,
			Status:        domain.TicketStatusPending,
			PaymentMethod: req.PaymentMethod,
		}

		if err := s.ticketRepo.CreateInTransaction(ctx, session, ticket, raffleID); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	s.auditRepo.Log(ctx, raffleID, participantID, "tickets_purchased", map[string]interface{}{
		"quantity":       req.Quantity,
		"payment_method": req.PaymentMethod,
		"ticket_ids":     getTicketIDs(tickets),
	}, "", "")

	result := make([]*dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		dto := ticketToDTO(t)
		result[i] = &dto
	}

	return result, nil
}

func (s *TicketService) ListByRaffle(ctx context.Context, raffleID bson.ObjectID, page, perPage int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	tickets, total, err := s.ticketRepo.ListByRaffle(ctx, raffleID, page, perPage)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		dtos[i] = ticketToDTO(t)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Data:       dtos,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TicketService) ListMy(ctx context.Context, participantID bson.ObjectID, page, perPage int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	tickets, total, err := s.ticketRepo.ListByParticipant(ctx, participantID, page, perPage)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		dtos[i] = ticketToDTO(t)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Data:       dtos,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TicketService) UploadProof(ctx context.Context, ticketID, userID bson.ObjectID, proofURL string) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if ticket.ParticipantID != userID {
		return nil, utils.ErrForbidden
	}

	if ticket.Status != domain.TicketStatusPending {
		return nil, utils.NewAPIError(400, "ticket must be in pending status to upload proof", "")
	}

	if ticket.PaymentMethod != domain.PaymentMethodManual {
		return nil, utils.NewAPIError(400, "only manual transfer tickets require payment proof", "")
	}

	if err := s.ticketRepo.Update(ctx, ticketID, bson.M{"payment_proof_url": proofURL}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, ticket.RaffleID, userID, "proof_uploaded", map[string]interface{}{
		"ticket_id": ticketID.Hex(),
		"proof_url": proofURL,
	}, "", "")

	ticket.PaymentProofURL = &proofURL
	resp := ticketToDTO(ticket)
	return &resp, nil
}

func (s *TicketService) Verify(ctx context.Context, ticketID, verifierID bson.ObjectID) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if ticket.ParticipantID == verifierID {
		return nil, utils.NewAPIError(403, "cannot verify your own ticket", "")
	}

	if ticket.Status != domain.TicketStatusPending {
		return nil, utils.NewAPIError(400, "ticket is not in pending status", "")
	}

	now := time.Now()
	if err := s.ticketRepo.Update(ctx, ticketID, bson.M{
		"status":      domain.TicketStatusVerified,
		"verified_by": verifierID,
		"verified_at": now,
	}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, ticket.RaffleID, verifierID, "ticket_verified", map[string]interface{}{
		"ticket_id":       ticketID.Hex(),
		"participant_id":  ticket.ParticipantID.Hex(),
	}, "", "")

	ticket.Status = domain.TicketStatusVerified
	ticket.VerifiedBy = &verifierID
	ticket.VerifiedAt = &now
	resp := ticketToDTO(ticket)
	return &resp, nil
}

func (s *TicketService) Reject(ctx context.Context, ticketID, verifierID bson.ObjectID, reason string) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if ticket.ParticipantID == verifierID {
		return nil, utils.NewAPIError(403, "cannot reject your own ticket", "")
	}

	if ticket.Status != domain.TicketStatusPending {
		return nil, utils.NewAPIError(400, "ticket is not in pending status", "")
	}

	if reason == "" {
		return nil, utils.NewAPIError(400, "rejection reason is required", "")
	}

	if err := s.ticketRepo.Update(ctx, ticketID, bson.M{
		"status":           domain.TicketStatusRejected,
		"verified_by":      verifierID,
		"rejection_reason": reason,
	}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, ticket.RaffleID, verifierID, "ticket_rejected", map[string]interface{}{
		"ticket_id":       ticketID.Hex(),
		"participant_id":  ticket.ParticipantID.Hex(),
		"rejection_reason": reason,
	}, "", "")

	ticket.Status = domain.TicketStatusRejected
	ticket.VerifiedBy = &verifierID
	ticket.RejectionReason = &reason
	resp := ticketToDTO(ticket)
	return &resp, nil
}

func getTicketIDs(tickets []*domain.Ticket) []string {
	ids := make([]string, len(tickets))
	for i, t := range tickets {
		ids[i] = t.ID.Hex()
	}
	return ids
}
