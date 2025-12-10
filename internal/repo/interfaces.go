/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package repo

import (
	"context"

	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CartRepository defines the interface for cart persistence
type CartRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (*model.Cart, error)
	Create(ctx context.Context, cart *model.Cart) error
	Update(ctx context.Context, cart *model.Cart) error
	DeleteItem(ctx context.Context, cartID bson.ObjectID, itemID string) error
}

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByStripeSessionID(ctx context.Context, stripeSessionID string) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderID bson.ObjectID, status string) error
}

// AddressRepository defines the interface for address persistence
type AddressRepository interface {
	Create(ctx context.Context, addr *model.Address) error
	ListByUserID(ctx context.Context, userID string) ([]model.Address, error)
	Get(ctx context.Context, id, userID string) (*model.Address, error)
	Update(ctx context.Context, addr *model.Address) error
	Delete(ctx context.Context, id, userID string) error
	UnsetDefaults(ctx context.Context, userID string) error
}
