# Build stage
FROM golang:1.20-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM alpine:3.18

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S nanorediz && \
    adduser -u 1001 -S nanorediz -G nanorediz

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/ ./bin/

# Copy any additional files needed
COPY --from=builder /app/web/templates ./web/templates/

# Change ownership to non-root user
RUN chown -R nanorediz:nanorediz /app

# Switch to non-root user
USER nanorediz

# Expose ports
EXPOSE 8080 8081 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./bin/nanorediz-client -host localhost -port 8080 ping || exit 1

# Default command
CMD ["./bin/nanorediz-server", "-host", "0.0.0.0", "-port", "8080"]