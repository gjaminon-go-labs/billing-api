package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)


func TestBillingService_GetClientByID_Success(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	validScenario := scenarios[0] // First scenario is "Valid Get Client Scenario"
	
	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	
	// Create and save a client first
	now := time.Now().UTC()
	client, err := entity.NewClientWithID(
		validScenario.Client.ID,
		validScenario.Client.Name,
		validScenario.Client.Email,
		validScenario.Client.Phone,
		validScenario.Client.Address,
		now,
		now,
	)
	require.NoError(t, err, "Failed to create test client")
	
	err = clientRepo.Save(client)
	require.NoError(t, err, "Failed to save test client")
	
	// Test GetClientByID
	retrievedClient, err := billingService.GetClientByID(validScenario.Client.ID)
	
	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "GetClientByID should succeed for valid ID")
	assert.NotNil(t, retrievedClient, "Retrieved client should not be nil")
	assert.Equal(t, validScenario.Client.ID, retrievedClient.ID(), "Client ID should match")
	assert.Equal(t, validScenario.Client.Name, retrievedClient.Name(), "Client name should match")
	assert.Equal(t, validScenario.Client.Email, retrievedClient.EmailString(), "Client email should match")
}

func TestBillingService_GetClientByID_NotFound(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	nonExistentID := scenarios[3].NonExistentIDs[0] // First non-existent ID
	
	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	
	// Test GetClientByID with non-existent ID
	retrievedClient, err := billingService.GetClientByID(nonExistentID)
	
	// Assertions - this should FAIL until implemented
	assert.Error(t, err, "GetClientByID should fail for non-existent ID")
	assert.Nil(t, retrievedClient, "Retrieved client should be nil for non-existent ID")
}

func TestBillingService_GetClientByID_InvalidUUID(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	invalidIDs := scenarios[2].InvalidIDs // Invalid UUID scenarios
	
	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	
	for _, invalidID := range invalidIDs {
		t.Run("InvalidID_"+invalidID, func(t *testing.T) {
			// Test GetClientByID with invalid UUID
			retrievedClient, err := billingService.GetClientByID(invalidID)
			
			// Assertions - this should FAIL until implemented
			assert.Error(t, err, "GetClientByID should fail for invalid UUID: %s", invalidID)
			assert.Nil(t, retrievedClient, "Retrieved client should be nil for invalid UUID: %s", invalidID)
		})
	}
}