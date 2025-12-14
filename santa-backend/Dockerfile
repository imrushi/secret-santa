# --- STAGE 1: Build the Binary ---
# We use Alpine for the builder because it's fast and small
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install SSL certificates (Needed if your app makes HTTPS calls)
RUN apk update && apk add --no-cache git ca-certificates

# Cache dependencies first (better build speed)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0: Disables C libraries to create a static binary (required for 'scratch')
# GOOS=linux: Targets Linux
# GOARCH=arm64: Targets Raspberry Pi 3, 4, 5 (64-bit OS). 
# -ldflags="-s -w": Strips debug info to reduce binary size by ~30%
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o santa-server main.go

# --- STAGE 2: The Lightweight Runtime ---
# 'scratch' is an empty image (0 MB)
FROM scratch

# Copy the SSL root certificates from the builder stage
# (Critical: without this, you cannot make external HTTPS requests)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=builder /app/santa-server /santa-server

# Expose your application port
EXPOSE 8007

# Run the binary
ENTRYPOINT ["/santa-server"]