// Test Helpers with Dependency Injection
//
// This file provides test helpers with clear separation of storage types:
// - Unit Tests: Always use in-memory storage (NewUnitTestServer, NewUnitTestStack)
// - Integration Tests: Always use PostgreSQL storage (NewIntegrationTestServer, NewIntegrationTestStack)
//
// IMPORTANT: Integration tests require local PostgreSQL running on localhost:5432
// with databases: billing_service_dev, billing_service_test
//
// Benefits: Memory efficient, thread-safe, clear storage separation, predictable behavior
// Pattern: Uses DI container with lazy initialization and storage type enforcement
package testhelpers

import (
	"fmt"

	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/di"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
)

// TestContainer provides a shared DI container for unit tests (singleton pattern)
var unitTestContainer *di.Container

// GetUnitTestContainer returns the shared unit test container, creating it if necessary
func GetUnitTestContainer() *di.Container {
	if unitTestContainer == nil {
		unitTestContainer = di.NewContainer(di.UnitTestConfig())
	}
	return unitTestContainer
}

// GetTestContainer returns the shared unit test container (backward compatibility)
// Deprecated: Use GetUnitTestContainer() instead
func GetTestContainer() *di.Container {
	return GetUnitTestContainer()
}

// ResetUnitTestContainer clears the shared unit test container (useful for test isolation)
func ResetUnitTestContainer() {
	if unitTestContainer != nil {
		unitTestContainer.Reset()
	}
	unitTestContainer = nil
}

// ResetTestContainer clears the shared unit test container (backward compatibility)
// Deprecated: Use ResetUnitTestContainer() instead
func ResetTestContainer() {
	ResetUnitTestContainer()
}

// NewUnitTestServer creates an HTTP server using in-memory storage for unit tests
// Performance: Uses singleton services, lazy initialization, no duplicate instances
func NewUnitTestServer() *httpserver.Server {
	container := GetUnitTestContainer()
	server, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to create unit test server: " + err.Error())
	}
	return server
}

// NewInMemoryTestServer creates an HTTP server using in-memory storage (backward compatibility)
// Deprecated: Use NewUnitTestServer() instead
func NewInMemoryTestServer() *httpserver.Server {
	return NewUnitTestServer()
}

// NewUnitTestContainer creates a DI container configured for unit testing (memory storage)
func NewUnitTestContainer() *di.Container {
	return di.NewContainer(di.UnitTestConfig())
}

// NewInMemoryTestContainer creates a DI container configured for in-memory testing (backward compatibility)
// Deprecated: Use NewUnitTestContainer() instead
func NewInMemoryTestContainer() *di.Container {
	return NewUnitTestContainer()
}

// NewIntegrationTestContainer creates a DI container configured for integration testing (PostgreSQL)
func NewIntegrationTestContainer() *di.Container {
	return di.NewContainer(di.IntegrationTestConfig())
}

// NewIntegrationTestServer creates an HTTP server using PostgreSQL for integration tests
func NewIntegrationTestServer() *httpserver.Server {
	container := NewIntegrationTestContainer()
	config := container.GetConfig()

	// Clean up test data if enabled in configuration
	if config.TestCleanupEnabled && config.TestCleanupOnSetup {
		stack := createTestStack(container)
		if stack.DatabaseCleaner != nil {
			if err := stack.DatabaseCleaner.CleanupTestData(); err != nil {
				panic("Failed to cleanup test data: " + err.Error())
			}
		}
	}

	server, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to create integration test server: " + err.Error())
	}
	return server
}

// NewIsolatedUnitTestServer creates an HTTP server with isolated dependencies for unit tests
// Use this when you need fresh instances for each test
func NewIsolatedUnitTestServer() *httpserver.Server {
	container := NewUnitTestContainer()
	server, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to create isolated unit test server: " + err.Error())
	}
	return server
}

// NewIsolatedTestServer creates an HTTP server with isolated dependencies (backward compatibility)
// Deprecated: Use NewIsolatedUnitTestServer() instead
func NewIsolatedTestServer() *httpserver.Server {
	return NewIsolatedUnitTestServer()
}

// NewIsolatedIntegrationTestServer creates an HTTP server with isolated dependencies for integration tests
// Use this when you need fresh instances for each integration test
func NewIsolatedIntegrationTestServer() *httpserver.Server {
	container := NewIntegrationTestContainer()
	server, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to create isolated integration test server: " + err.Error())
	}
	return server
}

// TestStack provides access to all components in the DI container
type TestStack struct {
	Container       *di.Container
	Storage         storage.Storage
	ClientRepo      repository.ClientRepository
	BillingService  *application.BillingService
	HTTPServer      *httpserver.Server
	DatabaseCleaner *DatabaseCleaner // Added for test data cleanup
}

// NewUnitTestStack creates a complete unit test stack using DI container
func NewUnitTestStack() *TestStack {
	container := GetUnitTestContainer()
	return createTestStack(container)
}

// NewInMemoryTestStack creates a complete test stack using DI container (backward compatibility)
// Deprecated: Use NewUnitTestStack() instead
func NewInMemoryTestStack() *TestStack {
	return NewUnitTestStack()
}

// NewIsolatedUnitTestStack creates a unit test stack with isolated dependencies
func NewIsolatedUnitTestStack() *TestStack {
	container := NewUnitTestContainer()
	return createTestStack(container)
}

