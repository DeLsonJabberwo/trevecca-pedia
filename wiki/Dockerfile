# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o wiki-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Create data directory for persistent storage
RUN mkdir -p /data && chown -R appuser:appuser /data

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/wiki-service .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set data directory environment variable
ENV WIKI_DATA_DIR=/data

# Expose the port the app runs on
EXPOSE 9454

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9454/health || exit 1

# Run the binary
CMD ["./wiki-service"]
