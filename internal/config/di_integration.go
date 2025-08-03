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
		
		// Database configuration
		DatabaseURL:      c.buildDatabaseURL(),
		DatabaseHost:     c.Database.Host,
		DatabasePort:     c.Database.Port,
		DatabaseName:     c.Database.DBName,
		DatabaseUser:     c.Database.User,
		DatabasePassword: c.Database.Password,
		DatabaseSchema:   c.Database.Schema,
		
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

// buildDatabaseURL constructs a PostgreSQL connection URL
func (c *Config) buildDatabaseURL() string {
	// postgresql://user:password@host:port/dbname?sslmode=disable
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password, 
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode)
}

// detectEnvironment determines the environment from configuration
func detectEnvironment(c *Config) string {
	// Try to detect from database name patterns
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

// NewProductionContainerFromEnvironment loads config and creates DI container
func NewProductionContainerFromEnvironment(environment string) (*di.Container, error) {
	config, err := LoadConfig(environment)
	if err != nil {
		return nil, err
	}
	
	container := NewProductionContainer(config)
	return container, nil
}