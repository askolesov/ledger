FROM alpine:latest

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
COPY ledger .

# Switch to non-root user
USER ledger

# Set entrypoint
ENTRYPOINT ["./ledger"]
