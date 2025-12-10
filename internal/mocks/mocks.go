/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package mocks

import (
	"context"
	"time"

	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockCartRepo is a mock implementation of cart repository for testing
type MockCartRepo struct {
	Carts       map[string]*model.Cart
	InsertError error
	FindError   error
	UpdateError error
}

// NewMockCartRepo creates a new mock cart repository
func NewMockCartRepo() *MockCartRepo {
	return &MockCartRepo{
		Carts: make(map[string]*model.Cart),
	}
}

// FindBySessionID finds a cart by session ID
func (m *MockCartRepo) FindBySessionID(ctx context.Context, sessionID string) (*model.Cart, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if cart, ok := m.Carts[sessionID]; ok {
		return cart, nil
	}
	return nil, nil
}

// Insert inserts a new cart
func (m *MockCartRepo) Insert(ctx context.Context, cart *model.Cart) error {
	if m.InsertError != nil {
		return m.InsertError
	}
	cart.ID = bson.NewObjectID()
	m.Carts[cart.SessionID] = cart
	return nil
}

// Update updates an existing cart
func (m *MockCartRepo) Update(ctx context.Context, cart *model.Cart) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	cart.UpdatedAt = time.Now()
	m.Carts[cart.SessionID] = cart
	return nil
}

// Reset clears all mock data
func (m *MockCartRepo) Reset() {
	m.Carts = make(map[string]*model.Cart)
	m.InsertError = nil
	m.FindError = nil
	m.UpdateError = nil
}

// MockOrderRepo is a mock implementation of order repository for testing
type MockOrderRepo struct {
	Orders      map[string]*model.Order
	InsertError error
	FindError   error
	UpdateError error
}

// NewMockOrderRepo creates a new mock order repository
func NewMockOrderRepo() *MockOrderRepo {
	return &MockOrderRepo{
		Orders: make(map[string]*model.Order),
	}
}

// FindByStripeSessionID finds an order by Stripe session ID
func (m *MockOrderRepo) FindByStripeSessionID(ctx context.Context, stripeSessionID string) (*model.Order, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if order, ok := m.Orders[stripeSessionID]; ok {
		return order, nil
	}
	return nil, nil
}

// Insert inserts a new order
func (m *MockOrderRepo) Insert(ctx context.Context, order *model.Order) error {
	if m.InsertError != nil {
		return m.InsertError
	}
	order.ID = bson.NewObjectID()
	m.Orders[order.StripeSessionID] = order
	return nil
}

// UpdateStatus updates order status
func (m *MockOrderRepo) UpdateStatus(ctx context.Context, orderID bson.ObjectID, status string) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	for _, order := range m.Orders {
		if order.ID == orderID {
			order.Status = status
			return nil
		}
	}
	return nil
}

// Reset clears all mock data
func (m *MockOrderRepo) Reset() {
	m.Orders = make(map[string]*model.Order)
	m.InsertError = nil
	m.FindError = nil
	m.UpdateError = nil
}
