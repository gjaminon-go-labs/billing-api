// Configuration to DI Integration
//
// This file integrates the configuration system with the DI container.
// Provides: Configuration-driven DI container setup, production configuration mapping
// Pattern: Adapter pattern to convert Config to ContainerConfig
// Used by: Production main.go, DI container builders
package config

import (
	"fmt"

	"github.com/gjaminon-go-labs/billing-api/internal/di"
)

// ToDIConfig converts application config to DI container config
func (c *Config) ToDIConfig() *di.ContainerConfig {
	return &di.ContainerConfig{
		// Storage configuration - read from YAML/environment variables
		StorageType: c.Storage.Type,

		// Database configuration (application user)
		DatabaseURL:      c.buildDatabaseURL(),
		DatabaseHost:     c.Database.Host,
		DatabasePort:     c.Database.Port,
		DatabaseName:     c.Database.DBName,
		DatabaseUser:     c.Database.User,
		DatabasePassword: c.Database.Password,
		DatabaseSchema:   c.Database.Schema,

		// Migration database configuration (migration user)
		MigrationDatabaseURL:      c.buildMigrationDatabaseURL(),
		MigrationDatabaseHost:     c.MigrationDatabase.Host,
		MigrationDatabasePort:     c.MigrationDatabase.Port,
		MigrationDatabaseName:     c.MigrationDatabase.DBName,
		MigrationDatabaseUser:     c.MigrationDatabase.User,
		MigrationDatabasePassword: c.MigrationDatabase.Password,
		MigrationDatabaseSchema:   c.MigrationDatabase.Schema,

		// Migration configuration
		MigrationEnabled:     c.Migration.Enabled,
		MigrationPath:        c.Migration.Path,
		MigrationAutoMigrate: c.Migration.AutoMigrate,
		MigrationTableName:   c.Migration.TableName,

		// Logging configuration
		LogLevel: c.Logging.Level,

		// Server configuration
		ServerPort: c.Server.Port,
		ServerHost: c.Server.Host,

		// Environment detection
		Environment: detectEnvironment(c),
	}
}

// buildDatabaseURL constructs a PostgreSQL connection URL for application user
func (c *Config) buildDatabaseURL() string {
	// postgresql://user:password@host:port/dbname?sslmode=disable&search_path=schema
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode)

	if c.Database.Schema != "" {
		url += "&search_path=" + c.Database.Schema
	}

	return url
}

// buildMigrationDatabaseURL constructs a PostgreSQL connection URL for migration user
func (c *Config) buildMigrationDatabaseURL() string {
	// If migration database is not configured, return empty string (fallback to main database)
	if c.MigrationDatabase.Host == "" || c.MigrationDatabase.User == "" {
		return ""
	}

	// postgresql://user:password@host:port/dbname?sslmode=disable&search_path=schema
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.MigrationDatabase.User,
		c.MigrationDatabase.Password,
		c.MigrationDatabase.Host,
		c.MigrationDatabase.Port,
		c.MigrationDatabase.DBName,
		c.MigrationDatabase.SSLMode)

	if c.MigrationDatabase.Schema != "" {
		url += "&search_path=" + c.MigrationDatabase.Schema
	}

	return url
}

// detectEnvironment determines the environment from configuration
func detectEnvironment(c *Config) string {
	// Try to detect from new shared database name patterns
	if c.Database.DBName == "go-labs-dev" {
		return "development"
	}
	if c.Database.DBName == "go-labs-tst" {
		return "test"
	}
	if c.Database.DBName == "go-labs-qua" {
		return "staging"
	}
	if c.Database.DBName == "go-labs-prd" {
		return "production"
	}

	// Legacy patterns for backward compatibility
	if c.Database.DBName == "billing_service_dev" {
		return "development"
	}
	if c.Database.DBName == "billing_service_test" {
		return "test"
	}
	if c.Database.DBName == "billing_service" {
		return "production"
	}

	// Default to development for safety
	return "development"
}

// NewProductionContainer creates a DI container from application configuration
func NewProductionContainer(config *Config) *di.Container {
	diConfig := config.ToDIConfig()
	return di.NewContainer(diConfig)
}

// NewProductionContainerWithVersion creates a DI container with version information
func NewProductionContainerWithVersion(config *Config, version string) *di.Container {
	diConfig := config.ToDIConfig()
	diConfig.Version = version
	return di.NewContainer(diConfig)
}

// NewProductionContainerFromEnvironment loads config and creates DI container
func NewProductionContainerFromEnvironment(environment string) (*di.Container, error) {
	config, err := LoadConfig(environment)
	if err != nil {
		return nil, err
	}

	container := NewProductionContainer(config)
	return container, nil
}

// NewProductionContainerFromEnvironmentWithVersion loads config and creates DI container with version
func NewProductionContainerFromEnvironmentWithVersion(environment string, version string) (*di.Container, error) {
	config, err := LoadConfig(environment)
	if err != nil {
		return nil, err
	}

	container := NewProductionContainerWithVersion(config, version)
	return container, nil
}
