# Simplified Testing Strategy - "Bring Your Own PostgreSQL"

## Overview

This document outlines the simplified testing strategy for the billing-api service, focusing on clear separation of concerns and predictable developer experience.

## Core Principles

### 1. Clear Storage Type Separation
- **Unit Tests**: Always use in-memory storage
- **Integration Tests**: Always use local PostgreSQL
- **No Exceptions**: Storage type is determined by test type, not configuration

### 2. Developer Responsibility
- Developers manage their own PostgreSQL installation
- No Docker complexity in the development workflow
- Clear setup requirements and error messages

### 3. Simplified Tooling
- Minimal Makefile with essential commands only
- No smart detection or auto-provisioning logic
- Predictable behavior across all environments

## Test Types

### Unit Tests (`tests/unit/`)
**Purpose**: Test individual components in isolation
- **Storage**: In-memory storage only
- **Database**: No real database required
- **Speed**: Very fast (< 1 second)
- **Command**: `make test-unit`

**Characteristics:**
- Test domain logic and business rules
- Test individual service methods
- No external dependencies
- Perfect for TDD and rapid feedback

### Integration Tests (`tests/integration/`)
**Purpose**: Test component interactions with real database
- **Storage**: Local PostgreSQL only
- **Database**: `billing_service_test` on localhost:5432
- **Speed**: Moderate (depends on test complexity)
- **Command**: `make test-integration`

**Characteristics:**
- Test database interactions and queries
- Test API endpoints with real storage
- Verify component integration
- Test migration behavior

## Developer Setup Requirements

### Prerequisites
1. **PostgreSQL 12+ installed locally**
   - Running on standard port 5432
   - Accessible with standard credentials

2. **Required Databases**
   ```sql
   CREATE DATABASE billing_service_dev;
   CREATE DATABASE billing_service_test;
   ```

3. **Environment Setup**
   ```bash
   # Clone and setup
   cd billing-api
   make restore
   
   # Run migrations
   make migrate-up                    # Development database
   ENVIRONMENT=test make migrate-up   # Test database
   ```

### Quick Setup Verification
```bash
# Test PostgreSQL connectivity
psql -h localhost -p 5432 -U postgres -l

# Verify databases exist
psql -h localhost -p 5432 -U postgres -c "\l" | grep billing_service

# Run tests to verify setup
make test-unit         # Should always work (no PostgreSQL needed)
make test-integration  # Should work if PostgreSQL setup correct
```

## Daily Development Workflow

### Standard Commands
```bash
make test-unit         # Fast unit tests (memory storage)
make test-integration  # Integration tests (requires local PostgreSQL)
make test-all          # Run both unit and integration tests
make run-dev           # Start development server
make migrate-up        # Apply pending migrations
```

### Error Handling
If PostgreSQL is not running or databases don't exist:
```bash
$ make test-integration
Running integration tests...
Error: Failed to connect to PostgreSQL at localhost:5432
Please ensure:
1. PostgreSQL is running locally
2. Database 'billing_service_test' exists
3. Connection credentials are correct
```

## What We Removed

### Docker Complexity
- **Removed**: `docker-up`, `docker-down`, `recreate-docker-postgres`
- **Removed**: Smart PostgreSQL detection logic
- **Removed**: Port conflict management (5433 vs 5432)
- **Removed**: Container recreation for test isolation

### Smart Detection Logic
- **Removed**: ~40 lines of complex shell logic in Makefile
- **Removed**: Multiple PostgreSQL scenario handling
- **Removed**: Automatic Docker container management
- **Removed**: CI/CD specific PostgreSQL provisioning

### Configuration Complexity
- **Removed**: Multiple storage backend options for integration tests
- **Removed**: Docker-specific environment variables
- **Removed**: Port configuration complexity

## What We Kept Simple

### Essential Makefile Commands
```makefile
test-unit:
    go test -v ./tests/unit/...

test-integration:
    go test -v ./tests/integration/...

test-all: test-unit test-integration

migrate-up:
    go run cmd/migrator/main.go up

run-dev:
    ENVIRONMENT=development go run cmd/api/main.go
```

