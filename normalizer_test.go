package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNormalizeProviderA(t *testing.T) {
	rawJson := `{
		"provider_name":"ProviderA",
		"id":"A-123",
		"cost":"123.45",
		"currency":"USD",
		"depart":"2025-11-01T09:00:00Z",
		"class":"economy"
	}`

	var raw map[string]interface{}
	json.Unmarshal([]byte(rawJson), &raw)

	offer, err := Normalize(raw)
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if offer.ProviderID != "ProviderA" {
		t.Fatalf("expected ProviderA got %s", offer.ProviderID)
	}
	if offer.Price != 123.45 {
		t.Fatalf("expected price 123.45 got %f", offer.Price)
	}
	if offer.DepartAt.IsZero() {
		t.Fatalf("depart time empty")
	}
}

func TestDedupeStore(t *testing.T) {
	store := NewDedupeStore()
	offer := &Offer{
		OfferID:    "hash1",
		ProviderID: "P",
		Price:      10,
		CreatedAt:  time.Now(),
	}
	if store.IsDuplicate(offer) {
		t.Fatalf("store empty but found duplicate")
	}
	store.Add(offer)
	if !store.IsDuplicate(offer) {
		t.Fatalf("expected duplicate after add")
	}
}
