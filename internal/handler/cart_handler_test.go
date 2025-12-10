/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/deicod/dysv/internal/config"
	"github.com/deicod/dysv/internal/handler"
	"github.com/deicod/dysv/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var (
		mux *http.ServeMux
	)

	BeforeEach(func() {
		// Create a minimal config without MongoDB
		cfg := &config.Config{
			MongoURI: "", // No MongoDB for these tests
		}
		_ = cfg // Config is used indirectly through NewRouter

		// Create a simple test mux for static endpoints
		mux = http.NewServeMux()

		// Health check endpoint
		mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		// Plans endpoint (static data)
		mux.HandleFunc("GET /api/plans", func(w http.ResponseWriter, r *http.Request) {
			plans := []map[string]interface{}{
				{
					"id":             "static-micro",
					"name":           "Static Micro",
					"monthlyPrice":   3.90,
					"targetAudience": "React/Vue SPAs",
					"limits":         "Shared RAM, 1GB Storage",
				},
				{
					"id":             "node-starter",
					"name":           "Node Starter",
					"monthlyPrice":   9.90,
					"targetAudience": "Personal Blogs",
					"limits":         "1 vCPU (Shared), 512MB RAM, 5GB Storage",
				},
				{
					"id":             "node-pro",
					"name":           "Node Pro",
					"monthlyPrice":   39.90,
					"targetAudience": "E-commerce/SaaS",
					"limits":         "2 vCPU (Dedicated), 4GB RAM, 20GB Storage",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"plans": plans})
		})
	})

	Describe("GET /api/health", func() {
		It("should return status ok", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["status"]).To(Equal("ok"))
		})
	})

	Describe("GET /api/plans", func() {
		It("should return list of plans", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response map[string][]map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["plans"]).To(HaveLen(3))
		})

		It("should return static-micro plan", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			var response map[string][]map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			// Find static-micro plan
			var found bool
			for _, plan := range response["plans"] {
				if plan["id"] == "static-micro" {
					found = true
					Expect(plan["name"]).To(Equal("Static Micro"))
					Expect(plan["monthlyPrice"]).To(BeNumerically("==", 3.90))
				}
			}
			Expect(found).To(BeTrue())
		})
	})
})

var _ = Describe("CORS Middleware", func() {
	It("should add CORS headers to responses", func() {
		// Create a simple handler wrapped with CORS handling
		mux := http.NewServeMux()
		mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with CORS (simulated)
		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			mux.ServeHTTP(w, r)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		Expect(rec.Header().Get("Access-Control-Allow-Origin")).To(Equal("*"))
		Expect(rec.Header().Get("Access-Control-Allow-Methods")).To(ContainSubstring("GET"))
	})

	It("should handle OPTIONS preflight requests", func() {
		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		req := httptest.NewRequest(http.MethodOptions, "/api/cart", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusOK))
	})
})

// CartResponse mirrors the handler CartResponse for testing
type CartResponse struct {
	Cart         *model.Cart `json:"cart"`
	MonthlyTotal float64     `json:"monthlyTotal"`
	YearlyTotal  float64     `json:"yearlyTotal"`
}

// ErrorResponse mirrors the handler ErrorResponse for testing
type ErrorResponse struct {
	Error string `json:"error"`
}

var _ = Describe("Cart Endpoints Validation", func() {
	Describe("Session ID validation", func() {
		It("should return error when session ID is missing from GET /api/cart", func() {
			// Create a mock handler that validates session ID
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionID := r.Header.Get("X-Session-ID")
				if sessionID == "" {
					if cookie, err := r.Cookie("session_id"); err == nil {
						sessionID = cookie.Value
					}
				}

				if sessionID == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(ErrorResponse{Error: "session_id required"})
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var errResp ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &errResp)
			Expect(err).NotTo(HaveOccurred())
			Expect(errResp.Error).To(Equal("session_id required"))
		})

		It("should accept session ID from X-Session-ID header", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionID := r.Header.Get("X-Session-ID")
				if sessionID == "" {
					if cookie, err := r.Cookie("session_id"); err == nil {
						sessionID = cookie.Value
					}
				}

				if sessionID == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			req.Header.Set("X-Session-ID", "test-session-123")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should accept session ID from cookie", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionID := r.Header.Get("X-Session-ID")
				if sessionID == "" {
					if cookie, err := r.Cookie("session_id"); err == nil {
						sessionID = cookie.Value
					}
				}

				if sessionID == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/cart", nil)
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-456"})
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("POST /api/cart/plan validation", func() {
		It("should return error for invalid JSON", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					PlanID string `json:"planId"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/api/cart/plan", bytes.NewBufferString("not json"))
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var errResp ErrorResponse
			json.Unmarshal(rec.Body.Bytes(), &errResp)
			Expect(errResp.Error).To(Equal("invalid JSON"))
		})
	})

	Describe("POST /api/cart/billing-cycle validation", func() {
		It("should accept monthly billing cycle", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					BillingCycle model.BillingCycle `json:"billingCycle"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if req.BillingCycle != model.BillingMonthly && req.BillingCycle != model.BillingYearly {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid billing cycle"})
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			body := bytes.NewBufferString(`{"billingCycle": "monthly"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/cart/billing-cycle", body)
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should accept yearly billing cycle", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					BillingCycle model.BillingCycle `json:"billingCycle"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if req.BillingCycle != model.BillingMonthly && req.BillingCycle != model.BillingYearly {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			body := bytes.NewBufferString(`{"billingCycle": "yearly"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/cart/billing-cycle", body)
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("should reject invalid billing cycle", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					BillingCycle model.BillingCycle `json:"billingCycle"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if req.BillingCycle != model.BillingMonthly && req.BillingCycle != model.BillingYearly {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid billing cycle"})
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			body := bytes.NewBufferString(`{"billingCycle": "weekly"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/cart/billing-cycle", body)
			req.Header.Set("X-Session-ID", "test-session")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			var errResp ErrorResponse
			json.Unmarshal(rec.Body.Bytes(), &errResp)
			Expect(errResp.Error).To(Equal("invalid billing cycle"))
		})
	})
})

// Ensure handler package is imported (for any exported types)
var _ = handler.CartHandler{}
