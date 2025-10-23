# Provider Normalizer (Go) — PoC

This is a minimal proof-of-concept that demonstrates:
- Accepting provider-specific payloads
- Normalizing into a canonical `Offer` model
- Deterministic deduplication
- In-memory storage and HTTP endpoints
- Unit tests and a Dockerfile

## Files
- `main.go` - HTTP server with `/normalize` and `/offers`
- `normalizer.go` - provider-specific normalization logic
- `dedupe.go` - in-memory dedupe store with TTL
- `normalizer_test.go` - unit tests
- `kafka_publisher.go` - optional Kafka publisher stub (commented)
- `Dockerfile` - build image

## Quick start (local)
```bash
# build
go build -o provider-normalizer

# run
./provider-normalizer

# health
curl http://localhost:8080/health

# example: ProviderA
curl -X POST http://localhost:8080/normalize \
  -H "Content-Type: application/json" \
  -d '{
    "provider_name":"ProviderA",
    "id":"A-123",
    "cost":"123.45",
    "currency":"USD",
    "depart":"2025-11-01T09:00:00Z",
    "class":"economy"
  }'

# list accepted canonical offers
curl http://localhost:8080/offers
