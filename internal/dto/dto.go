package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Email         string            `json:"email"`
	EmailVerified bool              `json:"email_verified"`
	Role          string            `json:"role"`
	Status        string            `json:"status"`
	SocialLinks   map[string]string `json:"social_links,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

type CreateRaffleRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	TicketPrice    float64  `json:"ticket_price"`
	Currency       string   `json:"currency"`
	PaymentMethods []string `json:"payment_methods"`
	TicketsTotal   int      `json:"tickets_total"`
}

type UpdateRaffleRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type RaffleResponse struct {
	ID              string                 `json:"id"`
	RafflerID       string                 `json:"raffler_id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Status          string                 `json:"status"`
	PricingTier     string                 `json:"pricing_tier"`
	TicketPrice     float64                `json:"ticket_price"`
	Currency        string                 `json:"currency"`
	PaymentMethods  []string               `json:"payment_methods"`
	TicketsTotal    int                    `json:"tickets_total"`
	TicketsSold     int                    `json:"tickets_sold"`
	WinnerID        *string                `json:"winner_id,omitempty"`
	WinningTicketID *string                `json:"winning_ticket_id,omitempty"`
	SortCommitHash  *string                `json:"sort_commit_hash,omitempty"`
	DrawDate        *time.Time             `json:"draw_date,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type RaffleDetailResponse struct {
	RaffleResponse
	Tickets []TicketResponse `json:"tickets,omitempty"`
}

type PurchaseTicketsRequest struct {
	Quantity      int    `json:"quantity"`
	PaymentMethod string `json:"payment_method"`
}

type TicketResponse struct {
	ID                    string     `json:"id"`
	RaffleID              string     `json:"raffle_id"`
	ParticipantID         string     `json:"participant_id"`
	TicketNumber          int        `json:"ticket_number"`
	Status                string     `json:"status"`
	PaymentMethod         string     `json:"payment_method"`
	PaymentProofURL       *string    `json:"payment_proof_url,omitempty"`
	StripePaymentIntentID *string    `json:"stripe_payment_intent_id,omitempty"`
	RejectionReason       *string    `json:"rejection_reason,omitempty"`
	VerifiedBy            *string    `json:"verified_by,omitempty"`
	VerifiedAt            *time.Time `json:"verified_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

type UploadProofRequest struct {
	PaymentProofURL string `json:"payment_proof_url"`
}

type RejectTicketRequest struct {
	Reason string `json:"reason"`
}

type SortCommitResponse struct {
	CommitHash     string `json:"commit_hash"`
	TicketsCount   int    `json:"tickets_count"`
	Timestamp      int64  `json:"timestamp"`
}

type SortRevealResponse struct {
	CommitHash     string `json:"commit_hash"`
	Seed           string `json:"seed"`
	WinnerTicketID string `json:"winning_ticket_id"`
	WinnerUserID   string `json:"winner_user_id"`
	TicketNumber   int    `json:"ticket_number"`
	DrawDate       string `json:"draw_date"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func ObjectIDToString(id bson.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}

func ObjectIDPtrToString(id *bson.ObjectID) *string {
	if id == nil || id.IsZero() {
		return nil
	}
	s := id.Hex()
	return &s
}
