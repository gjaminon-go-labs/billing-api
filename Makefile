# Makefile for billing-service

.DEFAULT_GOAL := help

help:
	@echo "Available commands:"
	@echo "  restore          - Install/update package dependencies"
	@echo "  test-unit        - Run unit tests only (in-memory storage)"
	@echo "  test-integration - Run integration tests only (PostgreSQL test DB)"
	@echo "  test-all         - Run all tests (smart: fresh Docker or existing local PostgreSQL)"
	@echo "  check-postgres   - Check PostgreSQL status, recreate Docker for test isolation"
	@echo "  test-setup       - Set up PostgreSQL test environment"
	@echo "  migrate-up       - Run all pending database migrations"
	@echo "  migrate-down     - Roll back one database migration"
	@echo "  migrate-status   - Show current migration status"
	@echo "  migrate-reset    - Reset database migrations (development only)"
	@echo "  run-dev          - Run application in development mode"
	@echo "  run-prod         - Run application in production mode"
	@echo ""
	@echo "Docker PostgreSQL commands:"
	@echo "  docker-up        - Start PostgreSQL container (port 5433)"
	@echo "  docker-down      - Stop PostgreSQL container"
	@echo "  recreate-docker-postgres - Recreate container for fresh state"
	@echo "  dev-setup        - Docker + migrations for development"
	@echo ""
	@echo "Environment variables:"
	@echo "  DOCKER_DB_PORT   - PostgreSQL port for Docker (default: 5433)"

restore:
	go mod tidy

test-unit:
	@echo "Running unit tests (domain layer only)..."
	go test -v ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Check PostgreSQL availability and ensure fresh Docker container for tests
check-postgres:
	@echo "🔍 Checking PostgreSQL status on port $(DOCKER_DB_PORT)..."
	@if command -v docker >/dev/null 2>&1; then \
		if docker ps --format "{{.Names}}" 2>/dev/null | grep -q "^billing-postgres$$"; then \
			echo "🔄 Docker PostgreSQL container found - recreating for fresh test isolation..."; \
			$(MAKE) recreate-docker-postgres; \
		elif command -v nc >/dev/null 2>&1 && nc -z localhost $(DOCKER_DB_PORT) 2>/dev/null; then \
			echo "✅ Local PostgreSQL detected on port $(DOCKER_DB_PORT) - using existing instance"; \
		elif command -v timeout >/dev/null 2>&1 && timeout 1 bash -c "</dev/tcp/localhost/$(DOCKER_DB_PORT)" 2>/dev/null; then \
			echo "✅ Local PostgreSQL detected on port $(DOCKER_DB_PORT) - using existing instance"; \
		else \
			echo "🐘 No PostgreSQL found - creating fresh Docker container..."; \
			$(MAKE) docker-up; \
		fi \
	else \
		echo "⚠️  Docker not available. Please ensure PostgreSQL is running on port $(DOCKER_DB_PORT)"; \
		echo "   Or install Docker to use automated PostgreSQL setup"; \
	fi

test-all: check-postgres
	@echo "🧪 Running all tests (unit + integration)..."
	go test -v ./tests/unit/... ./tests/integration/...

# Migration commands
migrate-up:
	@echo "Running all pending database migrations..."
	go run cmd/migrator/main.go up

migrate-down:
	@echo "Rolling back one database migration..."
	go run cmd/migrator/main.go down

migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrator/main.go status

migrate-reset:
	@echo "⚠️  WARNING: This will reset all migrations (development only)"
	@echo "This command should only be used in development environments"
	go run cmd/migrator/main.go force 0
	go run cmd/migrator/main.go up

# Application commands  
run-dev:
	@echo "Starting application in development mode..."
	ENVIRONMENT=development go run cmd/api/main.go

run-prod:
	@echo "Starting application in production mode..."
	ENVIRONMENT=production go run cmd/api/main.go

# Build commands
build:
	@echo "Building application binaries..."
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrator cmd/migrator/main.go

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Docker commands (when PostgreSQL is available)
# Set DOCKER_DB_PORT environment variable to use a different port (default: 5433)
DOCKER_DB_PORT ?= 5433

docker-up:
	@echo "Starting PostgreSQL with Docker on port $(DOCKER_DB_PORT) (avoids conflicts with local PostgreSQL)..."
	docker run --name billing-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=billing_service_dev -p $(DOCKER_DB_PORT):5432 -d postgres:15
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 5
	@echo "Creating test database..."
	docker exec billing-postgres psql -U postgres -c "CREATE DATABASE billing_service_test;" || true

docker-down:
	@echo "Stopping PostgreSQL container..."
	docker stop billing-postgres || true
	docker rm billing-postgres || true

# Recreate Docker PostgreSQL for fresh state (integration test isolation)
recreate-docker-postgres:
	@echo "🔄 Recreating Docker PostgreSQL for fresh database state..."
	@echo "🗑️  Removing existing container..."
	$(MAKE) docker-down
	@echo "🐘 Creating fresh PostgreSQL container..."
	$(MAKE) docker-up
	@echo "✅ Fresh PostgreSQL container ready for testing"

# Test-specific commands
test-setup: docker-up
	@echo "Setting up test databases..."
	$(MAKE) migrate-up
	ENVIRONMENT=test go run cmd/migrator/main.go up
	@echo "Test environment ready!"

# Development workflow
dev-setup: docker-up
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 5
	$(MAKE) migrate-up
	@echo "Development environment ready!"

dev-teardown: docker-down clean
	@echo "Development environment cleaned up!"

.PHONY: help restore test-unit test-integration test-all check-postgres recreate-docker-postgres migrate-up migrate-down migrate-status migrate-reset run-dev run-prod build clean docker-up docker-down dev-setup dev-teardown test-setup