package application

import (
	"strings"
	
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
)

// BillingService orchestrates billing domain operations and use cases
type BillingService struct {
	clientRepo repository.ClientRepository
}

// NewBillingService creates a new billing service
func NewBillingService(clientRepo repository.ClientRepository) *BillingService {
	return &BillingService{
		clientRepo: clientRepo,
	}
}

// CreateClient creates a new client with the provided details and persists it
func (s *BillingService) CreateClient(name, email, phone, address string) (*entity.Client, error) {
	client, err := entity.NewClient(name, email, phone, address)
	if err != nil {
		return nil, err
	}
	
	err = s.clientRepo.Save(client)
	if err != nil {
		return nil, err
	}
	
	return client, nil
}

// ListClients retrieves all clients from the repository
func (s *BillingService) ListClients() ([]*entity.Client, error) {
	return s.clientRepo.GetAll()
}

// GetClientByID retrieves a client by ID
func (s *BillingService) GetClientByID(id string) (*entity.Client, error) {
	// Basic UUID validation
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewValidationError("id", id, errors.ValidationRequired, "client ID is required")
	}
	
	// Simple UUID format validation (basic check)
	if !isValidUUID(id) {
		return nil, errors.NewValidationError("id", id, errors.ValidationFormat, "client ID must be a valid UUID")
	}
	
	// Delegate to repository
	return s.clientRepo.GetByID(id)
}

// isValidUUID performs basic UUID format validation
func isValidUUID(id string) bool {
	// Basic UUID format check (36 characters with dashes at positions 8, 13, 18, 23)
	if len(id) != 36 {
		return false
	}
	
	// Check dash positions
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	
	// Check that other characters are hex digits
	hexChars := "0123456789abcdefABCDEF"
	for i, char := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // Skip dash positions
		}
		
		found := false
		for _, hexChar := range hexChars {
			if char == hexChar {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

// DeleteClient removes a client by ID
func (s *BillingService) DeleteClient(id string) error {
	// Basic UUID validation (reuse validation logic)
	if strings.TrimSpace(id) == "" {
		return errors.NewValidationError("id", id, errors.ValidationRequired, "client ID is required")
	}
	
	if !isValidUUID(id) {
		return errors.NewValidationError("id", id, errors.ValidationFormat, "client ID must be a valid UUID")
	}
	
	// Delegate to repository
	return s.clientRepo.Delete(id)
}

// UpdateClient updates a client by ID
func (s *BillingService) UpdateClient(id string, req dtos.UpdateClientRequest) (*entity.Client, error) {
	// Basic UUID validation (reuse validation logic)
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewValidationError("id", id, errors.ValidationRequired, "client ID is required")
	}
	
	if !isValidUUID(id) {
		return nil, errors.NewValidationError("id", id, errors.ValidationFormat, "client ID must be a valid UUID")
	}
	
	// Validate request data
	if err := validateUpdateRequest(req); err != nil {
		return nil, err
	}
	
	// Get existing client
	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		return nil, err // Repository error (including not found)
	}
	
	// Update client details using domain method
	err = client.UpdateDetails(req.Name, req.Phone, req.Address)
	if err != nil {
		return nil, err // Domain validation error
	}
	
	// Save updated client
	err = s.clientRepo.Save(client)
	if err != nil {
		return nil, err // Repository error
	}
	
	return client, nil
}

// validateUpdateRequest validates the update request data
func validateUpdateRequest(req dtos.UpdateClientRequest) error {
	// Validate name (required)
	if strings.TrimSpace(req.Name) == "" {
		return errors.NewValidationError("name", req.Name, errors.ValidationRequired, "name is required")
	}
	
	if len(strings.TrimSpace(req.Name)) < 2 {
		return errors.NewValidationError("name", req.Name, errors.ValidationLength, "name must be at least 2 characters")
	}
	
	if len(strings.TrimSpace(req.Name)) > 100 {
		return errors.NewValidationError("name", req.Name, errors.ValidationLength, "name must not exceed 100 characters")
	}
	
	// Validate phone (optional, but if provided must be valid)
	if req.Phone != "" && len(req.Phone) > 20 {
		return errors.NewValidationError("phone", req.Phone, errors.ValidationLength, "phone number must not exceed 20 characters")
	}
	
	// Basic phone format validation if provided
	if req.Phone != "" && !isValidPhoneFormat(req.Phone) {
		return errors.NewValidationError("phone", req.Phone, errors.ValidationFormat, "phone number format is invalid")
	}
	
	// Validate address (optional)
	if len(req.Address) > 500 {
		return errors.NewValidationError("address", req.Address, errors.ValidationLength, "address must not exceed 500 characters")
	}
	
	return nil
}

// isValidPhoneFormat performs basic phone format validation
func isValidPhoneFormat(phone string) bool {
	// Basic check: starts with + and contains only digits, spaces, dashes, parentheses
	if !strings.HasPrefix(phone, "+") {
		return false
	}
	
	// Check that it has at least 7 digits (minimum phone number length)
	digitCount := 0
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			digitCount++
		} else if char != '+' && char != ' ' && char != '-' && char != '(' && char != ')' {
			return false // Invalid character
		}
	}
	
	return digitCount >= 7 && digitCount <= 15 // International phone number range
}