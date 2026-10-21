package db

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func RunMigrations(ctx context.Context, db *mongo.Database) error {
	slog.Info("running database index migrations")

	if err := createUsersIndexes(ctx, db.Collection("users")); err != nil {
		return err
	}
	if err := createRafflesIndexes(ctx, db.Collection("raffles")); err != nil {
		return err
	}
	if err := createTicketsIndexes(ctx, db.Collection("tickets")); err != nil {
		return err
	}
	if err := createAuditLogsIndexes(ctx, db.Collection("audit_logs")); err != nil {
		return err
	}
	if err := createRefreshTokensIndexes(ctx, db.Collection("refresh_tokens")); err != nil {
		return err
	}

	slog.Info("database index migrations completed")
	return nil
}

func createUsersIndexes(ctx context.Context, col *mongo.Collection) error {
	emailIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, emailIdx); err != nil {
		return err
	}

	oauthIdx := mongo.IndexModel{
		Keys: bson.D{{Key: "oauth_provider", Value: 1}, {Key: "oauth_id", Value: 1}},
	}
	if _, err := col.Indexes().CreateOne(ctx, oauthIdx); err != nil {
		return err
	}

	return nil
}

func createRafflesIndexes(ctx context.Context, col *mongo.Collection) error {
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "raffler_id", Value: 1}, {Key: "status", Value: 1}},
	}
	if _, err := col.Indexes().CreateOne(ctx, idx); err != nil {
		return err
	}

	return nil
}

func createTicketsIndexes(ctx context.Context, col *mongo.Collection) error {
	uniqueIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "raffle_id", Value: 1}, {Key: "ticket_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, uniqueIdx); err != nil {
		return err
	}

	raffleStatusIdx := mongo.IndexModel{
		Keys: bson.D{{Key: "raffle_id", Value: 1}, {Key: "status", Value: 1}},
	}
	if _, err := col.Indexes().CreateOne(ctx, raffleStatusIdx); err != nil {
		return err
	}

	participantIdx := mongo.IndexModel{
		Keys: bson.D{{Key: "participant_id", Value: 1}},
	}
	if _, err := col.Indexes().CreateOne(ctx, participantIdx); err != nil {
		return err
	}

	return nil
}

func createAuditLogsIndexes(ctx context.Context, col *mongo.Collection) error {
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "raffle_id", Value: 1}, {Key: "created_at", Value: -1}},
	}
	if _, err := col.Indexes().CreateOne(ctx, idx); err != nil {
		return err
	}

	return nil
}

func createRefreshTokensIndexes(ctx context.Context, col *mongo.Collection) error {
	tokenIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "token_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, tokenIdx); err != nil {
		return err
	}

	expIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	if _, err := col.Indexes().CreateOne(ctx, expIdx); err != nil {
		return err
	}

	return nil
}
