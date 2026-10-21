package repository

import (
	"context"
	"time"

	"raffle-saas/internal/domain"
	"raffle-saas/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TicketRepository struct {
	col  *mongo.Collection
	db   *mongo.Database
}

func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{col: db.Collection("tickets"), db: db}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	ticket.CreatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, ticket)
	if err != nil {
		return err
	}
	ticket.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *TicketRepository) CreateInTransaction(ctx context.Context, session *mongo.Session, ticket *domain.Ticket, raffleID bson.ObjectID) error {
	ticket.CreatedAt = time.Now()

	raffleCol := r.db.Collection("raffles")

	_, err := session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		var raffle domain.Raffle
		err := raffleCol.FindOne(sc, bson.M{"_id": raffleID}).Decode(&raffle)
		if err != nil {
			return nil, err
		}

		if raffle.TicketsSold >= raffle.TicketsTotal {
			return nil, utils.ErrConflict
		}

		ticket.TicketNumber = raffle.TicketsSold + 1

		result, err := r.col.InsertOne(sc, ticket)
		if mongo.IsDuplicateKeyError(err) {
			return nil, utils.ErrConflict
		}
		if err != nil {
			return nil, err
		}
		ticket.ID = result.InsertedID.(bson.ObjectID)

		_, err = raffleCol.UpdateOne(sc,
			bson.M{"_id": raffleID, "tickets_sold": bson.M{"$lt": raffle.TicketsTotal}},
			bson.M{
				"$inc": bson.M{"tickets_sold": 1},
				"$set": bson.M{"updated_at": time.Now()},
			},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (r *TicketRepository) FindByID(ctx context.Context, id bson.ObjectID) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) FindByRaffleAndNumber(ctx context.Context, raffleID bson.ObjectID, number int) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.col.FindOne(ctx, bson.M{"raffle_id": raffleID, "ticket_number": number}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) FindVerifiedByRaffle(ctx context.Context, raffleID bson.ObjectID) ([]*domain.Ticket, error) {
	filter := bson.M{"raffle_id": raffleID, "status": domain.TicketStatusVerified}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []*domain.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) ListByRaffle(ctx context.Context, raffleID bson.ObjectID, page, perPage int) ([]*domain.Ticket, int64, error) {
	filter := bson.M{"raffle_id": raffleID}
	return r.list(ctx, filter, page, perPage)
}

func (r *TicketRepository) ListByParticipant(ctx context.Context, participantID bson.ObjectID, page, perPage int) ([]*domain.Ticket, int64, error) {
	filter := bson.M{"participant_id": participantID}
	return r.list(ctx, filter, page, perPage)
}

func (r *TicketRepository) list(ctx context.Context, filter bson.M, page, perPage int) ([]*domain.Ticket, int64, error) {
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

	var tickets []*domain.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *TicketRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	return err
}

func (r *TicketRepository) BulkUpdateStatus(ctx context.Context, raffleID bson.ObjectID, fromStatus, toStatus string) error {
	_, err := r.col.UpdateMany(ctx,
		bson.M{"raffle_id": raffleID, "status": fromStatus},
		bson.M{"$set": bson.M{"status": toStatus}},
	)
	return err
}
