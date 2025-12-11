package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/deicod/auth/core"
	"github.com/deicod/dysv/internal/handler"
	"github.com/deicod/dysv/internal/mocks"
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressHandler", func() {
	var (
		mockRepo    *mocks.MockAddressRepo
		mockAuth    *mocks.MockAuthService
		addrService *service.AddressService
		addrHandler *handler.AddressHandler
		ctx         context.Context
		userID      string
	)

	BeforeEach(func() {
		mockRepo = mocks.NewMockAddressRepo()
		mockAuth = &mocks.MockAuthService{}
		addrService = service.NewAddressService(mockRepo)
		addrHandler = handler.NewAddressHandler(addrService, mockAuth)
		ctx = context.Background()
		userID = "user_123"

		// Default successful auth
		mockAuth.AuthenticateSessionFunc = func(ctx context.Context, token string) (core.UserPublic, core.SessionPublic, error) {
			if token == "valid-token" {
				return core.UserPublic{ID: core.ID(userID)}, core.SessionPublic{}, nil
			}
			return core.UserPublic{}, core.SessionPublic{}, errors.New("invalid token")
		}
	})

	Describe("List", func() {
		It("should list addresses for authenticated user", func() {
			// Seed data
			err := mockRepo.Create(ctx, &model.Address{UserID: userID, Label: "Home"})
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest(http.MethodGet, "/api/user/addresses", nil)
			req.Header.Set("Authorization", "Bearer valid-token")
			rec := httptest.NewRecorder()

			addrHandler.List(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response map[string][]model.Address
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["addresses"]).To(HaveLen(1))
			Expect(response["addresses"][0].Label).To(Equal("Home"))
		})

		It("should return 401 for invalid token", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/user/addresses", nil)
			req.Header.Set("Authorization", "Bearer invalid-token")
			rec := httptest.NewRecorder()

			addrHandler.List(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("Create", func() {
		It("should create address", func() {
			body := map[string]interface{}{
				"label":       "Work",
				"street":      "456 Work St",
				"countryCode": "DE",
			}
			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/addresses", bytes.NewBuffer(jsonBody))
			req.Header.Set("Authorization", "Bearer valid-token")
			rec := httptest.NewRecorder()

			addrHandler.Create(rec, req)

			Expect(rec.Code).To(Equal(http.StatusCreated))

			var created model.Address
			err := json.Unmarshal(rec.Body.Bytes(), &created)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Label).To(Equal("Work"))
			Expect(created.UserID).To(Equal(userID))
		})

		It("should set default label if missing", func() {
			body := map[string]interface{}{
				"street": "No Label",
			}
			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/addresses", bytes.NewBuffer(jsonBody))
			req.Header.Set("Authorization", "Bearer valid-token")
			rec := httptest.NewRecorder()

			addrHandler.Create(rec, req)

			Expect(rec.Code).To(Equal(http.StatusCreated))

			var created model.Address
			err := json.Unmarshal(rec.Body.Bytes(), &created)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Label).To(Equal("Default"))
		})
	})

	Describe("Update", func() {
		var existingID string

		BeforeEach(func() {
			addr := &model.Address{UserID: userID, Label: "Old"}
			err := mockRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())
			existingID = addr.ID
		})

		It("should update address", func() {
			body := map[string]interface{}{
				"label": "New Label",
			}
			jsonBody, _ := json.Marshal(body)
			url := "/api/user/addresses/" + existingID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Authorization", "Bearer valid-token")
			rec := httptest.NewRecorder()

			addrHandler.Update(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var updated model.Address
			err := json.Unmarshal(rec.Body.Bytes(), &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Label).To(Equal("New Label"))
			Expect(updated.ID).To(Equal(existingID))
		})
	})

	Describe("Delete", func() {
		var existingID string

		BeforeEach(func() {
			addr := &model.Address{UserID: userID}
			err := mockRepo.Create(ctx, addr)
			Expect(err).NotTo(HaveOccurred())
			existingID = addr.ID
		})

		It("should delete address", func() {
			url := "/api/user/addresses/" + existingID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Authorization", "Bearer valid-token")
			rec := httptest.NewRecorder()

			addrHandler.Delete(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			// Verify deletion
			stored, _ := mockRepo.Get(ctx, existingID, userID)
			Expect(stored).To(BeNil())
		})
	})
})
