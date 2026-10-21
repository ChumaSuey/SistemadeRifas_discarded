package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Email         string        `bson:"email" json:"email"`
	EmailVerified bool          `bson:"email_verified" json:"email_verified"`
	PasswordHash  *string       `bson:"password_hash,omitempty" json:"-"`
	OAuthProvider *string       `bson:"oauth_provider,omitempty" json:"oauth_provider,omitempty"`
	OAuthID       *string       `bson:"oauth_id,omitempty" json:"oauth_id,omitempty"`
	Role          string        `bson:"role" json:"role"`
	Status        string        `bson:"status" json:"status"`
	SocialLinks   SocialLinks   `bson:"social_links" json:"social_links"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}

type SocialLinks struct {
	Twitch    string `bson:"twitch,omitempty" json:"twitch,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
}

type Raffle struct {
	ID              bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	RafflerID       bson.ObjectID  `bson:"raffler_id" json:"raffler_id"`
	Title           string         `bson:"title" json:"title"`
	Description     string         `bson:"description" json:"description"`
	Status          string         `bson:"status" json:"status"`
	PricingTier     string         `bson:"pricing_tier" json:"pricing_tier"`
	TicketPrice     float64        `bson:"ticket_price" json:"ticket_price"`
	Currency        string         `bson:"currency" json:"currency"`
	PaymentMethods  []string       `bson:"payment_methods" json:"payment_methods"`
	TicketsTotal    int            `bson:"tickets_total" json:"tickets_total"`
	TicketsSold     int            `bson:"tickets_sold" json:"tickets_sold"`
	WinnerID        *bson.ObjectID `bson:"winner_id,omitempty" json:"winner_id,omitempty"`
	WinningTicketID *bson.ObjectID `bson:"winning_ticket_id,omitempty" json:"winning_ticket_id,omitempty"`
	SortSeed        *string        `bson:"sort_seed,omitempty" json:"sort_seed,omitempty"`
	SortCommitHash  *string        `bson:"sort_commit_hash,omitempty" json:"sort_commit_hash,omitempty"`
	DrawDate        *time.Time     `bson:"draw_date,omitempty" json:"draw_date,omitempty"`
	CreatedAt       time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `bson:"updated_at" json:"updated_at"`
}

type Ticket struct {
	ID                    bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	RaffleID              bson.ObjectID  `bson:"raffle_id" json:"raffle_id"`
	ParticipantID         bson.ObjectID  `bson:"participant_id" json:"participant_id"`
	TicketNumber          int            `bson:"ticket_number" json:"ticket_number"`
	Status                string         `bson:"status" json:"status"`
	PaymentMethod         string         `bson:"payment_method" json:"payment_method"`
	PaymentProofURL       *string        `bson:"payment_proof_url,omitempty" json:"payment_proof_url,omitempty"`
	StripePaymentIntentID *string        `bson:"stripe_payment_intent_id,omitempty" json:"stripe_payment_intent_id,omitempty"`
	RejectionReason       *string        `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	VerifiedBy            *bson.ObjectID `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	VerifiedAt            *time.Time     `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	CreatedAt             time.Time      `bson:"created_at" json:"created_at"`
}

type AuditLog struct {
	ID        bson.ObjectID          `bson:"_id,omitempty" json:"id"`
	RaffleID  bson.ObjectID          `bson:"raffle_id" json:"raffle_id"`
	UserID    bson.ObjectID          `bson:"user_id" json:"user_id"`
	Action    string                 `bson:"action" json:"action"`
	Details   map[string]interface{} `bson:"details" json:"details"`
	IPAddress string                 `bson:"ip_address" json:"ip_address"`
	UserAgent string                 `bson:"user_agent" json:"user_agent"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}

type RefreshToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	TokenHash string        `bson:"token_hash" json:"-"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

const (
	RoleAdmin       = "admin"
	RoleRaffler     = "raffler"
	RoleParticipant = "participant"
	RoleVerifier    = "verifier"

	StatusActive    = "active"
	StatusSuspended = "suspended"
	StatusBanned    = "banned"

	RaffleStatusDraft              = "draft"
	RaffleStatusPendingVerification = "pending_verification"
	RaffleStatusActive             = "active"
	RaffleStatusFinished           = "finished"
	RaffleStatusCancelled          = "cancelled"

	TicketStatusPending   = "pending"
	TicketStatusVerified  = "verified"
	TicketStatusRejected  = "rejected"
	TicketStatusCancelled = "cancelled"
	TicketStatusWinner    = "winner"

	PaymentMethodManual = "manual_transfer"
	PaymentMethodStripe = "stripe"

	CurrencyUSD = "USD"
	CurrencyVES = "VES"
	CurrencyEUR = "EUR"

	PricingTierFree = "free"
	PricingTierPro  = "pro"

	MaxFreeTickets = 100
)
