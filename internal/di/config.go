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
	
	// Database configuration (for PostgreSQL)
	DatabaseURL      string `yaml:"database_url" json:"database_url"`
	DatabaseHost     string `yaml:"database_host" json:"database_host"`
	DatabasePort     int    `yaml:"database_port" json:"database_port"`
	DatabaseName     string `yaml:"database_name" json:"database_name"`
	DatabaseUser     string `yaml:"database_user" json:"database_user"`
	DatabasePassword string `yaml:"database_password" json:"database_password"`
	DatabaseSchema   string `yaml:"database_schema" json:"database_schema"`
	
	// Migration configuration
	MigrationEnabled     bool   `yaml:"migration_enabled" json:"migration_enabled"`
	MigrationPath        string `yaml:"migration_path" json:"migration_path"`
	MigrationAutoMigrate bool   `yaml:"migration_auto_migrate" json:"migration_auto_migrate"`
	MigrationTableName   string `yaml:"migration_table_name" json:"migration_table_name"`
	
	// Logging configuration
	LogLevel string `yaml:"log_level" json:"log_level"`
	
	// Server configuration
	ServerPort int    `yaml:"server_port" json:"server_port"`
	ServerHost string `yaml:"server_host" json:"server_host"`
	
	// Environment
	Environment string `yaml:"environment" json:"environment"`
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
		StorageType:      "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     5433, // Docker PostgreSQL port to avoid conflicts
		DatabaseName:     "billing_service_test",
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseSchema:   "billing",
		MigrationEnabled:     true,
		MigrationPath:        "../../database/migrations", // Relative to test directory
		MigrationAutoMigrate: true,
		MigrationTableName:   "schema_migrations",
		LogLevel:         "debug",
		ServerPort:       8080,
		ServerHost:       "localhost",
		Environment:      "test",
	}
}

// TestConfig returns a configuration suitable for unit testing (backward compatibility)
// Deprecated: Use UnitTestConfig() instead
func TestConfig() *ContainerConfig {
	return UnitTestConfig()
}

// DevelopmentConfig returns a configuration suitable for development
func DevelopmentConfig() *ContainerConfig {
	return &ContainerConfig{
		StorageType:      "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "billing_service_dev",
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		LogLevel:         "debug",
		ServerPort:       8080,
		ServerHost:       "0.0.0.0",
		Environment:      "development",
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