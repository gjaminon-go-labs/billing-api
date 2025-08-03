# Integration Tests Explained

## Overview

This document explains the **3 types of integration tests** implemented in our Go DDD billing service. Each type tests different levels of component interaction, from business logic orchestration to complete HTTP server behavior.

## Test Classification

- **Unit Tests**: Test single components in isolation (Domain layer only)
- **Integration Tests**: Test multiple components working together
- **End-to-End Tests**: Test complete user workflows (not yet implemented)

## Integration Test Types

### 1. Application Layer Integration Tests

**Location:** `tests/integration/application/client_operations_test.go`

**Components Integrated:**
- Application Service (`BillingService`)
- Repository Layer (`ClientRepository`) 
- Storage Infrastructure (`InMemoryStorage`)

**What it tests:**
- Business orchestration through the application service
- Data persistence through repository abstraction
- Domain validation enforcement at service level
- Complete Create Client workflow: Service → Repository → Storage

**Integration Scope:** Application + Repository + Storage (no HTTP layer)

**Example Test:**
```go
func TestBillingService_CreateClient(t *testing.T) {
    // Set up dependencies with in-memory storage
    storage := infrastructure.NewInMemoryStorage()
    clientRepo := repository.NewClientRepository(storage)
    service := application.NewBillingService(clientRepo)

    // Test service orchestration
    client, err := service.CreateClient(name, email, phone, address)
    // Assertions...
}
```

### 2. API Handler Integration Tests

**Location:** `tests/integration/api/client_handler_test.go`

**Components Integrated:**
- HTTP Handler (`ClientHandler`)
- Application Service (`BillingService`)
- Repository Layer (`ClientRepository`)
- Storage Infrastructure (`InMemoryStorage`)

**What it tests:**
- HTTP request/response handling and transformation
- DTO conversion (HTTP JSON ↔ Domain models)
- Error handling and HTTP status code mapping
- Handler → Service → Repository → Storage workflow

**Integration Scope:** API + Application + Repository + Storage (no HTTP server/routing)

**Example Test:**
```go
func TestClientHandler_CreateClient(t *testing.T) {
    // Set up complete handler stack
    storage := infrastructure.NewInMemoryStorage()
    clientRepo := repository.NewClientRepository(storage)
    billingService := application.NewBillingService(clientRepo)
    handler := handlers.NewClientHandler(billingService)

    // Create HTTP request
    req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", requestBody)
    rr := httptest.NewRecorder()

    // Test handler directly
    handler.CreateClient(rr, req)
    // Assertions on HTTP response...
}
```

### 3. Full HTTP Server Integration Tests

**Location:** `tests/integration/http/server_integration_test.go`

**Components Integrated:**
- Complete HTTP Server with routing
- Middleware (CORS, logging, error handling)
- HTTP Handler (`ClientHandler`)
- Application Service (`BillingService`) 
- Repository Layer (`ClientRepository`)
- Storage Infrastructure (`InMemoryStorage`)
- Test Helpers (`NewInMemoryTestServer()`)

**What it tests:**
- End-to-end HTTP requests (real network calls)
- Complete server routing and middleware stack
- CORS functionality
- Health check endpoints
- Request persistence across multiple HTTP calls
- Full production-like server behavior

**Integration Scope:** Complete HTTP stack from network → handler → service → repository → storage

**Example Test:**
```go
func TestHTTPServer_Integration_CreateClient(t *testing.T) {
    // Set up complete HTTP server using test helpers
    server := testhelpers.NewInMemoryTestServer()
    testServer := httptest.NewServer(server.Handler())
    defer testServer.Close()

    // Make actual HTTP request to test server
    resp, err := http.Post(testServer.URL+"/api/v1/clients", "application/json", requestBody)
    // Assertions on real HTTP response...
}
```

## Integration Test Hierarchy

```
Unit Tests (Domain Only)
├── Domain Models (Client, Email, etc.)
└── Domain Errors (ValidationError, etc.)
    ↓
Application Integration (Service + Repository + Storage)
├── Business orchestration
├── Data persistence
└── Domain validation enforcement
    ↓  
API Integration (Handler + Service + Repository + Storage)
├── HTTP request/response handling
├── DTO conversion
└── Error mapping to HTTP status codes
    ↓
HTTP Server Integration (Full HTTP Stack + All Components)
├── End-to-end HTTP requests
├── Routing and middleware
├── CORS functionality
└── Multi-request persistence
```

## Key Benefits

### Progressive Testing
Each level adds more components and tests broader integration scenarios:

1. **Application Level**: Validates business logic orchestration
2. **API Level**: Ensures proper HTTP contract implementation  
3. **Server Level**: Confirms production-like behavior

### Isolation of Concerns
- **Application tests** focus on business logic without HTTP concerns
- **API tests** focus on HTTP handling without server infrastructure
- **Server tests** validate complete system behavior

### Fast Feedback Loop
- Unit tests run fastest (domain only)
- Application integration tests are medium speed
- API integration tests add HTTP processing overhead
- Server integration tests are slowest (full network stack)

## Test Data Strategy

All integration tests use **external JSON test data** following Kubernetes patterns:

```
tests/
├── testdata/
│   ├── client/
│   │   └── client_test_cases.json      # Shared client test scenarios
│   └── http/
│       └── create_client_requests.json # HTTP-specific test cases
```

### Benefits:
- **Reusable**: Same test scenarios across different integration levels
- **Maintainable**: Update test cases in one place
- **Readable**: JSON format is easy to understand and modify
- **Version Controlled**: Test data evolution is tracked

## Test Execution

```bash
# Run only unit tests (domain layer)
make test-unit

# Run only integration tests (all levels)
make test-integration

# Run complete test suite
make test-all
```

## Best Practices Learned

### Test Classification
- **Unit**: Single component, no external dependencies
- **Integration**: Multiple components, shared infrastructure
- **End-to-End**: Complete user workflows (future implementation)

### Path Resolution
- Use `runtime.Caller(0)` for reliable test data file loading
- Adjust relative paths based on test file location in directory structure

### Dependency Management
- Use test bootstraps to separate test infrastructure from production code
- Inject dependencies through constructors for testability

### Error Testing
- Test both success and failure scenarios at each integration level
- Verify error propagation through the complete stack

---

*This document reflects the learning journey of implementing proper integration testing in a Go DDD architecture, following clean architecture principles and Go community best practices.*