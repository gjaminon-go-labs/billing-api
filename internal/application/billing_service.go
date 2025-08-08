package application

import (
	"strings"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/google/uuid"
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

// CreateClient creates a new client from a DTO request
func (s *BillingService) CreateClient(req dtos.CreateClientRequest) (*entity.Client, error) {
	client, err := entity.NewClient(req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		return nil, err
	}

	err = s.clientRepo.Save(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateClientLegacy creates a new client with the provided details and persists it (for backward compatibility)
func (s *BillingService) CreateClientLegacy(name, email, phone, address string) (*entity.Client, error) {
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

// PaginatedClients represents clients with pagination metadata
type PaginatedClients struct {
	Clients    []*entity.Client
	Pagination PaginationMeta
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int
	Limit      int
	TotalCount int
	TotalPages int
}

// ListClientsWithPagination retrieves clients with pagination
func (s *BillingService) ListClientsWithPagination(page, limit int) (*PaginatedClients, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	totalCount, err := s.clientRepo.CountClients()
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := totalCount / limit
	if totalCount%limit > 0 {
		totalPages++
	}

	// Get paginated results
	clients, err := s.clientRepo.ListClientsWithPagination(offset, limit)
	if err != nil {
		return nil, err
	}

	return &PaginatedClients{
		Clients: clients,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	}, nil
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

// isValidUUID validates UUID format using the standard library
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
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
// NOTE: For production use, consider using a dedicated phone validation library
// like github.com/nyaruka/phonenumbers for comprehensive international format support
func isValidPhoneFormat(phone string) bool {
	// Basic check: starts with + and contains only digits, spaces, dashes, parentheses
	if !strings.HasPrefix(phone, "+") {
		return false
	}

	// Must have at least a country code (1-3 digits after +)
	if len(phone) < 4 {
		return false
	}

	// Check that it has valid characters and digit count
	digitCount := 0
	for i, char := range phone {
		if char >= '0' && char <= '9' {
			digitCount++
		} else if i == 0 && char == '+' {
			continue // Allow leading +
		} else if char == ' ' || char == '-' || char == '(' || char == ')' || char == '.' {
			continue // Allow common formatting characters
		} else {
			return false // Invalid character
		}
	}

	// International phone number length range (E.164 standard allows 4-15 digits including country code)
	return digitCount >= 4 && digitCount <= 15
}