### Clear Test Helpers
```go
// Unit tests - always memory storage
server := testhelpers.NewUnitTestServer()

// Integration tests - always PostgreSQL
server := testhelpers.NewIntegrationTestServer()
```

## Migration Strategy

### Development Migrations
```bash
make migrate-up      # Apply to billing_service_dev
make migrate-down    # Rollback from billing_service_dev
make migrate-status  # Check migration status
```

### Test Migrations
```bash
ENVIRONMENT=test make migrate-up    # Apply to billing_service_test
ENVIRONMENT=test make migrate-down  # Rollback from billing_service_test
```

### Auto-Migration for Tests
Integration tests automatically run migrations:
- `MigrationEnabled: true`
- `MigrationAutoMigrate: true`
- Uses `billing_service_test` database

## Benefits of This Approach

### For Developers
- **Predictable**: Always know where PostgreSQL is (localhost:5432)
- **Simple**: No Docker complexity to understand
- **Fast**: Local PostgreSQL typically faster than Docker
- **Realistic**: Most developers already have PostgreSQL installed

### For Project Maintenance
- **Smaller Makefile**: ~60 lines instead of 145
- **Less Complexity**: No smart detection logic to maintain
- **Clear Separation**: Development vs CI/CD concerns
- **Easier Debugging**: Fewer moving parts

### for Testing
- **Clear Rules**: Unit = memory, Integration = PostgreSQL
- **No Confusion**: Storage type determined by test type
- **Fast Feedback**: Unit tests remain very fast
- **Real Testing**: Integration tests use real database

## Error Scenarios and Solutions

### PostgreSQL Not Running
**Error**: `connection refused to localhost:5432`

**Solution**:
```bash
# Start PostgreSQL service
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS

# Or install PostgreSQL if not present
# See PostgreSQL installation docs for your OS
```

### Database Doesn't Exist
**Error**: `database "billing_service_test" does not exist`

**Solution**:
```bash
# Create missing databases
psql -U postgres -c "CREATE DATABASE billing_service_dev;"
psql -U postgres -c "CREATE DATABASE billing_service_test;"

# Run migrations
make migrate-up
ENVIRONMENT=test make migrate-up
```

### Migration Failures
**Error**: `migration failed`

**Solution**:
```bash
# Check migration status
make migrate-status

# Check database connection
psql -U postgres -d billing_service_dev -c "SELECT version();"

# Reset migrations (development only)
make migrate-reset
```

## Comparison: Before vs After

### Before (Complex)
- 145-line Makefile with smart detection
- Docker container recreation logic
- Port conflict management
- CI/CD and development mixed
- Integration tests could use memory storage
- Multiple PostgreSQL scenarios to handle

### After (Simple)
- ~60-line Makefile with essential commands
- No Docker management
- Standard PostgreSQL port (5432)
- Clear separation of concerns
- Integration tests always use PostgreSQL
- One PostgreSQL scenario: localhost:5432

## FAQ

### Q: What if I prefer Docker for PostgreSQL?
**A**: You can still use Docker for your local PostgreSQL, just run:
```bash
docker run --name my-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15
```
The key is that it's your choice, not managed by the project.

### Q: How do we handle CI/CD environments?
**A**: CI/CD will have its own PostgreSQL provisioning logic in the infrastructure project. This keeps CI/CD complexity separate from development workflow.

### Q: What about test isolation?
**A**: Integration tests use the `billing_service_test` database with auto-migration. Each test run gets a clean migrated state.

### Q: Performance impact of removing Docker recreation?
**A**: Local PostgreSQL is typically faster than Docker. Test isolation is maintained through the dedicated test database and auto-migration.

## Migration Path

### Phase 1: Update Documentation ✅
- Create this simplified testing strategy document

### Phase 2: Simplify Makefile
- Remove Docker commands and smart detection logic
- Keep only essential commands

### Phase 3: Enforce Storage Types
- Update test helpers to enforce PostgreSQL for integration tests
- Remove deprecated/confusing helper methods

### Phase 4: Update Main Documentation
- Update README.md with new approach
- Replace old TEST_STRATEGY.md with this document

---

**Result**: A clean, predictable, and maintainable testing strategy that puts complexity where it belongs and keeps the development experience simple and fast.