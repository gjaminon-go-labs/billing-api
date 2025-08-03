package repository

import (
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
)

// UserRepository defines the interface for user data persistence operations
// Following the Repository pattern from Domain-Driven Design
type UserRepository interface {
	// Create stores a new user in the repository
	Create(user *entity.User) error
	
	// GetByID retrieves a user by their unique identifier
	// Returns nil if user not found (not an error in DDD terms)
	GetByID(id string) (*entity.User, error)
	
	// GetByEmail retrieves a user by their email address
	// Used for authentication and preventing duplicate emails
	GetByEmail(email string) (*entity.User, error)
	
	// Update modifies an existing user in the repository
	Update(user *entity.User) error
	
	// Delete removes a user from the repository
	Delete(id string) error
	
	// List retrieves multiple users with pagination support
	List(limit, offset int) ([]*entity.User, error)
	
	// Count returns the total number of users in the repository
	Count() (int64, error)
}