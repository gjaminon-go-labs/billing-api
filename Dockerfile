# Multi-stage build for optimal size
FROM docker.io/golang:1.24-alpine AS builder

# Build argument for version injection
ARG VERSION=dev
ARG BUILD_DATE
ARG GIT_COMMIT

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT}" \
    -o billing-api cmd/api/main.go

# Final stage - minimal alpine image
FROM docker.io/alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Copy binary from builder
COPY --from=builder /app/billing-api /billing-api

# Copy migrations if needed in container
COPY --from=builder /app/database/migrations /database/migrations

# Copy configuration files
COPY --from=builder /app/configs /configs

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/billing-api"]