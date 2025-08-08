package testhelpers

import (
	"fmt"
	"testing"

	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/di"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TransactionTest runs a test within a database transaction that's automatically rolled back.
// This provides perfect test isolation allowing parallel test execution without data conflicts.
func TransactionTest(t *testing.T, testFunc func(*testing.T, *gorm.DB)) {
	// Get base database connection from integration test config
	config := di.IntegrationTestConfig()
	db := setupPostgreSQLConnection(config)

	// Start transaction - this is our isolation boundary
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// CRITICAL: Always rollback - never commit test data
	defer func() {
		if err := tx.Rollback().Error; err != nil {
			// Only log if it's not "transaction already closed" error
			if err != gorm.ErrInvalidTransaction {
				t.Logf("Warning: Failed to rollback transaction: %v", err)
			}
		}
	}()

	// Run the test with transaction
	// All database operations in testFunc will use this transaction
	testFunc(t, tx)
}

// IntegrationTestStack provides access to all components for integration tests
type IntegrationTestStack struct {
	Container       *di.Container
	Storage         storage.Storage
	BillingService  *application.BillingService
	ClientRepo      repository.ClientRepository
	HTTPServer      *httpserver.Server
	DB              *gorm.DB
	DatabaseCleaner *DatabaseCleaner
}

// NewTransactionalTestStack creates an IntegrationTestStack that uses a transaction
// instead of a direct database connection. This ensures test isolation.
func NewTransactionalTestStack(t *testing.T, tx *gorm.DB) *IntegrationTestStack {
	// Build container with transaction instead of creating new DB connection
	config := di.IntegrationTestConfig()

	// Create container with transaction using our custom method
	container := di.NewContainerWithDB(config, tx)

	// Extract services from container - using error-returning methods
	billingService, err := container.GetBillingService()
	if err != nil {
		t.Fatalf("Failed to get billing service: %v", err)
	}

	clientRepo, err := container.GetClientRepository()
	if err != nil {
		t.Fatalf("Failed to get client repository: %v", err)
	}

	httpServer, err := container.GetHTTPServer()
	if err != nil {
		t.Fatalf("Failed to get HTTP server: %v", err)
	}

	// Get storage from container (it's already using the transaction)
	stor, err := container.GetStorage()
	if err != nil {
		t.Fatalf("Failed to get storage: %v", err)
	}

	return &IntegrationTestStack{
		Container:      container,
		Storage:        stor,
		BillingService: billingService,
		ClientRepo:     clientRepo,
		HTTPServer:     httpServer,
		DB:             tx, // Use transaction as DB
		// DatabaseCleaner not needed - transaction rollback handles cleanup
		DatabaseCleaner: nil,
	}
}

// WithTransaction provides an easy adapter for existing tests to use transactions.
// Returns a test stack and a cleanup function that must be deferred.
func WithTransaction(t *testing.T) (*IntegrationTestStack, func()) {
	// Get base database connection
	config := di.IntegrationTestConfig()
	db := setupPostgreSQLConnection(config)

	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// Create stack with transaction
	stack := NewTransactionalTestStack(t, tx)

	// Cleanup function - just rollback the transaction
	cleanup := func() {
		if err := tx.Rollback().Error; err != nil {
			if err != gorm.ErrInvalidTransaction {
				t.Logf("Warning: Failed to rollback transaction: %v", err)
			}
		}
	}

	return stack, cleanup
}

// setupPostgreSQLConnection creates a base database connection for transactions
func setupPostgreSQLConnection(config *di.ContainerConfig) *gorm.DB {
	// Build database connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
		config.DatabaseSchema,
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		// Use simple logger for tests
		Logger: logger.Default.LogMode(logger.Silent),
		// Prepare statements for better performance
		PrepareStmt: true,
		// Skip default transaction for better control
		SkipDefaultTransaction: true,
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Configure connection pool for parallel tests
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get SQL DB: " + err.Error())
	}

	// Increase pool size for parallel test execution
	sqlDB.SetMaxOpenConns(50) // Allow more parallel connections
	sqlDB.SetMaxIdleConns(10) // Keep more connections idle

	return db
}
