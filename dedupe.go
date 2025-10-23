package main

import (
	"sync"
	"time"
)

// Simple in-memory dedupe store with TTL
type DedupeStore struct {
	mu    sync.RWMutex
	data  map[string]*Offer
	ttl   time.Duration
	order []string // for simple eviction
}

func NewDedupeStore() *DedupeStore {
	s := &DedupeStore{
		data: make(map[string]*Offer),
		ttl:  10 * time.Minute, // keep offers for 10 minutes by default
	}
	// start a janitor goroutine to evict old entries
	go s.janitor()
	return s
}

func (s *DedupeStore) janitor() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.data {
			if now.Sub(v.CreatedAt) > s.ttl {
				delete(s.data, k)
			}
		}
		s.mu.Unlock()
	}
}

// IsDuplicate checks if offer's canonical ID already exists
func (s *DedupeStore) IsDuplicate(o *Offer) bool {
	s.mu.RLock()
	_, exists := s.data[o.OfferID]
	s.mu.RUnlock()
	return exists
}

// Add stores canonical offer
func (s *DedupeStore) Add(o *Offer) {
	s.mu.Lock()
	s.data[o.OfferID] = o
	s.mu.Unlock()
}

// List returns all current offers
func (s *DedupeStore) List() []*Offer {
	s.mu.RLock()
	out := make([]*Offer, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	s.mu.RUnlock()
	return out
}
