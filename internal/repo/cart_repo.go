/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Ensure CartRepo implements CartRepository
var _ CartRepository = (*CartRepo)(nil)

// CartRepo is the MongoDB implementation of CartRepository
type CartRepo struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// NewCartRepo creates a new cart repository
func NewCartRepo(db *mongo.Database, timeout time.Duration) *CartRepo {
	return &CartRepo{
		coll:    db.Collection("carts"),
		timeout: timeout,
	}
}

// FindBySessionID finds a cart by session ID
func (r *CartRepo) FindBySessionID(ctx context.Context, sessionID string) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var cart model.Cart
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&cart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		fmt.Printf("Repo: FindBySessionID error: %v\n", err)
		return nil, err
	}
	return &cart, nil
}

// Create inserts a new cart
func (r *CartRepo) Create(ctx context.Context, cart *model.Cart) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.coll.InsertOne(ctx, cart)
	if err != nil {
		fmt.Printf("Repo: Create error: %v\n", err)
		return err
	}
	cart.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

// Update replaces an existing cart
func (r *CartRepo) Update(ctx context.Context, cart *model.Cart) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	spew.Dump("Repo Update Cart", cart)
	// fmt.Printf("Repo: Updating cart %s\n", cart.ID.Hex())
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": cart.ID}, cart)
	if err != nil {
		fmt.Printf("Repo: ReplaceOne error: %v\n", err)
	}
	return err
}

// DeleteItem removes an item from a cart by item ID
func (r *CartRepo) DeleteItem(ctx context.Context, cartID bson.ObjectID, itemID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": cartID},
		bson.M{"$pull": bson.M{"items": bson.M{"item_id": itemID}}},
	)
	if err != nil {
		fmt.Printf("Repo: DeleteItem error: %v\n", err)
	}
	return err
}
