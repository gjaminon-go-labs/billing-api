// Dependency Injection Providers
//
// This file contains provider functions for creating services and components.
// Provides: Factory functions for all services, repositories, and infrastructure
// Pattern: Each provider function creates and configures a specific component
// Used by: DI container for lazy initialization of dependencies
package di

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	infrarepo "github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
	"github.com/gjaminon-go-labs/billing-api/internal/migration"
	testinfra "github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

// StorageProvider creates a storage instance based on configuration
func StorageProvider(config *ContainerConfig) (storage.Storage, error) {
	switch config.StorageType {
	case "memory":
		return testinfra.NewInMemoryStorage(), nil
	case "postgres":
		return createPostgreSQLStorage(config)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.StorageType)
	}
}

// createPostgreSQLStorage creates a PostgreSQL-backed storage instance
func createPostgreSQLStorage(config *ContainerConfig) (storage.Storage, error) {
	log.Printf("🐘 Connecting to PostgreSQL at %s:%d...", config.DatabaseHost, config.DatabasePort)

	// Run migrations first if enabled and auto-migrate is true
	if config.MigrationEnabled && config.MigrationAutoMigrate {
		if err := runMigrations(config); err != nil {
			return nil, NewProviderError("postgresql-storage", fmt.Errorf("failed to run migrations: %w", err))
		}
	}

	// Configure GORM with PostgreSQL driver
	gormConfig := &gorm.Config{
		// Disable default transaction for better performance
		SkipDefaultTransaction: true,

		// Prepare statements for better performance
		PrepareStmt: true,
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), gormConfig)
	if err != nil {
		return nil, NewProviderError("postgresql-storage", fmt.Errorf("failed to connect to database: %w", err))
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, NewProviderError("postgresql-storage", fmt.Errorf("failed to get underlying SQL DB: %w", err))
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)                 // Maximum open connections
	sqlDB.SetMaxIdleConns(5)                  // Maximum idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Connection max lifetime
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Connection max idle time

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, NewProviderError("postgresql-storage", fmt.Errorf("failed to ping database: %w", err))
	}

	log.Printf("✅ PostgreSQL connection established successfully")

	// Create PostgreSQL storage with GORM
	return storage.NewPostgreSQLStorage(db), nil
}

// runMigrations runs database migrations if enabled
func runMigrations(config *ContainerConfig) error {
	// Use migration database URL if available, fallback to main database URL for backward compatibility
	databaseURL := config.MigrationDatabaseURL
	if databaseURL == "" {
		databaseURL = config.DatabaseURL
	}

	schema := config.MigrationDatabaseSchema
	if schema == "" {
		schema = config.DatabaseSchema
	}

	migrationConfig := &migration.Config{
		DatabaseURL:    databaseURL,
		MigrationsPath: config.MigrationPath,
		SchemaName:     schema,
	}

	migrationService, err := migration.NewService(migrationConfig)
	if err != nil {
		return fmt.Errorf("failed to create migration service: %w", err)
	}
	defer migrationService.Close()

	// Run migrations
	if err := migrationService.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// MigrationServiceProvider creates a migration service
func MigrationServiceProvider(config *ContainerConfig) (*migration.Service, error) {
	if !config.MigrationEnabled {
		return nil, fmt.Errorf("migrations are disabled in configuration")
	}

	// Use migration database URL if available, fallback to main database URL for backward compatibility
	databaseURL := config.MigrationDatabaseURL
	if databaseURL == "" {
		databaseURL = config.DatabaseURL
	}

	schema := config.MigrationDatabaseSchema
	if schema == "" {
		schema = config.DatabaseSchema
	}

	migrationConfig := &migration.Config{
		DatabaseURL:    databaseURL,
		MigrationsPath: config.MigrationPath,
		SchemaName:     schema,
	}

	service, err := migration.NewService(migrationConfig)
	if err != nil {
		return nil, NewProviderError("migration-service", err)
	}

	return service, nil
}

// ClientRepositoryProvider creates a client repository with the given storage
func ClientRepositoryProvider(storage storage.Storage) repository.ClientRepository {
	return infrarepo.NewClientRepository(storage)
}

// BillingServiceProvider creates a billing service with the given repository
func BillingServiceProvider(clientRepo repository.ClientRepository) *application.BillingService {
	return application.NewBillingService(clientRepo)
}

// HTTPServerProvider creates an HTTP server with the given services
func HTTPServerProvider(billingService *application.BillingService) *httpserver.Server {
	return httpserver.NewServer(billingService)
}

// ProviderError represents an error in provider creation
type ProviderError struct {
	Component string
	Err       error
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("failed to create %s: %v", e.Component, e.Err)
}

// NewProviderError creates a new provider error
func NewProviderError(component string, err error) *ProviderError {
	return &ProviderError{
		Component: component,
		Err:       err,
	}
}
