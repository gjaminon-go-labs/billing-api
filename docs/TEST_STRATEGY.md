# Test Strategy and Organization

## Overview

This document outlines the test strategy for the billing service, including the proper use of different storage backends for different test types.

## Test Types and Storage

### Unit Tests (`tests/unit/`)
- **Purpose**: Test individual components in isolation
- **Storage**: In-memory storage (fast, isolated)
- **Database**: No real database required
- **Bootstrap**: Use `testhelpers.NewUnitTestServer()` or `testhelpers.UnitTestConfig()`
- **Run Command**: `make test-unit` or `go test ./tests/unit/...`

### Integration Tests (`tests/integration/`)
- **Purpose**: Test component interactions and real database behavior
- **Storage**: PostgreSQL storage (real database testing)
- **Database**: Separate test database (`billing_service_test`)
- **Bootstrap**: Use `testhelpers.NewIntegrationTestServer()` or `testhelpers.IntegrationTestConfig()`
- **Run Command**: `make test-integration` or `go test ./tests/integration/...`

## Environment Configuration

### Test Configurations
- **Unit Tests**: Use `di.UnitTestConfig()` - memory storage, no database
- **Integration Tests**: Use `di.IntegrationTestConfig()` - PostgreSQL storage, test database
- **Integration Tests (YAML)**: Use `ENVIRONMENT=test` - loads `configs/test.yaml`

### Deployment Configurations
- **Development**: `configs/development.yaml` - PostgreSQL storage (`billing_service_dev`)
- **Production**: `configs/production.yaml` - PostgreSQL storage (`billing_service`)
- **Base**: `configs/base.yaml` - PostgreSQL storage (default for deployments)

## Database Layout

```
Production:         billing_service      (PostgreSQL)
Development:        billing_service_dev  (PostgreSQL)
Integration Tests:  billing_service_test (PostgreSQL)
Unit Tests:         in-memory            (no real database)
```

## Available Test Helpers

### Unit Test Helpers (Memory Storage)
```go
// Shared singleton container (performance optimized)
server := testhelpers.NewUnitTestServer()
stack := testhelpers.NewUnitTestStack()

// Isolated instances (for tests that modify state)  
server := testhelpers.NewIsolatedUnitTestServer()
stack := testhelpers.NewIsolatedUnitTestStack()
```

### Integration Test Helpers (PostgreSQL Storage)
```go
// PostgreSQL storage with test database
server := testhelpers.NewIntegrationTestServer()
stack := testhelpers.NewIntegrationTestStack()

// Alternative names (backward compatibility)
server := testhelpers.NewPostgresTestServer()
stack := testhelpers.NewPostgresTestStack()
```

### Deprecated Helpers (Backward Compatibility)
```go
// These still work but are deprecated
server := testhelpers.NewInMemoryTestServer()  // → Use NewUnitTestServer()
server := testhelpers.NewIsolatedTestServer()  // → Use NewIsolatedUnitTestServer()
```

## Setup Commands

### Development Setup
```bash
make dev-setup      # Start PostgreSQL + run dev migrations
make run-dev        # Run application in development mode
```

### Test Setup  
```bash
make test-setup     # Start PostgreSQL + create test database + run test migrations
make test-unit      # Run unit tests (memory storage)
make test-integration  # Run integration tests (PostgreSQL test DB)
make test-all       # Run both unit and integration tests (smart: auto-starts PostgreSQL)
make check-postgres # Check PostgreSQL availability, start Docker if needed
```

### Docker PostgreSQL Setup
```bash
make docker-up      # Start PostgreSQL container on port 5433 (avoids local conflicts)
make docker-down    # Stop PostgreSQL container

# Custom port usage
DOCKER_DB_PORT=15432 make docker-up  # Use port 15432 instead of 5433
```

### Database Management
```bash
make migrate-up     # Run migrations on development database
ENVIRONMENT=test make migrate-up  # Run migrations on test database
make migrate-status # Check migration status
```

