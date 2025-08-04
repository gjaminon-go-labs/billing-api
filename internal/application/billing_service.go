package application

import (
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
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