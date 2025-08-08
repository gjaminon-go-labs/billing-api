// Migration Service
//
// This file implements database migration management using golang-migrate.
// Provides: Schema migration, version tracking, rollback support
// Pattern: Service layer for database schema management
// Used by: Application startup, migration CLI, health checks
package migration

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Service handles database migrations
type Service struct {
	migrator       *migrate.Migrate
	migrationsPath string
	databaseURL    string
}

// Config holds migration service configuration
type Config struct {
	DatabaseURL    string
	MigrationsPath string
	SchemaName     string
}

// NewService creates a new migration service
func NewService(config *Config) (*Service, error) {
	// Validate configuration
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}
	if config.MigrationsPath == "" {
		return nil, fmt.Errorf("migrations path is required")
	}

	// Convert to absolute path for file source
	absPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	service := &Service{
		migrationsPath: absPath,
		databaseURL:    config.DatabaseURL,
	}

	// Initialize the migrator
	if err := service.initMigrator(config.SchemaName); err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %w", err)
	}

	return service, nil
}

// initMigrator initializes the golang-migrate instance
func (s *Service) initMigrator(schemaName string) error {
	// Open database connection
	db, err := sql.Open("postgres", s.databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create postgres driver instance with schema support
	driverConfig := &postgres.Config{
		MigrationsTable: "schema_migrations",
		SchemaName:      schemaName, // Use billing schema
	}

	driver, err := postgres.WithInstance(db, driverConfig)
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance with file source
	sourceURL := fmt.Sprintf("file://%s", s.migrationsPath)
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	s.migrator = migrator
	return nil
}

// Up runs all pending migrations
func (s *Service) Up() error {
	log.Println("🚀 Running database migrations...")

	if err := s.migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ Database schema is up to date")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// Down rolls back one migration
func (s *Service) Down() error {
	log.Println("🔄 Rolling back one migration...")

	if err := s.migrator.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("ℹ️ No migrations to roll back")
			return nil
		}
		return fmt.Errorf("failed to roll back migration: %w", err)
	}

	log.Println("✅ Migration rolled back successfully")
	return nil
}

// Steps runs a specific number of migrations (positive = up, negative = down)
func (s *Service) Steps(n int) error {
	direction := "up"
	if n < 0 {
		direction = "down"
	}

	log.Printf("🔄 Running %d migrations %s...", abs(n), direction)

	if err := s.migrator.Steps(n); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("ℹ️ No migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run %d migrations: %w", n, err)
	}

	log.Printf("✅ %d migrations completed successfully", abs(n))
	return nil
}

// Version returns the current migration version
func (s *Service) Version() (uint, bool, error) {
	version, dirty, err := s.migrator.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil // No migrations have been run
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Status returns the current migration status
func (s *Service) Status() (*Status, error) {
	version, dirty, err := s.Version()
	if err != nil {
		return nil, err
	}

	status := &Status{
		Version:   version,
		Dirty:     dirty,
		HasSchema: version > 0,
	}

	if dirty {
		status.Message = "Database is in dirty state - manual intervention required"
	} else if version == 0 {
		status.Message = "No migrations have been applied"
	} else {
		status.Message = "Database schema is up to date"
	}

	return status, nil
}

// Force sets the migration version without running migrations (use with caution)
func (s *Service) Force(version int) error {
	log.Printf("⚠️ Forcing migration version to %d (this skips migration execution)", version)

	if err := s.migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("✅ Migration version forced to %d", version)
	return nil
}

// Close closes the migration service and releases resources
func (s *Service) Close() error {
	if s.migrator != nil {
		sourceErr, databaseErr := s.migrator.Close()
		if sourceErr != nil {
			return fmt.Errorf("failed to close migration source: %w", sourceErr)
		}
		if databaseErr != nil {
			return fmt.Errorf("failed to close migration database: %w", databaseErr)
		}
	}
	return nil
}

// Status represents the current migration status
type Status struct {
	Version   uint   `json:"version"`
	Dirty     bool   `json:"dirty"`
	HasSchema bool   `json:"has_schema"`
	Message   string `json:"message"`
}

// abs returns the absolute value of an integer
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Validate validates the migration files and database connection
func (s *Service) Validate() error {
	// Check if migrations directory exists and has files
	// This could be implemented to scan the migrations directory
	// and validate migration file format

	// For now, just verify we can get the version (validates DB connection)
	_, _, err := s.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	return nil
}