// NewIsolatedTestStack creates a test stack with isolated dependencies (backward compatibility)
// Deprecated: Use NewIsolatedUnitTestStack() instead
func NewIsolatedTestStack() *TestStack {
	return NewIsolatedUnitTestStack()
}

// NewIntegrationTestStack creates a complete integration test stack using PostgreSQL
func NewIntegrationTestStack() *TestStack {
	container := NewIntegrationTestContainer()
	config := container.GetConfig()

	stack := createTestStack(container)

	// Clean up test data if enabled in configuration
	if config.TestCleanupEnabled && config.TestCleanupOnSetup {
		if stack.DatabaseCleaner != nil {
			if err := stack.DatabaseCleaner.CleanupTestData(); err != nil {
				panic("Failed to cleanup test data: " + err.Error())
			}
		}
	}

	return stack
}

// createTestStack is a helper function to create a TestStack from a container
func createTestStack(container *di.Container) *TestStack {
	// Get all components from container (lazy initialization)
	stor, err := container.GetStorage()
	if err != nil {
		panic("Failed to get storage: " + err.Error())
	}

	clientRepo, err := container.GetClientRepository()
	if err != nil {
		panic("Failed to get client repository: " + err.Error())
	}

	billingService, err := container.GetBillingService()
	if err != nil {
		panic("Failed to get billing service: " + err.Error())
	}

	httpServer, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to get HTTP server: " + err.Error())
	}

	// Create database cleaner if using PostgreSQL storage
	var dbCleaner *DatabaseCleaner
	if postgresStorage, ok := stor.(*storage.PostgreSQLStorage); ok {
		// Get the underlying GORM DB from PostgreSQL storage
		db := postgresStorage.GetDB()
		if db != nil {
			dbCleaner = NewDatabaseCleaner(db)
		}
	}

	return &TestStack{
		Container:       container,
		Storage:         stor,
		ClientRepo:      clientRepo,
		BillingService:  billingService,
		HTTPServer:      httpServer,
		DatabaseCleaner: dbCleaner,
	}
}

// PostgresTestStack provides a test stack for PostgreSQL integration tests
// Note: This is an alias for TestStack - kept for backward compatibility
type PostgresTestStack = TestStack

// NewPostgresTestServer creates an HTTP server with PostgreSQL storage for integration tests
func NewPostgresTestServer() *httpserver.Server {
	return NewIntegrationTestServer()
}

// NewPostgresTestStack creates a complete PostgreSQL test stack for integration tests
func NewPostgresTestStack() *PostgresTestStack {
	return NewIntegrationTestStack()
}

// Configuration helpers for different test scenarios

// WithSharedDependencies returns a server using shared singleton dependencies for unit tests
// Use for unit tests where dependency state can be shared
func WithSharedDependencies() *httpserver.Server {
	return NewUnitTestServer()
}

// WithIsolatedDependencies returns a server with fresh, isolated dependencies for unit tests
// Use for unit tests that need clean state or modify dependencies
func WithIsolatedDependencies() *httpserver.Server {
	return NewIsolatedUnitTestServer()
}

// WithIntegrationDependencies returns a server with PostgreSQL for integration tests
// Use for integration tests that need to test real database behavior
func WithIntegrationDependencies() *httpserver.Server {
	return NewIntegrationTestServer()
}

// NewCleanIntegrationTestServer creates an HTTP server with clean database state
// This function automatically cleans up test data before returning the server
func NewCleanIntegrationTestServer() *httpserver.Server {
	stack := NewIntegrationTestStack()

	// Clean up any existing test data
	if stack.DatabaseCleaner != nil {
		if err := stack.DatabaseCleaner.CleanupTestData(); err != nil {
			panic("Failed to cleanup test data: " + err.Error())
		}
	}

	return stack.HTTPServer
}

// NewCleanIntegrationTestStack creates a complete integration test stack with clean database
// This function automatically cleans up test data before returning the stack
func NewCleanIntegrationTestStack() *TestStack {
	stack := NewIntegrationTestStack()

	// Clean up any existing test data
	if stack.DatabaseCleaner != nil {
		if err := stack.DatabaseCleaner.CleanupTestData(); err != nil {
			panic("Failed to cleanup test data: " + err.Error())
		}
	}

	return stack
}

// NewIntegrationTestServerNoCleanup creates an integration test server without automatic cleanup
// Use this for debugging when you need to inspect test data between runs
func NewIntegrationTestServerNoCleanup() *httpserver.Server {
	container := NewIntegrationTestContainer()
	server, err := container.GetHTTPServer()
	if err != nil {
		panic("Failed to create integration test server: " + err.Error())
	}
	return server
}

// NewIntegrationTestStackNoCleanup creates an integration test stack without automatic cleanup
// Use this for debugging when you need to inspect test data between runs
func NewIntegrationTestStackNoCleanup() *TestStack {
	container := NewIntegrationTestContainer()
	return createTestStack(container)
}

// CleanupIntegrationTestData provides a standalone cleanup function
// Use this in test setup/teardown when you need manual control over cleanup timing
func CleanupIntegrationTestData() error {
	stack := NewIntegrationTestStackNoCleanup()

	if stack.DatabaseCleaner == nil {
		return fmt.Errorf("database cleaner not available - not using PostgreSQL storage")
	}

	return stack.DatabaseCleaner.CleanupTestData()
}
