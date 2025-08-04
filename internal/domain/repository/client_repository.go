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
}