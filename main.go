package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var store *DedupeStore

func main() {
	store = NewDedupeStore()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/normalize", normalizeHandler)
	http.HandleFunc("/offers", offersHandler)

	addr := ":8080"
	log.Printf("Provider Normalizer service starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func normalizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// read body
	var raw map[string]interface{}
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Try provider-specific normalizers. In real system you might pass provider id/header.
	offer, err := Normalize(raw)
	if err != nil {
		http.Error(w, "normalize error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// dedupe check
	if store.IsDuplicate(offer) {
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"status": "duplicate",
			"offer":  offer,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Save & (optionally) publish to Kafka (see kafka_publisher.go)
	store.Add(offer)

	w.WriteHeader(http.StatusCreated)
	resp := map[string]interface{}{
		"status": "accepted",
		"offer":  offer,
	}
	json.NewEncoder(w).Encode(resp)
}

func offersHandler(w http.ResponseWriter, r *http.Request) {
	all := store.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}
