/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BillingCycle represents monthly or yearly billing
type BillingCycle string

const (
	BillingMonthly BillingCycle = "monthly"
	BillingYearly  BillingCycle = "yearly"
)

// LineItem represents a single item in a cart or order
type LineItem struct {
	ItemID   string  `bson:"item_id" json:"itemId"`
	ItemType string  `bson:"item_type" json:"itemType"` // "plan" or "addon"
	Name     string  `bson:"name" json:"name"`
	Price    float64 `bson:"price" json:"price"` // Monthly price
	Quantity int     `bson:"quantity" json:"quantity"`
}

// Cart represents a shopping cart
type Cart struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID    string        `bson:"session_id" json:"sessionId"`
	Items        []LineItem    `bson:"items" json:"items"`
	BillingCycle BillingCycle  `bson:"billing_cycle" json:"billingCycle"`
	CreatedAt    time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updatedAt"`
}

// Order represents a completed order
type Order struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CartID          bson.ObjectID `bson:"cart_id" json:"cartId"`
	StripeSessionID string        `bson:"stripe_session_id" json:"stripeSessionId"`
	CustomerEmail   string        `bson:"customer_email" json:"customerEmail"`
	Items           []LineItem    `bson:"items" json:"items"`
	BillingCycle    BillingCycle  `bson:"billing_cycle" json:"billingCycle"`
	TotalAmount     float64       `bson:"total_amount" json:"totalAmount"`
	Status          string        `bson:"status" json:"status"` // pending, paid, cancelled
	CreatedAt       time.Time     `bson:"created_at" json:"createdAt"`
	PaidAt          *time.Time    `bson:"paid_at,omitempty" json:"paidAt,omitempty"`
}

// Plan represents a hosting plan (for reference, not stored in DB)
type Plan struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	MonthlyPrice   float64  `json:"monthlyPrice"`
	TargetAudience string   `json:"targetAudience"`
	Limits         string   `json:"limits"`
	Features       []string `json:"features"`
}

// Addon represents an add-on product
type Addon struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	MonthlyPrice float64 `json:"monthlyPrice"`
}
