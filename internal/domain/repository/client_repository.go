package repository

import (
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
)

// ClientRepository defines the contract for client persistence operations
type ClientRepository interface {
	// Save persists a client entity
	Save(client *entity.Client) error

	// GetAll retrieves all client entities
	GetAll() ([]*entity.Client, error)

	// GetByID retrieves a client entity by ID
	GetByID(id string) (*entity.Client, error)

	// Delete removes a client entity by ID
	Delete(id string) error
}
