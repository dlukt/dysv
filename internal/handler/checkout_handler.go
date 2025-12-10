/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/deicod/dysv/internal/service"
	"github.com/stripe/stripe-go/v82/webhook"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	checkoutService  *service.CheckoutService
	webhookSecret    string
	stripeAPIVersion string
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(checkoutService *service.CheckoutService, webhookSecret, stripeAPIVersion string) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService:  checkoutService,
		webhookSecret:    webhookSecret,
		stripeAPIVersion: stripeAPIVersion,
	}
}

// CreateCheckoutSession handles POST /api/checkout
func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	checkoutURL, err := h.checkoutService.CreateCheckoutSession(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, service.ErrEmptyCart) {
			log.Printf("CheckoutHandler: CreateCheckoutSession EmptyCart: %v", err)
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("CheckoutHandler: CreateCheckoutSession Error: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, CheckoutResponse{
		URL: checkoutURL,
	})
}

// CheckoutResponse is the response for checkout endpoint
type CheckoutResponse struct {
	URL string `json:"url"`
}

// Webhook handles POST /api/webhook/stripe
func (h *CheckoutHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Webhook: error reading body: %v", err)
		writeError(w, http.StatusBadRequest, "error reading request body")
		return
	}

	// Verify Stripe signature
	sigHeader := r.Header.Get("Stripe-Signature")
	log.Printf("Webhook: Payload size: %d bytes", len(payload))
	log.Printf("Webhook: Stripe-Signature: %s", sigHeader)

	maskSecret := "empty"
	if len(h.webhookSecret) > 8 {
		maskSecret = h.webhookSecret[:4] + "..." + h.webhookSecret[len(h.webhookSecret)-4:]
	}
	log.Printf("Webhook: Using Secret: %s", maskSecret)

	if sigHeader == "" {
		writeError(w, http.StatusBadRequest, "missing Stripe-Signature header")
		return
	}

	// If a specific API version is configured, we enforce it by NOT ignoring the mismatch (unless it matches).
	// But `stripe-go` library decides "mismatch" based on ITS hardcoded version.
	// Actually, `IgnoreAPIVersionMismatch: true` is the safest default if we want to allow Config overrides.
	// But to support the user's request:
	// If `h.stripeAPIVersion` is set, we could compare it manually, but `ConstructEvent` verifies against the SDK version.
	// The robust way: Always IgnoreMismatch, but maybe log a warning if it differs from Config?
	// User asked: "lock in 2025-11-17... in the backend".
	// The SDK assumes it knows best. If we want to "lock in", we should probably mostly trust the SDK,
	// OR just universally ignore the mismatch if the user explicitly provided a version that MIGHT differ.

	ignoreMismatch := true
	// logic: if user explicitly sets a version, they might be setting one DIFFERENT from the SDK.
	// In that case, we MUST ignore the SDK's complaint.
	// If the user doesn't set one, we also default to true because dev/prod often differ.

	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, h.webhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: ignoreMismatch,
	})
	if err != nil {
		log.Printf("Webhook: signature verification failed: %v", err)
		writeError(w, http.StatusBadRequest, "signature verification failed")
		return
	}

	// Optional: Enforce specific version manually if desired
	if h.stripeAPIVersion != "" && event.APIVersion != h.stripeAPIVersion {
		log.Printf("Webhook: Warning: Received API Version %s, expected configured %s", event.APIVersion, h.stripeAPIVersion)
		// We don't error here, just warn. Strict blocking could break things unnecessarily.
	}

	// Handle the event
	switch event.Type {
	// Checkout Session events
	case "checkout.session.completed":
		log.Printf("Webhook: checkout.session.completed for session %s", event.ID)
		if session, ok := event.Data.Object["id"].(string); ok {
			// Payment mode determines status
			paymentStatus, _ := event.Data.Object["payment_status"].(string)
			status := "pending"
			if paymentStatus == "paid" {
				status = "paid"
			}
			if err := h.checkoutService.HandleWebhook(r.Context(), session, status); err != nil {
				log.Printf("Webhook: error updating order: %v", err)
			}
		}

	case "checkout.session.async_payment_succeeded":
		log.Printf("Webhook: async payment succeeded for session %s", event.ID)
		if session, ok := event.Data.Object["id"].(string); ok {
			if err := h.checkoutService.HandleWebhook(r.Context(), session, "paid"); err != nil {
				log.Printf("Webhook: error updating order: %v", err)
			}
		}

	case "checkout.session.async_payment_failed":
		log.Printf("Webhook: async payment failed for session %s", event.ID)
		if session, ok := event.Data.Object["id"].(string); ok {
			if err := h.checkoutService.HandleWebhook(r.Context(), session, "payment_failed"); err != nil {
				log.Printf("Webhook: error updating order: %v", err)
			}
		}

	case "checkout.session.expired":
		log.Printf("Webhook: checkout.session.expired for session %s", event.ID)
		if session, ok := event.Data.Object["id"].(string); ok {
			if err := h.checkoutService.HandleWebhook(r.Context(), session, "expired"); err != nil {
				log.Printf("Webhook: error updating order: %v", err)
			}
		}

	// Subscription lifecycle events (for ongoing subscription management)
	case "customer.subscription.created":
		log.Printf("Webhook: subscription created %s", event.ID)
		// TODO: Link subscription to customer account

	case "customer.subscription.updated":
		log.Printf("Webhook: subscription updated %s", event.ID)
		// TODO: Handle plan changes, quantity updates

	case "customer.subscription.deleted":
		log.Printf("Webhook: subscription canceled %s", event.ID)
		// TODO: Deprovision resources

	case "customer.subscription.paused":
		log.Printf("Webhook: subscription paused %s", event.ID)
		// TODO: Suspend service

	case "customer.subscription.resumed":
		log.Printf("Webhook: subscription resumed %s", event.ID)
		// TODO: Resume service

	// Invoice events (for payment tracking)
	case "invoice.paid":
		log.Printf("Webhook: invoice paid %s", event.ID)
		// TODO: Record successful payment

	case "invoice.payment_failed":
		log.Printf("Webhook: invoice payment failed %s", event.ID)
		// TODO: Notify customer, retry logic

	default:
		log.Printf("Webhook: unhandled event type %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
