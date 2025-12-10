/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deicod/dysv/internal/config"
	"github.com/deicod/dysv/internal/handler"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the dysv API server",
	Long:  `Starts the HTTP API server for dysv.de hosting platform.`,
	Run:   runAPI,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func runAPI(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	// Log loaded config
	fmt.Println("--- Loaded Config ---")
	fmt.Printf("PORT: %s\n", cfg.Port)
	fmt.Printf("BASE_URL: %s\n", cfg.BaseURL)
	fmt.Printf("MONGODB_URI: %s\n", cfg.MongoURI) // Ideally mask user/pass
	fmt.Printf("MONGODB_TIMEOUT: %v\n", cfg.MongoTimeout)
	mask := func(s string) string {
		if len(s) < 8 {
			return "***"
		}
		return s[:4] + "..." + s[len(s)-4:]
	}
	fmt.Printf("STRIPE_SECRET: %s\n", mask(cfg.StripeSecret))
	fmt.Printf("STRIPE_PUBLIC_KEY: %s\n", mask(cfg.StripePubKey))
	fmt.Printf("STRIPE_WEBHOOK_SECRET: %s\n", mask(cfg.StripeWebhookSecret))
	fmt.Println("---------------------")
	// Create router with handlers
	mux := handler.NewRouter(cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting API server on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
