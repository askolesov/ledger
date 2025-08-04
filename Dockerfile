# =============================================================================
# Builder Stage
# =============================================================================
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Set working directory
WORKDIR /build

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -a \
        -installsuffix cgo \
        -ldflags="-w -s" \
        -o ledger \
        ./cmd/ledger

# =============================================================================
# Runner Stage
# =============================================================================
FROM alpine:latest AS runner

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user and group
RUN addgroup -g 1001 -S ledger && \
    adduser -u 1001 -S ledger -G ledger

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=ledger:ledger /build/ledger .

# Switch to non-root user
USER ledger

# Set entrypoint
ENTRYPOINT ["./ledger"]
