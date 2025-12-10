/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package service

import (
	"context"
	"fmt"
	"math"

	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/repo"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

// CheckoutService handles Stripe checkout
type CheckoutService struct {
	cartService *CartService
	orderRepo   repo.OrderRepository
	successURL  string
	cancelURL   string
}

// NewCheckoutService creates a new checkout service
func NewCheckoutService(cartService *CartService, orderRepo repo.OrderRepository, stripeKey, successURL, cancelURL string) *CheckoutService {
	stripe.Key = stripeKey
	return &CheckoutService{
		cartService: cartService,
		orderRepo:   orderRepo,
		successURL:  successURL,
		cancelURL:   cancelURL,
	}
}

// CreateCheckoutSession creates a Stripe Checkout session for the cart
func (s *CheckoutService) CreateCheckoutSession(ctx context.Context, sessionID string) (string, error) {
	cart, err := s.cartService.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		return "", err
	}

	if len(cart.Items) == 0 {
		return "", ErrEmptyCart
	}

	// Build line items for Stripe
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(cart.Items))
	var totalCents int64

	for _, item := range cart.Items {
		var unitAmount int64
		var interval stripe.PriceRecurringInterval

		if cart.BillingCycle == model.BillingYearly {
			interval = stripe.PriceRecurringIntervalYear
			if item.ItemType == "plan" {
				// Plans: charge 10 months worth (2 months free) once per year
				yearlyPrice := item.Price * float64(12-YearlyDiscountMonths)
				unitAmount = int64(math.Round(yearlyPrice * 100))
			} else {
				// Addons: no discount, pay full 12 months
				yearlyPrice := item.Price * 12
				unitAmount = int64(math.Round(yearlyPrice * 100))
			}
		} else {
			// Monthly billing: charge monthly price each month
			unitAmount = int64(math.Round(item.Price * 100))
			interval = stripe.PriceRecurringIntervalMonth
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(item.Name),
					Description: stripe.String(fmt.Sprintf("%s - %s billing", item.ItemType, cart.BillingCycle)),
				},
				UnitAmount: stripe.Int64(unitAmount),
				Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
					Interval: stripe.String(string(interval)),
				},
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})

		totalCents += unitAmount * int64(item.Quantity)
	}

	// Create Stripe Checkout session
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(s.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.cancelURL),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"cart_session_id": sessionID,
		},
	}

	stripeSession, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	// Create order record
	orderTotal := float64(totalCents) / 100
	order := &model.Order{
		CartID:          cart.ID,
		StripeSessionID: stripeSession.ID,
		Items:           cart.Items,
		BillingCycle:    cart.BillingCycle,
		TotalAmount:     orderTotal,
		Status:          "pending",
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return stripeSession.URL, nil
}

// HandleWebhook processes Stripe webhook events
func (s *CheckoutService) HandleWebhook(ctx context.Context, stripeSessionID, status string) error {
	order, err := s.orderRepo.FindByStripeSessionID(ctx, stripeSessionID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	return s.orderRepo.UpdateStatus(ctx, order.ID, status)
}
