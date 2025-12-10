/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/repo"
)

// Pricing constants (must match frontend pricing-data.ts)
var Plans = map[string]model.Plan{
	"static-micro": {
		ID:             "static-micro",
		Name:           "Static Micro",
		MonthlyPrice:   3.90,
		TargetAudience: "React/Vue SPAs",
		Limits:         "Shared RAM, 1GB Storage",
	},
	"node-starter": {
		ID:             "node-starter",
		Name:           "Node Starter",
		MonthlyPrice:   9.90,
		TargetAudience: "Personal Blogs",
		Limits:         "1 vCPU (Shared), 512MB RAM, 5GB Storage",
	},
	"node-pro": {
		ID:             "node-pro",
		Name:           "Node Pro",
		MonthlyPrice:   39.90,
		TargetAudience: "E-commerce/SaaS",
		Limits:         "2 vCPU (Dedicated), 4GB RAM, 20GB Storage",
	},
}

var Addons = map[string]model.Addon{
	"de-domain": {
		ID:           "de-domain",
		Name:         ".de Domain",
		MonthlyPrice: 1.00,
	},
}

const YearlyDiscountMonths = 2

// CartService handles cart business logic
type CartService struct {
	cartRepo repo.CartRepository
}

// NewCartService creates a new cart service
func NewCartService(cartRepo repo.CartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

// GetOrCreateCart gets existing cart or creates a new one
func (s *CartService) GetOrCreateCart(ctx context.Context, sessionID string) (*model.Cart, error) {
	cart, err := s.cartRepo.FindBySessionID(ctx, sessionID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, repo.ErrNotFound) {
		fmt.Printf("Service: GetOrCreateCart Find Error: %v\n", err)
		return nil, err
	}

	// Create new cart
	cart = &model.Cart{
		SessionID:    sessionID,
		Items:        []model.LineItem{},
		BillingCycle: model.BillingMonthly,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.cartRepo.Create(ctx, cart); err != nil {
		fmt.Printf("Service: GetOrCreateCart Create Error: %v\n", err)
		return nil, err
	}
	return cart, nil
}

// AddPlan adds a plan to the cart (increments quantity if exists)
func (s *CartService) AddPlan(ctx context.Context, sessionID, planID string, quantity int) (*model.Cart, error) {
	plan, ok := Plans[planID]
	if !ok {
		return nil, ErrInvalidPlan
	}
	if quantity < 1 {
		quantity = 1
	}

	cart, err := s.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if plan already exists
	found := false
	for i, item := range cart.Items {
		if item.ItemType == "plan" && item.ItemID == planID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, model.LineItem{
			ItemID:   plan.ID,
			ItemType: "plan",
			Name:     plan.Name,
			Price:    plan.MonthlyPrice,
			Quantity: quantity,
		})
	}

	cart.UpdatedAt = time.Now()

	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// UpdateItemQuantity updates the quantity of an item
func (s *CartService) UpdateItemQuantity(ctx context.Context, sessionID, itemID string, quantity int) (*model.Cart, error) {
	cart, err := s.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if quantity <= 0 {
		return s.RemoveItem(ctx, sessionID, itemID)
	}

	found := false
	for i, item := range cart.Items {
		if item.ItemID == itemID {
			cart.Items[i].Quantity = quantity
			found = true
			break
		}
	}

	if !found {
		return cart, nil // Item not found, do nothing or return error? Current logic: idempotent success
	}

	cart.UpdatedAt = time.Now()

	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// AddAddon adds an addon to the cart
func (s *CartService) AddAddon(ctx context.Context, sessionID, addonID string) (*model.Cart, error) {
	addon, ok := Addons[addonID]
	if !ok {
		return nil, ErrInvalidAddon
	}

	cart, err := s.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if addon already exists
	for _, item := range cart.Items {
		if item.ItemID == addonID && item.ItemType == "addon" {
			return cart, nil // Already added
		}
	}

	cart.Items = append(cart.Items, model.LineItem{
		ItemID:   addon.ID,
		ItemType: "addon",
		Name:     addon.Name,
		Price:    addon.MonthlyPrice,
		Quantity: 1,
	})
	cart.UpdatedAt = time.Now()

	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(ctx context.Context, sessionID, itemID string) (*model.Cart, error) {
	cart, err := s.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var newItems []model.LineItem
	for _, item := range cart.Items {
		if item.ItemID != itemID {
			newItems = append(newItems, item)
		}
	}

	cart.Items = newItems
	cart.UpdatedAt = time.Now()

	if err := s.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// SetBillingCycle sets the billing cycle
func (s *CartService) SetBillingCycle(ctx context.Context, sessionID string, cycle model.BillingCycle) (*model.Cart, error) {
	cart, err := s.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	cart.BillingCycle = cycle
	cart.UpdatedAt = time.Now()

	fmt.Printf("SetBillingCycle: session=%s cycle=%s\n", sessionID, cycle)
	spew.Dump("Service Cart Before Update", cart)

	if err := s.cartRepo.Update(ctx, cart); err != nil {
		fmt.Printf("SetBillingCycle error: %v\n", err)
		return nil, err
	}
	return cart, nil
}

// GetCartTotal calculates the cart total
// For yearly: plans get 2 months free (×10), addons pay full 12 months
func (s *CartService) GetCartTotal(cart *model.Cart) (monthly float64, yearly float64) {
	var planMonthly, addonMonthly float64

	for _, item := range cart.Items {
		itemTotal := item.Price * float64(item.Quantity)
		monthly += itemTotal
		if item.ItemType == "plan" {
			planMonthly += itemTotal
		} else {
			addonMonthly += itemTotal
		}
	}

	// Yearly: plans get discount (10 months), addons pay full (12 months)
	yearly = (planMonthly * float64(12-YearlyDiscountMonths)) + (addonMonthly * 12)
	return
}
