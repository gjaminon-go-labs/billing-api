package di

import (
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
	"gorm.io/gorm"
)

// NewContainerWithDB creates a new DI container using a provided database connection or transaction.
// This is primarily used for testing with transaction isolation.
func NewContainerWithDB(config *ContainerConfig, db *gorm.DB) *Container {
	// Create a normal container
	container := NewContainer(config)
	
	// Override the storage with a custom PostgreSQL storage using the provided DB
	customStorage := storage.NewPostgreSQLStorageWithDB(db)
	container.storage = customStorage
	
	// Mark storage as initialized so it won't be created again
	container.storageOnce.Do(func() {})
	
	return container
}

// ContainerBuilder extensions for transaction support
func (b *ContainerBuilder) BuildWithDB(db *gorm.DB) *Container {
	return NewContainerWithDB(b.config, db)
}