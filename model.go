package main

import "time"

// Canonical Offer model
type Offer struct {
	OfferID     string    `json:"offer_id"`      // canonical id (could be hash)
	ProviderID  string    `json:"provider_id"`   // source provider
	ProviderRef string    `json:"provider_ref"`  // provider's native id/sku
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	FareClass   string    `json:"fare_class,omitempty"`
	DepartAt    time.Time `json:"depart_at,omitempty"`
	ReturnAt    time.Time `json:"return_at,omitempty"`
	Raw         any       `json:"raw,omitempty"` // original payload for traceability
	CreatedAt   time.Time `json:"created_at"`
}
