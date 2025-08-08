// Dependency Injection Builders
//
// This file implements builder patterns for creating DI containers in different environments.
// Provides: Environment-specific container builders, fluent API for configuration
// Pattern: Builder pattern with method chaining for flexible container creation
// Used by: Test setups, production initialization, development environments
package di

import (
	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
)

// ContainerBuilder provides a fluent API for building DI containers
type ContainerBuilder struct {
	config *ContainerConfig
}

// NewContainerBuilder creates a new container builder
func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		config: &ContainerConfig{},
	}
}

// WithConfig sets the entire configuration
func (b *ContainerBuilder) WithConfig(config *ContainerConfig) *ContainerBuilder {
	b.config = config
	return b
}

// WithStorageType sets the storage type
func (b *ContainerBuilder) WithStorageType(storageType string) *ContainerBuilder {
	b.config.StorageType = storageType
	return b
}

// WithEnvironment sets the environment
func (b *ContainerBuilder) WithEnvironment(env string) *ContainerBuilder {
	b.config.Environment = env
	return b
}

// WithLogLevel sets the log level
func (b *ContainerBuilder) WithLogLevel(level string) *ContainerBuilder {
	b.config.LogLevel = level
	return b
}

// WithServerConfig sets server configuration
func (b *ContainerBuilder) WithServerConfig(host string, port int) *ContainerBuilder {
	b.config.ServerHost = host
	b.config.ServerPort = port
	return b
}

// WithDatabaseConfig sets database configuration
func (b *ContainerBuilder) WithDatabaseConfig(host string, port int, name, user, password string) *ContainerBuilder {
	b.config.DatabaseHost = host
	b.config.DatabasePort = port
	b.config.DatabaseName = name
	b.config.DatabaseUser = user
	b.config.DatabasePassword = password
	return b
}

// Build creates and returns the configured container
func (b *ContainerBuilder) Build() *Container {
	return NewContainer(b.config)
}

// TestContainerBuilder provides pre-configured builders for testing
type TestContainerBuilder struct {
	*ContainerBuilder
}

// NewTestContainerBuilder creates a builder pre-configured for testing
func NewTestContainerBuilder() *TestContainerBuilder {
	builder := NewContainerBuilder().WithConfig(TestConfig())
	return &TestContainerBuilder{ContainerBuilder: builder}
}

// WithInMemoryStorage configures the builder for in-memory storage
func (b *TestContainerBuilder) WithInMemoryStorage() *TestContainerBuilder {
	b.WithStorageType("memory")
	return b
}

// WithPostgresStorage configures the builder for PostgreSQL storage (for future use)
func (b *TestContainerBuilder) WithPostgresStorage() *TestContainerBuilder {
	b.WithStorageType("postgres")
	return b
}

// BuildContainer builds and returns the test container
func (b *TestContainerBuilder) BuildContainer() *Container {
	return b.Build()
}

// BuildServer builds the container and returns an HTTP server
func (b *TestContainerBuilder) BuildServer() (*httpserver.Server, error) {
	container := b.Build()
	return container.GetHTTPServer()
}

// DevelopmentContainerBuilder provides pre-configured builders for development
type DevelopmentContainerBuilder struct {
	*ContainerBuilder
}

// NewDevelopmentContainerBuilder creates a builder pre-configured for development
func NewDevelopmentContainerBuilder() *DevelopmentContainerBuilder {
	builder := NewContainerBuilder().WithConfig(DevelopmentConfig())
	return &DevelopmentContainerBuilder{ContainerBuilder: builder}
}

// BuildContainer builds and returns the development container
func (b *DevelopmentContainerBuilder) BuildContainer() *Container {
	return b.Build()
}

// ProductionContainerBuilder provides pre-configured builders for production
type ProductionContainerBuilder struct {
	*ContainerBuilder
}

// NewProductionContainerBuilder creates a builder pre-configured for production
func NewProductionContainerBuilder() *ProductionContainerBuilder {
	builder := NewContainerBuilder().WithConfig(ProductionConfig())
	return &ProductionContainerBuilder{ContainerBuilder: builder}
}

// BuildContainer builds and returns the production container
func (b *ProductionContainerBuilder) BuildContainer() *Container {
	return b.Build()
}

// Convenience functions for quick container creation

// NewTestContainer creates a test container with default configuration
func NewTestContainer() *Container {
	return NewTestContainerBuilder().BuildContainer()
}

// NewTestServer creates a test HTTP server with default configuration
func NewTestServer() (*httpserver.Server, error) {
	return NewTestContainerBuilder().BuildServer()
}

// NewDevelopmentContainer creates a development container with default configuration
func NewDevelopmentContainer() *Container {
	return NewDevelopmentContainerBuilder().BuildContainer()
}

// NewProductionContainer creates a production container with default configuration
func NewProductionContainer() *Container {
	return NewProductionContainerBuilder().BuildContainer()
}
