/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"

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
		fmt.Printf("Handler: GetCart Error: %v\n", err)
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

// AddPlanRequest is the request body for adding a plan
type AddPlanRequest struct {
	PlanID   string `json:"planId"`
	Quantity int    `json:"quantity"`
}

// AddPlan handles POST /api/cart/plan
func (h *CartHandler) AddPlan(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	var req AddPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	qty := req.Quantity
	if qty < 1 {
		qty = 1
	}

	cart, err := h.cartService.AddPlan(r.Context(), sessionID, req.PlanID, qty)
	if err != nil {
		fmt.Printf("Handler: AddPlan Error: %v\n", err)
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

// UpdateItemRequest is the request body for updating item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

// UpdateItemQuantity handles PUT /api/cart/item/{itemId}
func (h *CartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(r.Context(), sessionID, itemID, req.Quantity)
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
	fmt.Println("SetBillingCycle start")
	sessionID := getSessionID(r)
	spew.Dump("sessionID", sessionID)
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id required")
		return
	}

	var req SetBillingCycleRequest
	fmt.Printf("Handler: Reading Body. Content-Length: %d\n", r.ContentLength)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Handler: SetBillingCycle JSON Decode Error: %v\n", err)
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	fmt.Println("r.Body decoded into req successfully")
	spew.Dump("Handler Request Body", req)
	fmt.Printf("Handler: SetBillingCycle to %s for %s\n", req.BillingCycle, sessionID)

	if req.BillingCycle != model.BillingMonthly && req.BillingCycle != model.BillingYearly {
		writeError(w, http.StatusBadRequest, "invalid billing cycle")
		return
	}
	fmt.Println("BillingCycle checked successfully")
	cart, err := h.cartService.SetBillingCycle(r.Context(), sessionID, req.BillingCycle)
	if err != nil {
		fmt.Printf("Handler: SetBillingCycle Error: %v\n", err)
		writeError(w, http.StatusInternalServerError, "SetBillingCycle failed: "+err.Error())
		return
	}
	spew.Dump("cart:", cart)
	fmt.Println("BillingCycle set successfully")

	monthly, yearly := h.cartService.GetCartTotal(cart)
	resp := CartResponse{
		Cart:         cart,
		MonthlyTotal: monthly,
		YearlyTotal:  yearly,
	}
	spew.Dump("monthly, yearly:", monthly, yearly)

	writeJSON(w, http.StatusOK, resp)
	fmt.Println("SetBillingCycle end")
}

// CartResponse is the response for cart endpoints
type CartResponse struct {
	Cart         *model.Cart `json:"cart"`
	MonthlyTotal float64     `json:"monthlyTotal"`
	YearlyTotal  float64     `json:"yearlyTotal"`
}

// Placeholder for dependency injection
var _ = config.Config{}
