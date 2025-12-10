/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MadAppGang/httplog"
	"github.com/deicod/auth"
	"github.com/deicod/auth/mgo"
	"github.com/deicod/dysv/internal/config"
	"github.com/deicod/dysv/internal/repo"
	"github.com/deicod/dysv/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewRouter creates the HTTP router with all handlers
func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Printf("Warning: MongoDB connection failed: %v (running without persistence)", err)
	} else {
		// Ping to verify connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Printf("Warning: MongoDB ping failed: %v (running without persistence)", err)
			client = nil
		}
	}

	var cartHandler *CartHandler
	var checkoutHandler *CheckoutHandler
	var authHandler *AuthHandler
	var addressHandler *AddressHandler

	if client != nil {
		db := client.Database("dysv")

		// Repositories
		cartRepo := repo.NewCartRepo(db, cfg.MongoTimeout)
		orderRepo := repo.NewOrderRepo(db, cfg.MongoTimeout)
		addressRepo := repo.NewAddressRepo(db, cfg.MongoTimeout)

		// Services
		cartService := service.NewCartService(cartRepo)
		addressService := service.NewAddressService(addressRepo)

		// Auth Service Initialization
		ac := auth.DefaultConfig()
		ac.Backend = auth.BackendMongo
		ac.Mongo = &mgo.Config{
			URI:              cfg.MongoURI,
			Database:         "dysv",
			OperationTimeout: cfg.MongoTimeout,
		}

		if cfg.AuthEmailHost != "" {
			ac.Email.Host = cfg.AuthEmailHost
			ac.Email.Port = cfg.AuthEmailPort
			ac.Email.User = cfg.AuthEmailUser
			ac.Email.Pass = cfg.AuthEmailPass
			ac.Email.UseSSL = cfg.AuthEmailUseSSL
		}
		if cfg.AuthEmailFrom != "" {
			ac.Email.From = cfg.AuthEmailFrom
		}

		var authSvc auth.Service
		authSvc, err := auth.NewService(context.Background(), ac)
		if err != nil {
			log.Printf("Error: Failed to initialize Auth Service: %v", err)
		} else {
			authHandler = NewAuthHandler(authSvc)
			addressHandler = NewAddressHandler(addressService, authSvc)
		}

		// Handlers & Checkout Service
		cartHandler = NewCartHandler(cartService)

		if cfg.StripeSecret != "" {
			successURL := cfg.BaseURL + "/checkout/success"
			cancelURL := cfg.BaseURL + "/cart"
			// CheckoutService needs AddressService
			checkoutService := service.NewCheckoutService(cartService, orderRepo, addressService, cfg.StripeSecret, successURL, cancelURL)

			// CheckoutHandler needs Auth Service (authSvc)
			// Ensure authSvc is not nil
			if authSvc != nil {
				checkoutHandler = NewCheckoutHandler(checkoutService, authSvc, cfg.StripeWebhookSecret, cfg.StripeAPIVersion)
			} else {
				log.Println("Warning: CheckoutHandler disabled because Auth Service failed to initialize")
			}
		}
	}

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Plans endpoint (static data)
	mux.HandleFunc("GET /api/plans", func(w http.ResponseWriter, r *http.Request) {
		plans := make([]map[string]interface{}, 0, len(service.Plans))
		for _, plan := range service.Plans {
			plans = append(plans, map[string]interface{}{
				"id":             plan.ID,
				"name":           plan.Name,
				"monthlyPrice":   plan.MonthlyPrice,
				"targetAudience": plan.TargetAudience,
				"limits":         plan.Limits,
			})
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"plans": plans})
	})

	// Auth endpoints
	if authHandler != nil {
		mux.HandleFunc("POST /api/auth/register", authHandler.Register)
		mux.HandleFunc("POST /api/auth/login", authHandler.Login)
		mux.HandleFunc("POST /api/auth/verify", authHandler.VerifyEmail)
		mux.HandleFunc("POST /api/auth/forgot-password", authHandler.ForgotPassword)
		mux.HandleFunc("POST /api/auth/reset-password", authHandler.ResetPassword)
		mux.HandleFunc("GET /api/auth/me", authHandler.Me)
	}

	// Address endpoints
	if addressHandler != nil {
		mux.HandleFunc("GET /api/user/addresses", addressHandler.List)
		mux.HandleFunc("POST /api/user/addresses", addressHandler.Create)
		mux.HandleFunc("PUT /api/user/addresses/", addressHandler.Update)    // Trailing slash for ID parsing in handler? No.
		mux.HandleFunc("DELETE /api/user/addresses/", addressHandler.Delete) // Handlers parse path manually
	}

	// Cart endpoints (require MongoDB)
	if cartHandler != nil {
		mux.HandleFunc("GET /api/cart", cartHandler.GetCart)
		mux.HandleFunc("POST /api/cart/plan", cartHandler.AddPlan)
		mux.HandleFunc("POST /api/cart/addon", cartHandler.AddAddon)
		mux.HandleFunc("DELETE /api/cart/item/{itemId}", cartHandler.RemoveItem)
		mux.HandleFunc("PUT /api/cart/item/{itemId}", cartHandler.UpdateItemQuantity)
		mux.HandleFunc("POST /api/cart/billing-cycle", cartHandler.SetBillingCycle)
	} else {
		// Return error if MongoDB not available
		mongoRequired := func(w http.ResponseWriter, r *http.Request) {
			writeError(w, http.StatusServiceUnavailable, "database not available")
		}
		mux.HandleFunc("GET /api/cart", mongoRequired)
		mux.HandleFunc("POST /api/cart/plan", mongoRequired)
		mux.HandleFunc("POST /api/cart/addon", mongoRequired)
		mux.HandleFunc("DELETE /api/cart/item/{itemId}", mongoRequired)
		mux.HandleFunc("POST /api/cart/billing-cycle", mongoRequired)
	}

	// Checkout endpoints (require MongoDB + Stripe)
	if checkoutHandler != nil {
		mux.HandleFunc("POST /api/checkout", checkoutHandler.CreateCheckoutSession)
		mux.HandleFunc("POST /api/webhook/stripe", checkoutHandler.Webhook)
	} else {
		stripeRequired := func(w http.ResponseWriter, r *http.Request) {
			writeError(w, http.StatusServiceUnavailable, "checkout not available")
		}
		mux.HandleFunc("POST /api/checkout", stripeRequired)
		mux.HandleFunc("POST /api/webhook/stripe", stripeRequired)
	}

	// CORS middleware
	handler := corsMiddleware(mux)

	// Logging middleware with headers
	httplog.ForceConsoleColor()
	return httplog.HandlerWithFormatter(httplog.DefaultLogFormatterWithRequestHeader, handler)
}

// corsMiddleware adds CORS headers for frontend development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
