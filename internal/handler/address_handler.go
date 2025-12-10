package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/deicod/auth"
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/service"
)

type AddressHandler struct {
	service *service.AddressService
	auth    auth.Service
}

func NewAddressHandler(service *service.AddressService, auth auth.Service) *AddressHandler {
	return &AddressHandler{service: service, auth: auth}
}

func (h *AddressHandler) getUserID(r *http.Request) string {
	token := getToken(r)
	if token == "" {
		return ""
	}
	user, _, err := h.auth.AuthenticateSession(r.Context(), token)
	if err != nil {
		return ""
	}
	return string(user.ID)
}

func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list addresses")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"addresses": addresses})
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var addr model.Address
	if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	addr.UserID = userID
	if addr.Label == "" {
		addr.Label = "Default"
	}

	if err := h.service.CreateAddress(r.Context(), &addr); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create address")
		return
	}

	writeJSON(w, http.StatusCreated, addr)
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	id := pathParts[4] // /api/user/addresses/{id}

	var addr model.Address
	if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	addr.ID = id
	addr.UserID = userID

	if err := h.service.UpdateAddress(r.Context(), &addr); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update address")
		return
	}

	writeJSON(w, http.StatusOK, addr)
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	id := pathParts[4]

	if err := h.service.DeleteAddress(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete address")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// UserID Extraction helper (Temporary until middleware injects it into context)
// For now, we assume the Auth handler validated token and client passes it?
// Actually, in `Me`, we decode session.
// We need middleware to validate token and extract UserID for these protected routes.
// OR, we just decode it here again using auth service?
// Since `auth` logic is encapsulated, we should arguably have an auth middleware.
// But for MVP, `getUserID` might need to look at header and call `authService.AuthenticateSession`.
// Wait, `AddressHandler` doesn't have access to `AuthService`.
//
// Plan adjustment:
// I need `AuthService` in `AddressHandler` OR a middleware that does verification.
// Simpler: Pass `authService` to `AddressHandler` and use it to verify token.
//
// Let's refactor `NewAddressHandler` to take `authService`.
