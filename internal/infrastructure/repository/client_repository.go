package repository

import (
	"encoding/json"
	"fmt"
	
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

// GetAll retrieves all client entities from storage
func (r *ClientRepositoryImpl) GetAll() ([]*entity.Client, error) {
	// Get all values from storage
	values, err := r.storage.ListAll()
	if err != nil {
		return nil, errors.NewRepositoryError(
			"get_all_clients",
			errors.RepositoryInternal,
			"failed to retrieve all clients",
			err,
		)
	}
	
	// Convert storage values to domain entities
	clients := make([]*entity.Client, 0, len(values))
	for _, value := range values {
		// Try direct type assertion first (for in-memory storage)
		if client, ok := value.(*entity.Client); ok {
			clients = append(clients, client)
			continue
		}
		
		// Handle JSON deserialization (for PostgreSQL storage)
		if clientMap, ok := value.(map[string]interface{}); ok {
			client, err := r.deserializeClient(clientMap)
			if err != nil {
				return nil, errors.NewRepositoryError(
					"deserialize_client",
					errors.RepositoryInternal,
					"failed to deserialize client",
					err,
				)
			}
			clients = append(clients, client)
		}
	}
	
	return clients, nil
}

// deserializeClient converts a map[string]interface{} back to a Client entity
func (r *ClientRepositoryImpl) deserializeClient(clientMap map[string]interface{}) (*entity.Client, error) {
	// Convert the map back to JSON and then unmarshal using custom unmarshaling
	jsonBytes, err := json.Marshal(clientMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client map to JSON: %w", err)
	}
	
	// Create a new client instance and unmarshal into it
	var client entity.Client
	if err := json.Unmarshal(jsonBytes, &client); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to client: %w", err)
	}
	
	return &client, nil
}

// GetByID retrieves a client entity by ID
func (r *ClientRepositoryImpl) GetByID(id string) (*entity.Client, error) {
	// Get value from storage
	value, err := r.storage.Get(id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == fmt.Sprintf("key not found: %s", id) {
			return nil, errors.ErrClientNotFound
		}
		
		return nil, errors.NewRepositoryError(
			"get_client",
			errors.RepositoryInternal,
			"failed to retrieve client",
			err,
		)
	}
	
	// Try direct type assertion first (for in-memory storage)
	if client, ok := value.(*entity.Client); ok {
		return client, nil
	}
	
	// Handle JSON deserialization (for PostgreSQL storage)
	if clientMap, ok := value.(map[string]interface{}); ok {
		client, err := r.deserializeClient(clientMap)
		if err != nil {
			return nil, errors.NewRepositoryError(
				"deserialize_client",
				errors.RepositoryInternal,
				"failed to deserialize client",
				err,
			)
		}
		return client, nil
	}
	
	return nil, errors.NewRepositoryError(
		"get_client",
		errors.RepositoryInternal,
		"unexpected value type in storage",
		nil,
	)
}

// Delete removes a client entity by ID
func (r *ClientRepositoryImpl) Delete(id string) error {
	// Use storage Delete method
	err := r.storage.Delete(id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == fmt.Sprintf("key not found: %s", id) {
			return errors.ErrClientNotFound
		}
		
		return errors.NewRepositoryError(
			"delete_client",
			errors.RepositoryInternal,
			"failed to delete client",
			err,
		)
	}
	
	return nil
}