/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deicod/dysv/internal/config"
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/service"
)

// CartHandler handles cart-related HTTP requests
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// getSessionID extracts session ID from request (cookie or header)
func getSessionID(r *http.Request) string {
	// Try cookie first
	if cookie, err := r.Cookie("session_id"); err == nil {
		return cookie.Value
	}
	// Fall back to header
	return r.Header.Get("X-Session-ID")
}

// GetCart handles GET /api/cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	cart, err := h.cartService.GetOrCreateCart(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}

	writeJSON(w, http.StatusOK, resp)
}

// SetPlanRequest is the request body for setting a plan
type SetPlanRequest struct {
	PlanID string `json:"planId"`
}

// SetPlan handles POST /api/cart/plan
func (h *CartHandler) SetPlan(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	var req SetPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	cart, err := h.cartService.SetPlan(r.Context(), sessionID, req.PlanID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPlan) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}

	writeJSON(w, http.StatusOK, resp)
}

// AddAddonRequest is the request body for adding an addon
type AddAddonRequest struct {
	AddonID string `json:"addonId"`
}

// AddAddon handles POST /api/cart/addon
func (h *CartHandler) AddAddon(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	var req AddAddonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	cart, err := h.cartService.AddAddon(r.Context(), sessionID, req.AddonID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAddon) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}

	writeJSON(w, http.StatusOK, resp)
}

// RemoveItem handles DELETE /api/cart/item/{itemId}
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	itemID := r.PathValue("itemId")
	if itemID == "" {
		writeError(w, http.StatusBadRequest, "itemId required")
		return
	}

	cart, err := h.cartService.RemoveItem(r.Context(), sessionID, itemID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}

	writeJSON(w, http.StatusOK, resp)
}

// SetBillingCycleRequest is the request body for setting billing cycle
type SetBillingCycleRequest struct {
	BillingCycle model.BillingCycle `json:"billingCycle"`
}

// SetBillingCycle handles POST /api/cart/billing-cycle
func (h *CartHandler) SetBillingCycle(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	var req SetBillingCycleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.BillingCycle != model.BillingMonthly && req.BillingCycle != model.BillingYearly {
		writeError(w, http.StatusBadRequest, "invalid billing cycle")
		return
	}

	cart, err := h.cartService.SetBillingCycle(r.Context(), sessionID, req.BillingCycle)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}

	writeJSON(w, http.StatusOK, resp)
}

// CartResponse is the response for cart endpoints
type CartResponse struct {
	Cart         *model.Cart `json:"cart"`
	MonthlyTotal float64     `json:"monthlyTotal"`
	YearlyTotal  float64     `json:"yearlyTotal"`
}

// Placeholder for dependency injection
var _ = config.Config{}