## Best Practices

### When to Use Unit Tests
- Testing domain logic and business rules
- Testing individual service methods
- Fast feedback loops during development
- No external dependencies needed

### When to Use Integration Tests
- Testing database interactions and queries
- Testing migration behavior
- Testing API endpoints with real storage
- Verifying component integration

### Storage Selection Rules
1. **Never use memory storage in deployment configs** - only PostgreSQL
2. **Unit tests always use memory storage** - for speed and isolation
3. **Integration tests always use PostgreSQL** - for real database testing
4. **Each test type uses separate database** - no cross-contamination

### PostgreSQL Port Configuration
- **Default Docker port**: 5433 (avoids conflicts with local PostgreSQL on 5432)
- **Custom port**: Set `DOCKER_DB_PORT=xxxx` environment variable
- **Integration tests**: Use port 5433 by default (configured in `test.yaml`)
- **CI/CD**: Works with any available port, no conflicts in clean CI environment

### Smart PostgreSQL Detection with Test Isolation
The `make test-all` command automatically manages PostgreSQL for optimal test isolation:

**Scenario 1: Docker PostgreSQL container exists**
```bash
make test-all
# 🔍 Checking PostgreSQL status on port 5433...
# 🔄 Docker PostgreSQL container found - recreating for fresh test isolation...
# 🗑️  Removing existing container...
# 🐘 Creating fresh PostgreSQL container...
# 🧪 Running all tests (unit + integration)...
```

**Scenario 2: Local PostgreSQL running**
```bash
make test-all
# 🔍 Checking PostgreSQL status on port 5433...
# ✅ Local PostgreSQL detected on port 5433 - using existing instance
# 🧪 Running all tests (unit + integration)...
```

**Scenario 3: No PostgreSQL running**
```bash
make test-all
# 🔍 Checking PostgreSQL status on port 5433...
# 🐘 No PostgreSQL found - creating fresh Docker container...
# 🧪 Running all tests (unit + integration)...
```

**Test Isolation Benefits:**
- **Fresh database state**: Each test run starts with clean data
- **No data contamination**: Previous test data completely removed
- **CI/CD optimized**: Perfect for automated testing pipelines
- **Local development friendly**: Preserves user's local PostgreSQL setup

## Migration Testing

### Development Migrations
```bash
make migrate-up      # Apply to billing_service_dev
make migrate-down    # Rollback from billing_service_dev
```

### Test Migrations
```bash
ENVIRONMENT=test go run cmd/migrator/main.go up    # Apply to billing_service_test
ENVIRONMENT=test go run cmd/migrator/main.go down  # Rollback from billing_service_test
```

### Integration Test Auto-Migration
Integration tests automatically run migrations via `IntegrationTestConfig()`:
- `MigrationEnabled: true`
- `MigrationAutoMigrate: true`
- Uses `billing_service_test` database

## File Organization

```
tests/
├── unit/                    # Unit tests (memory storage)
│   ├── domain/             # Domain logic tests
│   └── application/        # Application service tests
├── integration/            # Integration tests (PostgreSQL storage)
│   ├── api/               # API integration tests
│   ├── application/       # Application integration tests
│   └── http/              # HTTP server integration tests
├── testhelpers/           # Test bootstrap helpers
└── testdata/              # External test data files
```

## Configuration Examples

### Unit Test Example
```go
func TestClientCreation(t *testing.T) {
    // Uses memory storage automatically
    server := testhelpers.NewUnitTestServer()
    // ... test logic
}
```

### Integration Test Example  
```go
func TestClientCreationWithDatabase(t *testing.T) {
    // Uses PostgreSQL test database automatically
    server := testhelpers.NewIntegrationTestServer()
    // ... test logic that verifies real database behavior
}
```

This strategy ensures proper test isolation while maintaining realistic database testing where needed.