/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package mocks

import (
	"context"
	"time"

	"github.com/deicod/auth/core"
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

// MockAddressRepo is a mock implementation of address repository for testing
type MockAddressRepo struct {
	Addresses         map[string]*model.Address // Keyed by ID for simplicity, or we can scan
	CreateError       error
	ListError         error
	GetError          error
	UpdateError       error
	DeleteError       error
	UnsetDefaultError error
}

func NewMockAddressRepo() *MockAddressRepo {
	return &MockAddressRepo{
		Addresses: make(map[string]*model.Address),
	}
}

func (m *MockAddressRepo) Create(ctx context.Context, addr *model.Address) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if addr.ID == "" {
		addr.ID = bson.NewObjectID().Hex()
	}
	// Store a copy to avoid pointer aliasing issues in tests
	val := *addr
	m.Addresses[addr.ID] = &val
	return nil
}

func (m *MockAddressRepo) ListByUserID(ctx context.Context, userID string) ([]model.Address, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	var result []model.Address
	for _, addr := range m.Addresses {
		if addr.UserID == userID {
			result = append(result, *addr)
		}
	}
	return result, nil
}

func (m *MockAddressRepo) Get(ctx context.Context, id, userID string) (*model.Address, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	addr, ok := m.Addresses[id]
	if !ok || addr.UserID != userID {
		return nil, nil // Not found
	}
	// Return copy
	val := *addr
	return &val, nil
}

func (m *MockAddressRepo) Update(ctx context.Context, addr *model.Address) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	// Store a copy
	val := *addr
	m.Addresses[addr.ID] = &val
	return nil
}

func (m *MockAddressRepo) Delete(ctx context.Context, id, userID string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if addr, ok := m.Addresses[id]; ok && addr.UserID == userID {
		delete(m.Addresses, id)
	}
	return nil
}

func (m *MockAddressRepo) UnsetDefaults(ctx context.Context, userID string) error {
	if m.UnsetDefaultError != nil {
		return m.UnsetDefaultError
	}
	for _, addr := range m.Addresses {
		if addr.UserID == userID {
			addr.IsDefault = false
		}
	}
	return nil
}

// MockAuthService is a partial mock of auth.Service needed for handlers
type MockAuthService struct {
	AuthenticateSessionFunc func(ctx context.Context, token string) (core.UserPublic, core.SessionPublic, error)
}

func (m *MockAuthService) AuthenticateSession(ctx context.Context, token string) (core.UserPublic, core.SessionPublic, error) {
	if m.AuthenticateSessionFunc != nil {
		return m.AuthenticateSessionFunc(ctx, token)
	}
	return core.UserPublic{}, core.SessionPublic{}, nil
}

func (m *MockAuthService) Register(ctx context.Context, cmd core.RegisterCommand) (core.AuthResult, error) {
	return core.AuthResult{}, nil
}

func (m *MockAuthService) Login(ctx context.Context, cmd core.LoginCommand) (core.AuthResult, error) {
	return core.AuthResult{}, nil
}

func (m *MockAuthService) VerifyEmail(ctx context.Context, cmd core.VerifyEmailCommand) (core.VerifyEmailResult, error) {
	return core.VerifyEmailResult{}, nil
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, cmd core.ForgotPasswordCommand) error {
	return nil
}

func (m *MockAuthService) ResetPassword(ctx context.Context, cmd core.ResetPasswordCommand) (core.UserPublic, error) {
	return core.UserPublic{}, nil
}

func (m *MockAuthService) InitiateEmailChange(ctx context.Context, cmd core.ChangeEmailCommand) error {
	return nil
}

func (m *MockAuthService) ConfirmEmailChange(ctx context.Context, cmd core.ConfirmEmailChangeCommand) (core.ChangeEmailResult, error) {
	return core.ChangeEmailResult{}, nil
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return nil
}

// Implement other methods of auth.Service if required by interface,
// strictly we only need what AddressHandler uses.
// However, if the interface is large, we might need to implement stub methods.
// Let's assume for now we only need AuthenticateSession or we'll get a compile error if we pass it
// where the full interface is expected.
// AddressHandler uses `auth.Service`. We need to know what `auth.Service` looks like.
// It is likely an interface. If it has more methods, we need to stub them.
