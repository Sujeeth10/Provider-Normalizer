package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Normalize tries provider-specific parsers and returns canonical Offer.
func Normalize(raw map[string]interface{}) (*Offer, error) {
	// naive detection heuristics. A real system uses explicit provider id.
	if _, ok := raw["provider_name"]; ok && raw["provider_name"] == "ProviderA" {
		return normalizeProviderA(raw)
	}
	if _, ok := raw["vendor"]; ok {
		return normalizeProviderB(raw)
	}

	// fallback attempt: if fields like "price" exist use them
	if _, ok := raw["price"]; ok {
		return normalizeGeneric(raw)
	}
	return nil, errors.New("unknown provider/schema")
}

func normalizeProviderA(raw map[string]interface{}) (*Offer, error) {
	// ProviderA schema example:
	// { "provider_name": "ProviderA", "id": "abc123", "cost": "123.45", "currency": "USD", "depart": "2025-11-01T09:00:00Z" }

	id, _ := raw["id"].(string)
	cost := parseNumber(raw["cost"])
	curr, _ := raw["currency"].(string)
	depart := parseTime(raw["depart"])
	fare, _ := raw["class"].(string)
	offer := &Offer{
		ProviderID:  "ProviderA",
		ProviderRef: id,
		Price:       cost,
		Currency:    curr,
		FareClass:   fare,
		DepartAt:    depart,
		Raw:         raw,
		CreatedAt:   time.Now().UTC(),
	}
	offer.OfferID = canonicalID(offer)
	return offer, nil
}

func normalizeProviderB(raw map[string]interface{}) (*Offer, error) {
	// ProviderB schema example:
	// { "vendor": "ProviderB", "sku": 999, "pricing": { "amount": 12345, "currency_code": "USD", "units": "cents" }, "times": { "leave": 1698772345 } }

	sku := ""
	if v, ok := raw["sku"]; ok {
		switch t := v.(type) {
		case string:
			sku = t
		case float64:
			sku = strconv.FormatInt(int64(t), 10)
		case jsonNumber:
			sku = fmt.Sprintf("%v", t)
		}
	}
	// nested pricing
	price := 0.0
	curr := "USD"
	if p, ok := raw["pricing"].(map[string]interface{}); ok {
		if a, exists := p["amount"]; exists {
			price = parseNumber(a)
			// assume amount in cents if units == "cents"
			if u, ex := p["units"].(string); ex && u == "cents" {
				price = price / 100.0
			}
		}
		if c, exists := p["currency_code"].(string); exists {
			curr = c
		}
	}
	depart := time.Time{}
	if times, ok := raw["times"].(map[string]interface{}); ok {
		if l, ex := times["leave"]; ex {
			// unix timestamp seconds
			if f, ok := l.(float64); ok && f > 1e9 {
				depart = time.Unix(int64(f), 0).UTC()
			}
		}
	}

	offer := &Offer{
		ProviderID:  "ProviderB",
		ProviderRef: sku,
		Price:       price,
		Currency:    curr,
		DepartAt:    depart,
		Raw:         raw,
		CreatedAt:   time.Now().UTC(),
	}
	offer.OfferID = canonicalID(offer)
	return offer, nil
}

func normalizeGeneric(raw map[string]interface{}) (*Offer, error) {
	price := parseNumber(raw["price"])
	curr := ""
	if c, ok := raw["currency"].(string); ok {
		curr = c
	}
	offer := &Offer{
		ProviderID: "generic",
		Price:      price,
		Currency:   curr,
		Raw:        raw,
		CreatedAt:  time.Now().UTC(),
	}
	offer.OfferID = canonicalID(offer)
	return offer, nil
}

func parseNumber(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case string:
		if x, err := strconv.ParseFloat(t, 64); err == nil {
			return x
		}
	}
	return 0.0
}

func parseTime(v interface{}) time.Time {
	switch t := v.(type) {
	case string:
		// try RFC3339
		if tm, err := time.Parse(time.RFC3339, t); err == nil {
			return tm.UTC()
		}
		// try other parse
	case float64:
		if t > 1e9 {
			return time.Unix(int64(t), 0).UTC()
		}
	}
	return time.Time{}
}

func canonicalID(o *Offer) string {
	// deterministic hash of provider + providerRef + price + depart timestamp
	h := sha1.New()
	h.Write([]byte(o.ProviderID))
	h.Write([]byte("|"))
	h.Write([]byte(o.ProviderRef))
	h.Write([]byte("|"))
	h.Write([]byte(fmt.Sprintf("%.2f", o.Price)))
	h.Write([]byte("|"))
	if !o.DepartAt.IsZero() {
		h.Write([]byte(o.DepartAt.UTC().Format(time.RFC3339)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
