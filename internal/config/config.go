/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Port                string        `mapstructure:"PORT"`
	BaseURL             string        `mapstructure:"BASE_URL"`
	MongoURI            string        `mapstructure:"MONGODB_URI"`
	MongoTimeout        time.Duration `mapstructure:"MONGODB_TIMEOUT"`
	StripeSecret        string        `mapstructure:"STRIPE_SECRET"`
	StripePubKey        string        `mapstructure:"STRIPE_PUBLIC_KEY"`
	StripeWebhookSecret string        `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	StripeAPIVersion    string        `mapstructure:"STRIPE_API_VERSION"`
	AuthSessionSecret   string        `mapstructure:"AUTH_SESSION_SECRET"`
	AuthTokenSecret     string        `mapstructure:"AUTH_TOKEN_SECRET"`
	AuthEmailFrom       string        `mapstructure:"AUTH_EMAIL_FROM"`
	AuthEmailHost       string        `mapstructure:"AUTH_EMAIL_HOST"`
	AuthEmailPort       int           `mapstructure:"AUTH_EMAIL_PORT"`
	AuthEmailUser       string        `mapstructure:"AUTH_EMAIL_USER"`
	AuthEmailPass       string        `mapstructure:"AUTH_EMAIL_PASS"`
	AuthEmailUseSSL     bool          `mapstructure:"AUTH_EMAIL_USE_SSL"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("BASE_URL", "https://dysv.de")
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017/dysv")
	viper.SetDefault("MONGODB_TIMEOUT", "30s")

	timeout, err := time.ParseDuration(viper.GetString("MONGODB_TIMEOUT"))
	if err != nil {
		log.Printf("Config: ParseDuration error: %v\n", err)
		timeout = 30 * time.Second
	}

	cfg := &Config{
		Port:                viper.GetString("PORT"),
		BaseURL:             viper.GetString("BASE_URL"),
		MongoURI:            viper.GetString("MONGODB_URI"),
		MongoTimeout:        timeout,
		StripeSecret:        viper.GetString("STRIPE_SECRET"),
		StripePubKey:        viper.GetString("STRIPE_PUBLIC_KEY"),
		StripeWebhookSecret: viper.GetString("STRIPE_WEBHOOK_SECRET"),
		StripeAPIVersion:    viper.GetString("STRIPE_API_VERSION"),
		AuthSessionSecret:   viper.GetString("AUTH_SESSION_SECRET"),
		AuthTokenSecret:     viper.GetString("AUTH_TOKEN_SECRET"),
		AuthEmailFrom:       viper.GetString("AUTH_EMAIL_FROM"),
		AuthEmailHost:       viper.GetString("AUTH_EMAIL_HOST"),
		AuthEmailPort:       viper.GetInt("AUTH_EMAIL_PORT"),
		AuthEmailUser:       viper.GetString("AUTH_EMAIL_USER"),
		AuthEmailPass:       viper.GetString("AUTH_EMAIL_PASS"),
		AuthEmailUseSSL:     viper.GetBool("AUTH_EMAIL_USE_SSL"),
	}

	return cfg, nil
}
