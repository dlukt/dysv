package model

import "time"

// Address represents a customer's physical address.
type Address struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	UserID     string    `json:"userId" bson:"user_id"`
	Label      string    `json:"label" bson:"label"` // e.g. "Home", "Office"
	Line1      string    `json:"line1" bson:"line1"`
	Line2      string    `json:"line2,omitempty" bson:"line2,omitempty"`
	City       string    `json:"city" bson:"city"`
	PostalCode string    `json:"postalCode" bson:"postal_code"`
	State      string    `json:"state,omitempty" bson:"state,omitempty"`
	Country    string    `json:"country" bson:"country"` // ISO 3166-1 alpha-2
	IsDefault  bool      `json:"isDefault" bson:"is_default"`
	CreatedAt  time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updated_at"`
}
