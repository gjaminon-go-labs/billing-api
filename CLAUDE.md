# Billing API - Service-Specific Context

## Service Overview
This is the billing microservice for the go-labs platform. It handles client management, invoicing, and payment operations following Domain-Driven Design principles.

⚠️ **IMPORTANT**: All development MUST follow the TDD workflow defined in `../CLAUDE.md`

## Current Implementation Status

### ✅ Completed Features
- **UC-B-001**: Create Client - Full implementation with validation
  - Domain entity with email/phone value objects
  - Repository pattern with PostgreSQL storage
  - REST API endpoint: `POST /api/v1/clients`
  - Comprehensive test coverage

### 🎯 Next Implementation Priority
1. **UC-B-002**: Get Client by ID
2. **UC-B-003**: Update Client
3. **UC-B-004**: Delete Client  
4. **UC-B-005**: List Clients with pagination

## Domain-Specific Rules

### Client Domain
- Email addresses MUST be unique (enforced at domain level)
- Phone numbers use international format with validation
- Client IDs are UUIDs generated on creation
- Soft delete NOT supported (hard delete only)

### Error Handling Pattern
```go
// Domain errors (internal/domain/errors/)
var (
    ErrClientNotFound = NewNotFoundError("client", "client not found")
    ErrEmailExists = NewConflictError("email", "email already exists")
)

// Always wrap domain errors properly
if err == gorm.ErrRecordNotFound {
    return nil, errors.ErrClientNotFound
}
```

### Repository Pattern
- Interfaces in `internal/domain/repository/`
- Implementations in `internal/infrastructure/repository/`
- Always program against interfaces, not concrete types

### Testing Patterns
- Unit tests: Same package, in-memory storage
- Integration tests: Separate `_test` package, PostgreSQL
- Test data in `tests/testdata/` as JSON files
- Use table-driven tests with `t.Run()`

## API Conventions

### RESTful Endpoints
```
POST   /api/v1/clients          # Create
GET    /api/v1/clients/:id      # Get by ID  
PUT    /api/v1/clients/:id      # Update
DELETE /api/v1/clients/:id      # Delete
GET    /api/v1/clients          # List (with pagination)
```

### Response Format
```json
// Success
{
  "data": { ... },
  "status": "success"
}

// Error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": { ... }
  }
}
```

## Package Structure Reminder
```
internal/
├── domain/          # Business logic (entities, value objects)
├── application/     # Use cases (service layer)
├── infrastructure/  # External concerns (DB, HTTP)
└── api/            # HTTP handlers and routing
```

## Quick Commands
```bash
# Run tests with coverage
make test-all

# Run only domain tests
go test ./internal/domain/...

# Run integration tests
make test-integration

# Start development server
make run-dev
```

## Important Notes
1. This service uses schema-based isolation (`billing` schema)
2. Two database users: migration user (DDL) and app user (DML)
3. Test database has automatic cleanup between runs
4. All new features MUST follow TDD workflow from `../CLAUDE.md`

---
*For TDD workflow and general project rules, see `../CLAUDE.md`*