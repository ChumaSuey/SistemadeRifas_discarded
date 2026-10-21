package repository

import (
	"context"
	"time"

	"raffle-saas/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RaffleRepository struct {
	col *mongo.Collection
}

func NewRaffleRepository(db *mongo.Database) *RaffleRepository {
	return &RaffleRepository{col: db.Collection("raffles")}
}

func (r *RaffleRepository) Create(ctx context.Context, raffle *domain.Raffle) error {
	raffle.CreatedAt = time.Now()
	raffle.UpdatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, raffle)
	if err != nil {
		return err
	}
	raffle.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *RaffleRepository) FindByID(ctx context.Context, id bson.ObjectID) (*domain.Raffle, error) {
	var raffle domain.Raffle
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&raffle)
	if err != nil {
		return nil, err
	}
	return &raffle, nil
}

func (r *RaffleRepository) ListPublic(ctx context.Context, page, perPage int) ([]*domain.Raffle, int64, error) {
	filter := bson.M{"status": domain.RaffleStatusActive}
	return r.list(ctx, filter, page, perPage)
}

func (r *RaffleRepository) ListByRaffler(ctx context.Context, rafflerID bson.ObjectID, page, perPage int) ([]*domain.Raffle, int64, error) {
	filter := bson.M{"raffler_id": rafflerID}
	return r.list(ctx, filter, page, perPage)
}

func (r *RaffleRepository) list(ctx context.Context, filter bson.M, page, perPage int) ([]*domain.Raffle, int64, error) {
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * perPage)).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var raffles []*domain.Raffle
	if err := cursor.All(ctx, &raffles); err != nil {
		return nil, 0, err
	}
	return raffles, total, nil
}

func (r *RaffleRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	return err
}

func (r *RaffleRepository) IncrementTicketsSold(ctx context.Context, id bson.ObjectID, count int) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"tickets_sold": count},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}
