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
	cartService    *CartService
	orderRepo      repo.OrderRepository
	addressService *AddressService
	successURL     string
	cancelURL      string
}

// NewCheckoutService creates a new checkout service
func NewCheckoutService(cartService *CartService, orderRepo repo.OrderRepository, addressService *AddressService, stripeKey, successURL, cancelURL string) *CheckoutService {
	stripe.Key = stripeKey
	return &CheckoutService{
		cartService:    cartService,
		orderRepo:      orderRepo,
		addressService: addressService,
		successURL:     successURL,
		cancelURL:      cancelURL,
	}
}

// CreateCheckoutSession creates a Stripe Checkout session for the cart
func (s *CheckoutService) CreateCheckoutSession(ctx context.Context, sessionID, userID, addressID string) (string, error) {
	cart, err := s.cartService.GetOrCreateCart(ctx, sessionID)
	if err != nil {
		fmt.Printf("CheckoutService: GetOrCreateCart Error: %v\n", err)
		return "", err
	}

	if len(cart.Items) == 0 {
		return "", ErrEmptyCart
	}

	// Fetch Address
	// Using repo directly via interface or via service? Service!
	// CheckoutService uses AddressService.
	// But CheckoutService needs access to GetAddress logic.
	// Assume AddressService exposes GetAddress.
	// Wait, I implemented List, Create, Update, Delete in AddressService. Did I implement Get?
	// Let's verify AddressService has Get.
	// I forgot to add GetAddress to AddressService in implementation plan step...
	// I'll assume I need to ADD it now or use repo directly?
	// Better: Add GetAddress to Service now.

	// Assuming GetAddress exists or I'll add it.
	// Let's implement getting address here assuming service has it.
	// If it fails compile, I'll fix service.

	address, err := s.addressService.GetAddress(ctx, addressID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil {
		return "", fmt.Errorf("address not found or does not belong to user")
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
			"address_id":      addressID,
		},
		// Optional: Pre-fill customer email if we knew it, or address from our DB?
		// Stripe allows passing address collection fields.
	}

	stripeSession, err := session.New(params)
	if err != nil {
		fmt.Printf("CheckoutService: Stripe Session New Error: %v\n", err)
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	// Create order record
	orderTotal := float64(totalCents) / 100
	order := &model.Order{
		CartID:          cart.ID,
		StripeSessionID: stripeSession.ID,
		Items:           cart.Items,
		BillingCycle:    cart.BillingCycle,
		BillingAddress:  *address, // Store snapshot
		TotalAmount:     orderTotal,
		Status:          "pending",
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		fmt.Printf("CheckoutService: OrderRepo Create Error: %v\n", err)
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return stripeSession.URL, nil
}

// HandleWebhook processes Stripe webhook events
func (s *CheckoutService) HandleWebhook(ctx context.Context, stripeSessionID, status string) error {
	order, err := s.orderRepo.FindByStripeSessionID(ctx, stripeSessionID)
	if err != nil {
		fmt.Printf("CheckoutService: FindByStripeSessionID Error: %v\n", err)
		return fmt.Errorf("order not found: %w", err)
	}

	return s.orderRepo.UpdateStatus(ctx, order.ID, status)
}
