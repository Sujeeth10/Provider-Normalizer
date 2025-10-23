# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o provider-normalizer ./...

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/provider-normalizer .
EXPOSE 8080
ENTRYPOINT ["./provider-normalizer"]
