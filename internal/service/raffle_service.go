package service

import (
	"context"
	"time"

	"raffle-saas/internal/domain"
	"raffle-saas/internal/dto"
	"raffle-saas/internal/repository"
	"raffle-saas/internal/utils"
	"raffle-saas/pkg/crypto"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RaffleService struct {
	raffleRepo *repository.RaffleRepository
	ticketRepo *repository.TicketRepository
	auditRepo  *repository.AuditRepository
}

func NewRaffleService(raffleRepo *repository.RaffleRepository, ticketRepo *repository.TicketRepository, auditRepo *repository.AuditRepository) *RaffleService {
	return &RaffleService{
		raffleRepo: raffleRepo,
		ticketRepo: ticketRepo,
		auditRepo:  auditRepo,
	}
}

func (s *RaffleService) Create(ctx context.Context, rafflerID bson.ObjectID, req *dto.CreateRaffleRequest) (*dto.RaffleResponse, error) {
	if req.Title == "" {
		return nil, utils.NewAPIError(400, "title is required", "")
	}
	if req.TicketsTotal <= 0 {
		return nil, utils.NewAPIError(400, "tickets_total must be greater than 0", "")
	}
	if req.TicketPrice < 0 {
		return nil, utils.NewAPIError(400, "ticket_price cannot be negative", "")
	}

	raffle := &domain.Raffle{
		RafflerID:      rafflerID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         domain.RaffleStatusDraft,
		PricingTier:    domain.PricingTierFree,
		TicketPrice:    req.TicketPrice,
		Currency:       req.Currency,
		PaymentMethods: req.PaymentMethods,
		TicketsTotal:   req.TicketsTotal,
	}

	if raffle.Currency == "" {
		raffle.Currency = domain.CurrencyUSD
	}
	if len(raffle.PaymentMethods) == 0 {
		raffle.PaymentMethods = []string{domain.PaymentMethodManual}
	}

	if err := s.raffleRepo.Create(ctx, raffle); err != nil {
		return nil, err
	}

	return raffleToDTO(raffle), nil
}

func (s *RaffleService) GetByID(ctx context.Context, id bson.ObjectID) (*dto.RaffleDetailResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	resp := raffleToDTO(raffle)
	tickets, _, err := s.ticketRepo.ListByRaffle(ctx, id, 1, 100)
	if err != nil {
		return nil, err
	}

	ticketDTOs := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		ticketDTOs[i] = ticketToDTO(t)
	}

	return &dto.RaffleDetailResponse{
		RaffleResponse: *resp,
		Tickets:        ticketDTOs,
	}, nil
}

