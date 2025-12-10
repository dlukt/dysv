/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package repo

import (
	"context"
	"errors"
	"time"

	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Ensure OrderRepo implements OrderRepository
var _ OrderRepository = (*OrderRepo)(nil)

// OrderRepo is the MongoDB implementation of OrderRepository
type OrderRepo struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// NewOrderRepo creates a new order repository
func NewOrderRepo(db *mongo.Database, timeout time.Duration) *OrderRepo {
	return &OrderRepo{
		coll:    db.Collection("orders"),
		timeout: timeout,
	}
}

// Create inserts a new order
func (r *OrderRepo) Create(ctx context.Context, order *model.Order) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.coll.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	order.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

// FindByStripeSessionID finds an order by Stripe session ID
func (r *OrderRepo) FindByStripeSessionID(ctx context.Context, stripeSessionID string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var order model.Order
	err := r.coll.FindOne(ctx, bson.M{"stripe_session_id": stripeSessionID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

// UpdateStatus updates the status of an order
func (r *OrderRepo) UpdateStatus(ctx context.Context, orderID bson.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}
