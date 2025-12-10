/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package repo

import (
	"context"
	"sync"

	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Ensure MockCartRepo implements CartRepository
var _ CartRepository = (*MockCartRepo)(nil)

// MockCartRepo is an in-memory implementation for testing
type MockCartRepo struct {
	mu    sync.RWMutex
	carts map[string]*model.Cart
}

// NewMockCartRepo creates a new mock cart repository
func NewMockCartRepo() *MockCartRepo {
	return &MockCartRepo{
		carts: make(map[string]*model.Cart),
	}
}

func (m *MockCartRepo) FindBySessionID(ctx context.Context, sessionID string) (*model.Cart, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cart, ok := m.carts[sessionID]
	if !ok {
		return nil, ErrNotFound
	}
	return cart, nil
}

func (m *MockCartRepo) Create(ctx context.Context, cart *model.Cart) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cart.ID = bson.NewObjectID()
	m.carts[cart.SessionID] = cart
	return nil
}

func (m *MockCartRepo) Update(ctx context.Context, cart *model.Cart) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.carts[cart.SessionID] = cart
	return nil
}

func (m *MockCartRepo) DeleteItem(ctx context.Context, cartID bson.ObjectID, itemID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, cart := range m.carts {
		if cart.ID == cartID {
			var newItems []model.LineItem
			for _, item := range cart.Items {
				if item.ItemID != itemID {
					newItems = append(newItems, item)
				}
			}
			cart.Items = newItems
			break
		}
	}
	return nil
}

// Reset clears all data (for test cleanup)
func (m *MockCartRepo) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.carts = make(map[string]*model.Cart)
}

// Ensure MockOrderRepo implements OrderRepository
var _ OrderRepository = (*MockOrderRepo)(nil)

// MockOrderRepo is an in-memory implementation for testing
type MockOrderRepo struct {
	mu     sync.RWMutex
	orders map[string]*model.Order // keyed by StripeSessionID
}

// NewMockOrderRepo creates a new mock order repository
func NewMockOrderRepo() *MockOrderRepo {
	return &MockOrderRepo{
		orders: make(map[string]*model.Order),
	}
}

func (m *MockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	order.ID = bson.NewObjectID()
	m.orders[order.StripeSessionID] = order
	return nil
}

func (m *MockOrderRepo) FindByStripeSessionID(ctx context.Context, stripeSessionID string) (*model.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	order, ok := m.orders[stripeSessionID]
	if !ok {
		return nil, ErrNotFound
	}
	return order, nil
}

func (m *MockOrderRepo) UpdateStatus(ctx context.Context, orderID bson.ObjectID, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, order := range m.orders {
		if order.ID == orderID {
			order.Status = status
			break
		}
	}
	return nil
}

// Reset clears all data (for test cleanup)
func (m *MockOrderRepo) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders = make(map[string]*model.Order)
}
