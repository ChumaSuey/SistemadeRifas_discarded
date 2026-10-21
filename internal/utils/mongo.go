package utils

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	slog.Info("connecting to MongoDB", "uri", uri)
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	slog.Info("connected to MongoDB")
	return client, nil
}
