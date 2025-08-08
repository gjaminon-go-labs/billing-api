// Dependency Injection Configuration
//
// This file defines configuration structures for the DI container.
// Provides: Environment-specific DI configuration, component selection
// Used by: Container builders, test setups, production initialization
package di

// ContainerConfig defines configuration for dependency injection
type ContainerConfig struct {
	// Storage configuration
	StorageType string `yaml:"storage_type" json:"storage_type"`

	// Database configuration (for PostgreSQL) - Application user
	DatabaseURL      string `yaml:"database_url" json:"database_url"`
	DatabaseHost     string `yaml:"database_host" json:"database_host"`
	DatabasePort     int    `yaml:"database_port" json:"database_port"`
	DatabaseName     string `yaml:"database_name" json:"database_name"`
	DatabaseUser     string `yaml:"database_user" json:"database_user"`
	DatabasePassword string `yaml:"database_password" json:"database_password"`
	DatabaseSchema   string `yaml:"database_schema" json:"database_schema"`

	// Migration database configuration - Migration user for DDL operations
	MigrationDatabaseURL      string `yaml:"migration_database_url" json:"migration_database_url"`
	MigrationDatabaseHost     string `yaml:"migration_database_host" json:"migration_database_host"`
	MigrationDatabasePort     int    `yaml:"migration_database_port" json:"migration_database_port"`
	MigrationDatabaseName     string `yaml:"migration_database_name" json:"migration_database_name"`
	MigrationDatabaseUser     string `yaml:"migration_database_user" json:"migration_database_user"`
	MigrationDatabasePassword string `yaml:"migration_database_password" json:"migration_database_password"`
	MigrationDatabaseSchema   string `yaml:"migration_database_schema" json:"migration_database_schema"`

	// Migration configuration
	MigrationEnabled     bool   `yaml:"migration_enabled" json:"migration_enabled"`
	MigrationPath        string `yaml:"migration_path" json:"migration_path"`
	MigrationAutoMigrate bool   `yaml:"migration_auto_migrate" json:"migration_auto_migrate"`
	MigrationTableName   string `yaml:"migration_table_name" json:"migration_table_name"`

	// Test configuration
	TestCleanupEnabled bool `yaml:"test_cleanup_enabled" json:"test_cleanup_enabled"`
	TestCleanupOnSetup bool `yaml:"test_cleanup_on_setup" json:"test_cleanup_on_setup"`

	// Logging configuration
	LogLevel string `yaml:"log_level" json:"log_level"`

	// Server configuration
	ServerPort int    `yaml:"server_port" json:"server_port"`
	ServerHost string `yaml:"server_host" json:"server_host"`

	// Environment
	Environment string `yaml:"environment" json:"environment"`

	// Version information
	Version string `yaml:"version" json:"version"`
}

// UnitTestConfig returns a configuration suitable for unit testing (memory storage)
func UnitTestConfig() *ContainerConfig {
	return &ContainerConfig{
		StorageType: "memory",
		LogLevel:    "debug",
		ServerPort:  8080,
		ServerHost:  "localhost",
		Environment: "test",
	}
}

// IntegrationTestConfig returns a configuration suitable for integration testing (PostgreSQL)
func IntegrationTestConfig() *ContainerConfig {
	return &ContainerConfig{
		StorageType: "postgres",
		// Application database configuration (DML operations)
		DatabaseURL:      "postgres://billing_app_tst_user:billing_app_tst_2025@localhost:5432/go-labs-tst?sslmode=disable&search_path=billing",
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "go-labs-tst",
		DatabaseUser:     "billing_app_tst_user",
		DatabasePassword: "billing_app_tst_2025",
		DatabaseSchema:   "billing",
		// Migration database configuration (DDL operations)
		MigrationDatabaseURL:      "postgres://billing_migration_tst_user:billing_migration_tst_2025@localhost:5432/go-labs-tst?sslmode=disable&search_path=billing",
		MigrationDatabaseHost:     "localhost",
		MigrationDatabasePort:     5432,
		MigrationDatabaseName:     "go-labs-tst",
		MigrationDatabaseUser:     "billing_migration_tst_user",
		MigrationDatabasePassword: "billing_migration_tst_2025",
		MigrationDatabaseSchema:   "billing",
		// Test configuration
		TestCleanupEnabled:   true, // Enable test data cleanup by default
		TestCleanupOnSetup:   true, // Cleanup on test setup by default
		MigrationEnabled:     true,
		MigrationPath:        "database/migrations", // Relative to project root
		MigrationAutoMigrate: false,
		MigrationTableName:   "schema_migrations",
		LogLevel:             "debug",
		ServerPort:           8080,
		ServerHost:           "localhost",
		Environment:          "test",
	}
}

// DEPRECATED METHODS - These methods are kept for backward compatibility
// but should not be used in new code. They will be removed in a future version.

// TestConfig returns a configuration suitable for unit testing (backward compatibility)
// Deprecated: Use UnitTestConfig() instead
func TestConfig() *ContainerConfig {
	return UnitTestConfig()
}

// DevelopmentConfig returns a configuration suitable for development
func DevelopmentConfig() *ContainerConfig {
	return &ContainerConfig{
		StorageType: "postgres",
		// Application database configuration (DML operations)
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "go-labs-dev",
		DatabaseUser:     "billing_app_dev_user",
		DatabasePassword: "billing_app_dev_2025",
		DatabaseSchema:   "billing",
		// Migration database configuration (DDL operations)
		MigrationDatabaseHost:     "localhost",
		MigrationDatabasePort:     5432,
		MigrationDatabaseName:     "go-labs-dev",
		MigrationDatabaseUser:     "billing_migration_dev_user",
		MigrationDatabasePassword: "billing_migration_dev_2025",
		MigrationDatabaseSchema:   "billing",
		LogLevel:                  "debug",
		ServerPort:                8080,
		ServerHost:                "0.0.0.0",
		Environment:               "development",
	}
}

// ProductionConfig returns a configuration suitable for production
func ProductionConfig() *ContainerConfig {
	return &ContainerConfig{
		StorageType: "postgres",
		// Production values should come from environment variables
		LogLevel:    "warn",
		ServerPort:  8080,
		ServerHost:  "0.0.0.0",
		Environment: "production",
	}
}