func (s *RaffleService) ListPublic(ctx context.Context, page, perPage int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	raffles, total, err := s.raffleRepo.ListPublic(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.RaffleResponse, len(raffles))
	for i, r := range raffles {
		dtos[i] = *raffleToDTO(r)
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

func (s *RaffleService) ListMy(ctx context.Context, rafflerID bson.ObjectID, page, perPage int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	raffles, total, err := s.raffleRepo.ListByRaffler(ctx, rafflerID, page, perPage)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.RaffleResponse, len(raffles))
	for i, r := range raffles {
		dtos[i] = *raffleToDTO(r)
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

func (s *RaffleService) Update(ctx context.Context, raffleID, userID bson.ObjectID, req *dto.UpdateRaffleRequest) (*dto.RaffleResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status != domain.RaffleStatusDraft {
		return nil, utils.NewAPIError(400, "raffle can only be updated in draft status", "")
	}

	update := bson.M{}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}

	if len(update) > 0 {
		if err := s.raffleRepo.Update(ctx, raffleID, update); err != nil {
			return nil, err
		}
		raffle, _ = s.raffleRepo.FindByID(ctx, raffleID)
	}

	return raffleToDTO(raffle), nil
}

func (s *RaffleService) Publish(ctx context.Context, raffleID, userID bson.ObjectID) (*dto.RaffleResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status != domain.RaffleStatusDraft {
		return nil, utils.NewAPIError(400, "only draft raffles can be published", "")
	}

	if raffle.PricingTier == domain.PricingTierFree && raffle.TicketsTotal > domain.MaxFreeTickets {
		return nil, utils.NewAPIError(400, "free tier raffles are limited to 100 tickets", "")
	}

	if err := s.raffleRepo.Update(ctx, raffleID, bson.M{"status": domain.RaffleStatusActive}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, raffleID, userID, "raffle_published", map[string]interface{}{
		"previous_status": raffle.Status,
	}, "", "")

	raffle.Status = domain.RaffleStatusActive
	return raffleToDTO(raffle), nil
}

func (s *RaffleService) Cancel(ctx context.Context, raffleID, userID bson.ObjectID) (*dto.RaffleResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status == domain.RaffleStatusFinished {
		return nil, utils.NewAPIError(400, "finished raffles cannot be cancelled", "")
	}

	if err := s.raffleRepo.Update(ctx, raffleID, bson.M{"status": domain.RaffleStatusCancelled}); err != nil {
		return nil, err
	}

	if raffle.Status == domain.RaffleStatusActive {
		s.ticketRepo.BulkUpdateStatus(ctx, raffleID, domain.TicketStatusPending, domain.TicketStatusCancelled)
		s.ticketRepo.BulkUpdateStatus(ctx, raffleID, domain.TicketStatusVerified, domain.TicketStatusCancelled)
	}

	s.auditRepo.Log(ctx, raffleID, userID, "raffle_cancelled", map[string]interface{}{
		"previous_status": raffle.Status,
	}, "", "")

	raffle.Status = domain.RaffleStatusCancelled
	return raffleToDTO(raffle), nil
}

func (s *RaffleService) CommitSort(ctx context.Context, raffleID, userID bson.ObjectID) (*dto.SortCommitResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status != domain.RaffleStatusActive {
		return nil, utils.NewAPIError(400, "sort can only be started on active raffles", "")
	}

	seedInput, commitHash, err := crypto.GenerateSortSeed(raffleID.Hex())
	if err != nil {
		return nil, err
	}

	if err := s.raffleRepo.Update(ctx, raffleID, bson.M{
		"sort_seed":        seedInput,
		"sort_commit_hash": commitHash,
	}); err != nil {
		return nil, err
	}

	tickets, _ := s.ticketRepo.FindVerifiedByRaffle(ctx, raffleID)

	s.auditRepo.Log(ctx, raffleID, userID, "sort_committed", map[string]interface{}{
		"commit_hash":   commitHash,
		"tickets_count": len(tickets),
	}, "", "")

	return &dto.SortCommitResponse{
		CommitHash:   commitHash,
		TicketsCount: len(tickets),
		Timestamp:    time.Now().Unix(),
	}, nil
}

func (s *RaffleService) ExecuteSort(ctx context.Context, raffleID, userID bson.ObjectID) (*dto.SortRevealResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status != domain.RaffleStatusActive {
		return nil, utils.NewAPIError(400, "sort can only be executed on active raffles", "")
	}

	if raffle.SortSeed == nil || raffle.SortCommitHash == nil {
		return nil, utils.NewAPIError(400, "commit phase must be completed first", "")
	}

	verifiedHash := crypto.ComputeCommitHash(*raffle.SortSeed)
	if verifiedHash != *raffle.SortCommitHash {
		return nil, utils.NewAPIError(400, "seed does not match commit hash", "")
	}

	tickets, err := s.ticketRepo.FindVerifiedByRaffle(ctx, raffleID)
	if err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, utils.ErrUnprocessableEntity
	}

	numbers := make([]int, len(tickets))
	for i, t := range tickets {
		numbers[i] = t.TicketNumber
	}

	crypto.FisherYatesShuffle(*raffle.SortSeed, numbers)

	winningNumber := numbers[0]
	var winningTicket *domain.Ticket
	for _, t := range tickets {
		if t.TicketNumber == winningNumber {
			winningTicket = t
			break
		}
	}

	now := time.Now()
	if err := s.raffleRepo.Update(ctx, raffleID, bson.M{
		"status":             domain.RaffleStatusFinished,
		"winner_id":          winningTicket.ParticipantID,
		"winning_ticket_id":  winningTicket.ID,
		"draw_date":          now,
	}); err != nil {
		return nil, err
	}

	if err := s.ticketRepo.Update(ctx, winningTicket.ID, bson.M{"status": domain.TicketStatusWinner}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, raffleID, userID, "sort_executed", map[string]interface{}{
		"seed":                *raffle.SortSeed,
		"commit_hash":         *raffle.SortCommitHash,
		"winning_ticket_id":   winningTicket.ID.Hex(),
		"winning_number":      winningNumber,
		"tickets_participant": len(tickets),
	}, "", "")

	return &dto.SortRevealResponse{
		CommitHash:     *raffle.SortCommitHash,
		Seed:           *raffle.SortSeed,
		WinnerTicketID: winningTicket.ID.Hex(),
		WinnerUserID:   winningTicket.ParticipantID.Hex(),
		TicketNumber:   winningNumber,
		DrawDate:       now.Format(time.RFC3339),
	}, nil
}

func (s *RaffleService) ReSort(ctx context.Context, raffleID, userID bson.ObjectID, disqualifyReason string) (*dto.SortRevealResponse, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	if raffle.RafflerID != userID {
		return nil, utils.ErrForbidden
	}

	if raffle.Status != domain.RaffleStatusFinished {
		return nil, utils.NewAPIError(400, "can only re-sort finished raffles", "")
	}

	if raffle.WinningTicketID != nil {
		if err := s.ticketRepo.Update(ctx, *raffle.WinningTicketID, bson.M{"status": domain.TicketStatusCancelled}); err != nil {
			return nil, err
		}
	}

	seedInput, commitHash, err := crypto.GenerateSortSeed(raffleID.Hex())
	if err != nil {
		return nil, err
	}

	if err := s.raffleRepo.Update(ctx, raffleID, bson.M{
		"status":           domain.RaffleStatusActive,
		"sort_seed":        seedInput,
		"sort_commit_hash": commitHash,
		"winner_id":        nil,
		"winning_ticket_id": nil,
		"draw_date":        nil,
	}); err != nil {
		return nil, err
	}

	s.auditRepo.Log(ctx, raffleID, userID, "resort_initiated", map[string]interface{}{
		"previous_winner_ticket_id": raffle.WinningTicketID.Hex(),
		"disqualify_reason":         disqualifyReason,
	}, "", "")

	return s.ExecuteSort(ctx, raffleID, userID)
}

func raffleToDTO(raffle *domain.Raffle) *dto.RaffleResponse {
	return &dto.RaffleResponse{
		ID:              raffle.ID.Hex(),
		RafflerID:       raffle.RafflerID.Hex(),
		Title:           raffle.Title,
		Description:     raffle.Description,
		Status:          raffle.Status,
		PricingTier:     raffle.PricingTier,
		TicketPrice:     raffle.TicketPrice,
		Currency:        raffle.Currency,
		PaymentMethods:  raffle.PaymentMethods,
		TicketsTotal:    raffle.TicketsTotal,
		TicketsSold:     raffle.TicketsSold,
		WinnerID:        dto.ObjectIDPtrToString(raffle.WinnerID),
		WinningTicketID: dto.ObjectIDPtrToString(raffle.WinningTicketID),
		SortCommitHash:  raffle.SortCommitHash,
		DrawDate:        raffle.DrawDate,
		CreatedAt:       raffle.CreatedAt,
		UpdatedAt:       raffle.UpdatedAt,
	}
}

func ticketToDTO(ticket *domain.Ticket) dto.TicketResponse {
	return dto.TicketResponse{
		ID:                    ticket.ID.Hex(),
		RaffleID:              ticket.RaffleID.Hex(),
		ParticipantID:         ticket.ParticipantID.Hex(),
		TicketNumber:          ticket.TicketNumber,
		Status:                ticket.Status,
		PaymentMethod:         ticket.PaymentMethod,
		PaymentProofURL:       ticket.PaymentProofURL,
		StripePaymentIntentID: ticket.StripePaymentIntentID,
		RejectionReason:       ticket.RejectionReason,
		VerifiedBy:            dto.ObjectIDPtrToString(ticket.VerifiedBy),
		VerifiedAt:            ticket.VerifiedAt,
		CreatedAt:             ticket.CreatedAt,
	}
}
