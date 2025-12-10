/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/deicod/dysv/internal/handler"
	"github.com/deicod/dysv/internal/repo"
	"github.com/deicod/dysv/internal/service"
)

var _ = Describe("Checkout Handler", func() {
	var (
		mockCartRepo  *repo.MockCartRepo
		mockOrderRepo *repo.MockOrderRepo
		cartService   *service.CartService
		cartHandler   *handler.CartHandler
	)

	BeforeEach(func() {
		mockCartRepo = repo.NewMockCartRepo()
		mockOrderRepo = repo.NewMockOrderRepo()
		cartService = service.NewCartService(mockCartRepo)
		cartHandler = handler.NewCartHandler(cartService)
		_ = mockOrderRepo // Will be used when we test checkout
	})

	AfterEach(func() {
		mockCartRepo.Reset()
		mockOrderRepo.Reset()
	})

	Describe("GET /api/cart", func() {
		It("should return empty cart for new session", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			req.Header.Set("X-Session-ID", "new-session-123")
			rec := httptest.NewRecorder()

			cartHandler.GetCart(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			cart := resp["cart"].(map[string]interface{})
			items := cart["items"].([]interface{})
			Expect(items).To(BeEmpty())
		})

		It("should return error when session ID is missing", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			rec := httptest.NewRecorder()

			cartHandler.GetCart(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("POST /api/cart/plan", func() {
		It("should add plan to cart", func() {
			body := `{"planId": "node-starter"}`
			req := httptest.NewRequest(http.MethodPost, "/api/cart/plan", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			cartHandler.SetPlan(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			cart := resp["cart"].(map[string]interface{})
			items := cart["items"].([]interface{})
			Expect(items).To(HaveLen(1))

			item := items[0].(map[string]interface{})
			Expect(item["itemId"]).To(Equal("node-starter"))
			Expect(item["itemType"]).To(Equal("plan"))
		})

		It("should reject invalid plan", func() {
			body := `{"planId": "invalid-plan"}`
			req := httptest.NewRequest(http.MethodPost, "/api/cart/plan", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			cartHandler.SetPlan(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should replace existing plan", func() {
			// Add first plan
			body := `{"planId": "static-micro"}`
			req := httptest.NewRequest(http.MethodPost, "/api/cart/plan", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			cartHandler.SetPlan(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))

			// Replace with second plan
			body = `{"planId": "node-pro"}`
			req = httptest.NewRequest(http.MethodPost, "/api/cart/plan", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()
			cartHandler.SetPlan(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			cart := resp["cart"].(map[string]interface{})
			items := cart["items"].([]interface{})
			Expect(items).To(HaveLen(1))

			item := items[0].(map[string]interface{})
			Expect(item["itemId"]).To(Equal("node-pro"))
		})
	})

	Describe("POST /api/cart/addon", func() {
		It("should add addon to cart", func() {
			body := `{"addonId": "de-domain"}`
			req := httptest.NewRequest(http.MethodPost, "/api/cart/addon", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			cartHandler.AddAddon(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			cart := resp["cart"].(map[string]interface{})
			items := cart["items"].([]interface{})
			Expect(items).To(HaveLen(1))

			item := items[0].(map[string]interface{})
			Expect(item["itemId"]).To(Equal("de-domain"))
			Expect(item["itemType"]).To(Equal("addon"))
		})

		It("should not duplicate addons", func() {
			body := `{"addonId": "de-domain"}`

			// Add addon twice
			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodPost, "/api/cart/addon", stringReader(body))
				req.Header.Set("X-Session-ID", "test-session")
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				cartHandler.AddAddon(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))
			}

			// Verify only one addon
			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			req.Header.Set("X-Session-ID", "test-session")
			rec := httptest.NewRecorder()
			cartHandler.GetCart(rec, req)

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			cart := resp["cart"].(map[string]interface{})
			items := cart["items"].([]interface{})
			Expect(items).To(HaveLen(1))
		})
	})

	Describe("Cart Totals", func() {
		It("should calculate correct yearly totals", func() {
			// Add plan (€9.90/mo)
			body := `{"planId": "node-starter"}`
			req := httptest.NewRequest(http.MethodPost, "/api/cart/plan", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			cartHandler.SetPlan(rec, req)

			// Add addon (€1.00/mo)
			body = `{"addonId": "de-domain"}`
			req = httptest.NewRequest(http.MethodPost, "/api/cart/addon", stringReader(body))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()
			cartHandler.AddAddon(rec, req)

			// Get cart totals
			req = httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			req.Header.Set("X-Session-ID", "test-session")
			rec = httptest.NewRecorder()
			cartHandler.GetCart(rec, req)

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())

			monthlyTotal := resp["monthlyTotal"].(float64)
			yearlyTotal := resp["yearlyTotal"].(float64)

			// Monthly: 9.90 + 1.00 = 10.90
			Expect(monthlyTotal).To(BeNumerically("~", 10.90, 0.01))

			// Yearly: (9.90 * 10) + (1.00 * 12) = 99.00 + 12.00 = 111.00
			Expect(yearlyTotal).To(BeNumerically("~", 111.00, 0.01))
		})
	})
})

func stringReader(s string) *stringReaderType {
	return &stringReaderType{s: s, i: 0}
}

type stringReaderType struct {
	s string
	i int
}

func (r *stringReaderType) Read(p []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, nil
	}
	n = copy(p, r.s[r.i:])
	r.i += n
	return n, nil
}
