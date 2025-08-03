# Makefile for billing-api

.DEFAULT_GOAL := help

help:
	@echo "Available commands:"
	@echo "  dev-setup        - Setup development environment (provision databases)"
	@echo "  test-setup       - Setup test environment (provision test database)"
	@echo "  restore          - Install/update package dependencies"
	@echo "  test-unit        - Run unit tests only (domain layer validation)"
	@echo "  test-integration - Run integration tests only (requires local PostgreSQL)"
	@echo "  test-all         - Run all tests (unit + integration)"
	@echo "  migrate-up       - Run all pending database migrations (dev environment)"
	@echo "  migrate-down     - Roll back one database migration (dev environment)"
	@echo "  migrate-status   - Show current migration status (dev environment)"
	@echo "  migrate-reset    - Reset database migrations (development only)"
	@echo "  run-dev          - Run application in development mode"
	@echo "  build            - Build application binaries"
	@echo "  clean            - Clean build artifacts"
	@echo "  validate-env     - Validate environment setup (databases, infrastructure)"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - PostgreSQL running on localhost:5432"
	@echo "  - Databases: go-labs-dev, go-labs-tst (use 'make dev-setup' to provision)"
	@echo "  - Infrastructure: ../infrastructure/scripts/provision-database.sh"

# Infrastructure setup commands
dev-setup:
	@echo "Setting up development environment..."
	@echo "1. Provisioning development database (go-labs-dev)..."
	@if [ ! -f "../infrastructure/scripts/provision-database.sh" ]; then \
		echo "❌ Error: Infrastructure script not found at ../infrastructure/scripts/provision-database.sh"; \
		echo "   Please ensure you're in the billing-api directory and infrastructure project exists"; \
		exit 1; \
	fi
	@cd ../infrastructure && ./scripts/provision-database.sh dev
	@echo "2. Running development database migrations..."
	$(MAKE) migrate-up
	@echo "✅ Development environment setup complete!"
	@echo ""
	@echo "You can now:"
	@echo "  - Run tests: make test-all"
	@echo "  - Start API: make run-dev"

test-setup:
	@echo "Setting up test environment..."
	@echo "1. Provisioning test database (go-labs-tst)..."
	@if [ ! -f "../infrastructure/scripts/provision-database.sh" ]; then \
		echo "❌ Error: Infrastructure script not found at ../infrastructure/scripts/provision-database.sh"; \
		echo "   Please ensure you're in the billing-api directory and infrastructure project exists"; \
		exit 1; \
	fi
	@cd ../infrastructure && ./scripts/provision-database.sh tst
	@echo "2. Test database setup complete!"
	@echo "✅ Test environment ready (migrations run automatically during integration tests)"

restore:
	go mod tidy

test-unit:
	@echo "Running unit tests (domain layer validation)..."
	go test -v ./tests/unit/...

test-integration:
	@echo "Running integration tests (requires local PostgreSQL)..."
	@echo "Checking PostgreSQL connectivity..."
	@if ! command -v psql >/dev/null 2>&1; then \
		echo "❌ Error: psql command not found. Please install PostgreSQL client."; \
		exit 1; \
	fi
	@if ! PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d go-labs-tst -c "SELECT 1;" >/dev/null 2>&1; then \
		echo "❌ Error: Cannot connect to PostgreSQL at localhost:5432 or database 'go-labs-tst' does not exist"; \
		echo "   Please ensure:"; \
		echo "   1. PostgreSQL is running locally"; \
		echo "   2. Database 'go-labs-tst' exists (run 'make test-setup' to create)"; \
		echo "   3. Connection credentials are correct (postgres/postgres)"; \
		echo "   4. Test database migrations are up to date (run 'make migrate-up-test')"; \
		exit 1; \
	fi
	go test -v ./tests/integration/...

test-all:
	@echo "Running all tests (unit + integration)..."
	$(MAKE) test-unit
	$(MAKE) test-integration

# Migration commands (default to development environment)
migrate-up:
	@echo "Running all pending database migrations (development)..."
	ENVIRONMENT=development go run cmd/migrator/main.go up

migrate-down:
	@echo "Rolling back one database migration (development)..."
	ENVIRONMENT=development go run cmd/migrator/main.go down

migrate-status:
	@echo "Checking migration status (development)..."
	ENVIRONMENT=development go run cmd/migrator/main.go status

migrate-reset:
	@echo "⚠️  WARNING: This will reset all migrations (development only)"
	@echo "This command should only be used in development environments"
	ENVIRONMENT=development go run cmd/migrator/main.go force 0
	ENVIRONMENT=development go run cmd/migrator/main.go up

# Application commands  
run-dev:
	@echo "Starting application in development mode..."
	ENVIRONMENT=development go run cmd/api/main.go

# Build commands
build:
	@echo "Building application binaries..."
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrator cmd/migrator/main.go

# Validation and utility commands
validate-env:
	@echo "Validating environment setup..."
	@echo "1. Checking PostgreSQL connectivity..."
	@if ! command -v psql >/dev/null 2>&1; then \
		echo "❌ PostgreSQL client not found"; \
		exit 1; \
	fi
	@if ! PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -c "SELECT version();" >/dev/null 2>&1; then \
		echo "❌ Cannot connect to PostgreSQL"; \
		exit 1; \
	fi
	@echo "✅ PostgreSQL connectivity OK"
	@echo "2. Checking development database..."
	@if ! PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d go-labs-dev -c "SELECT 1;" >/dev/null 2>&1; then \
		echo "❌ Development database (go-labs-dev) not found - run 'make dev-setup'"; \
		exit 1; \
	fi
	@echo "✅ Development database OK"
	@echo "3. Checking test database..."
	@if ! PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d go-labs-tst -c "SELECT 1;" >/dev/null 2>&1; then \
		echo "❌ Test database (go-labs-tst) not found - run 'make test-setup'"; \
		exit 1; \
	fi
	@echo "✅ Test database OK"
	@echo "4. Checking infrastructure scripts..."
	@if [ ! -f "../infrastructure/scripts/provision-database.sh" ]; then \
		echo "❌ Infrastructure script not found"; \
		exit 1; \
	fi
	@echo "✅ Infrastructure scripts OK"
	@echo "🎉 Environment validation complete - all systems ready!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

.PHONY: help dev-setup test-setup restore test-unit test-integration test-all migrate-up migrate-down migrate-status migrate-reset run-dev build clean validate-env