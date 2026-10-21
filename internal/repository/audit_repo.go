package repository

import (
	"context"
	"time"

	"raffle-saas/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuditRepository struct {
	col *mongo.Collection
}

func NewAuditRepository(db *mongo.Database) *AuditRepository {
	return &AuditRepository{col: db.Collection("audit_logs")}
}

func (r *AuditRepository) Log(ctx context.Context, raffleID, userID bson.ObjectID, action string, details map[string]interface{}, ip, userAgent string) error {
	entry := &domain.AuditLog{
		RaffleID:  raffleID,
		UserID:    userID,
		Action:    action,
		Details:   details,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	_, err := r.col.InsertOne(ctx, entry)
	return err
}
