// Dependency Injection Container
//
// This file implements a Kubernetes-style DI container with lazy initialization.
// Provides: Singleton services, thread-safe initialization, performance optimization
// Pattern: Service container with sync.Once for lazy loading
// Benefits: Memory efficient, thread-safe, performance optimized
package di

import (
	"sync"

	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
	"github.com/gjaminon-go-labs/billing-api/internal/migration"
)

// Container manages all application dependencies using lazy initialization
type Container struct {
	config *ContainerConfig

	// Singleton instances (created once, reused)
	storage          storage.Storage
	migrationService *migration.Service
	clientRepo       repository.ClientRepository
	billingService   *application.BillingService
	httpServer       *httpserver.Server

	// Synchronization for thread-safe lazy initialization
	storageOnce          sync.Once
	migrationServiceOnce sync.Once
	clientRepoOnce       sync.Once
	billingServiceOnce   sync.Once
	httpServerOnce       sync.Once

	// Error tracking for failed initializations
	errors      map[string]error
	errorsMutex sync.RWMutex
}

// NewContainer creates a new DI container with the given configuration
func NewContainer(config *ContainerConfig) *Container {
	return &Container{
		config: config,
		errors: make(map[string]error),
	}
}

// GetStorage returns the storage instance, creating it if necessary
func (c *Container) GetStorage() (storage.Storage, error) {
	c.storageOnce.Do(func() {
		storage, err := StorageProvider(c.config)
		if err != nil {
			c.setError("storage", err)
			return
		}
		c.storage = storage
	})

	if err := c.getError("storage"); err != nil {
		return nil, err
	}
	return c.storage, nil
}

// GetMigrationService returns the migration service instance, creating it if necessary
func (c *Container) GetMigrationService() (*migration.Service, error) {
	c.migrationServiceOnce.Do(func() {
		service, err := MigrationServiceProvider(c.config)
		if err != nil {
			c.setError("migration_service", err)
			return
		}
		c.migrationService = service
	})

	if err := c.getError("migration_service"); err != nil {
		return nil, err
	}
	return c.migrationService, nil
}

// GetClientRepository returns the client repository instance, creating it if necessary
func (c *Container) GetClientRepository() (repository.ClientRepository, error) {
	c.clientRepoOnce.Do(func() {
		storage, err := c.GetStorage()
		if err != nil {
			c.setError("client_repository", NewProviderError("client_repository", err))
			return
		}
		c.clientRepo = ClientRepositoryProvider(storage)
	})

	if err := c.getError("client_repository"); err != nil {
		return nil, err
	}
	return c.clientRepo, nil
}

// GetBillingService returns the billing service instance, creating it if necessary
func (c *Container) GetBillingService() (*application.BillingService, error) {
	c.billingServiceOnce.Do(func() {
		clientRepo, err := c.GetClientRepository()
		if err != nil {
			c.setError("billing_service", NewProviderError("billing_service", err))
			return
		}
		c.billingService = BillingServiceProvider(clientRepo)
	})

	if err := c.getError("billing_service"); err != nil {
		return nil, err
	}
	return c.billingService, nil
}

// GetHTTPServer returns the HTTP server instance, creating it if necessary
func (c *Container) GetHTTPServer() (*httpserver.Server, error) {
	c.httpServerOnce.Do(func() {
		billingService, err := c.GetBillingService()
		if err != nil {
			c.setError("http_server", NewProviderError("http_server", err))
			return
		}
		c.httpServer = HTTPServerProvider(billingService)
	})

	if err := c.getError("http_server"); err != nil {
		return nil, err
	}
	return c.httpServer, nil
}

// Reset clears all cached instances and errors (useful for testing)
func (c *Container) Reset() {
	c.storage = nil
	c.migrationService = nil
	c.clientRepo = nil
	c.billingService = nil
	c.httpServer = nil

	c.storageOnce = sync.Once{}
	c.migrationServiceOnce = sync.Once{}
	c.clientRepoOnce = sync.Once{}
	c.billingServiceOnce = sync.Once{}
	c.httpServerOnce = sync.Once{}

	c.errorsMutex.Lock()
	c.errors = make(map[string]error)
	c.errorsMutex.Unlock()
}

// GetConfig returns the container configuration
func (c *Container) GetConfig() *ContainerConfig {
	return c.config
}

// setError stores an error for a component (thread-safe)
func (c *Container) setError(component string, err error) {
	c.errorsMutex.Lock()
	defer c.errorsMutex.Unlock()
	c.errors[component] = err
}

// getError retrieves an error for a component (thread-safe)
func (c *Container) getError(component string) error {
	c.errorsMutex.RLock()
	defer c.errorsMutex.RUnlock()
	return c.errors[component]
}

// HasErrors returns true if any component has initialization errors
func (c *Container) HasErrors() bool {
	c.errorsMutex.RLock()
	defer c.errorsMutex.RUnlock()
	return len(c.errors) > 0
}

// GetErrors returns all initialization errors
func (c *Container) GetErrors() map[string]error {
	c.errorsMutex.RLock()
	defer c.errorsMutex.RUnlock()

	// Return a copy to prevent concurrent access issues
	errorsCopy := make(map[string]error)
	for k, v := range c.errors {
		errorsCopy[k] = v
	}
	return errorsCopy
}
