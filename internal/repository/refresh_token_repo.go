package repository

import (
	"context"
	"time"

	"raffle-saas/internal/domain"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RefreshTokenRepository struct {
	col *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{col: db.Collection("refresh_tokens")}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID bson.ObjectID, token string, expiresAt time.Time) error {
	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: utils.HashToken(token),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	_, err := r.col.InsertOne(ctx, rt)
	return err
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.col.FindOne(ctx, bson.M{"token_hash": utils.HashToken(token)}).Decode(&rt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"token_hash": utils.HashToken(token)})
	return err
}

func (r *RefreshTokenRepository) DeleteAllForUser(ctx context.Context, userID bson.ObjectID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
