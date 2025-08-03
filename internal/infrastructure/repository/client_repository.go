package repository

import (
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
)

// ClientRepositoryImpl implements the ClientRepository interface using a storage backend
type ClientRepositoryImpl struct {
	storage storage.Storage
}

// NewClientRepository creates a new client repository with the given storage backend
func NewClientRepository(storage storage.Storage) repository.ClientRepository {
	return &ClientRepositoryImpl{
		storage: storage,
	}
}

// Save persists a client entity using the storage backend
func (r *ClientRepositoryImpl) Save(client *entity.Client) error {
	// Single Save logic - works with any storage backend
	err := r.storage.Store(client.ID(), client)
	if err != nil {
		// Wrap storage error with repository context
		return errors.NewRepositoryError(
			"save_client",
			errors.RepositoryInternal,
			"failed to save client",
			err,
		)
	}
	return nil
}