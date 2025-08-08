package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

// UpdateClientScenario represents test data for update client operations
type UpdateClientScenario struct {
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Request         dtos.UpdateClientRequest `json:"request"`
	ExpectedClient  ClientTestData           `json:"expected_client"`
	InvalidRequests []InvalidUpdateRequest   `json:"invalid_requests"`
}

// InvalidUpdateRequest represents an invalid update request scenario
type InvalidUpdateRequest struct {
	Description   string                   `json:"description"`
	Request       dtos.UpdateClientRequest `json:"request"`
	ExpectedError string                   `json:"expected_error"`
}

func loadUpdateClientScenarios(t *testing.T) []UpdateClientScenario {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")

	// Build path to testdata
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "update_client_requests.json")

	// Read and parse JSON
	data, err := os.ReadFile(testDataPath)
	require.NoError(t, err, "Failed to read update client scenarios file")

	var scenarios []UpdateClientScenario
	err = json.Unmarshal(data, &scenarios)
	require.NoError(t, err, "Failed to parse update client scenarios JSON")

	return scenarios
}

func TestBillingService_UpdateClient_Success(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	fullUpdateScenario := scenarios[0] // "Full Update Request"

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	// Create and save the original client
	now := time.Now().UTC()
	originalClient, err := entity.NewClientWithID(
		fullUpdateScenario.ExpectedClient.ID,
		"Alice Johnson", // Original name
		fullUpdateScenario.ExpectedClient.Email,
		"+1234567890",                        // Original phone
		"123 Main Street, Anytown, ST 12345", // Original address
		now,
		now,
	)
	require.NoError(t, err, "Failed to create test client")

	err = clientRepo.Save(originalClient)
	require.NoError(t, err, "Failed to save test client")

	// Test UpdateClient
	updatedClient, err := billingService.UpdateClient(
		fullUpdateScenario.ExpectedClient.ID,
		fullUpdateScenario.Request,
	)

	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "UpdateClient should succeed for valid request")
	assert.NotNil(t, updatedClient, "Updated client should not be nil")
	assert.Equal(t, fullUpdateScenario.ExpectedClient.ID, updatedClient.ID(), "Client ID should remain unchanged")
	assert.Equal(t, fullUpdateScenario.Request.Name, updatedClient.Name(), "Client name should be updated")
	assert.Equal(t, fullUpdateScenario.Request.Phone, updatedClient.PhoneString(), "Client phone should be updated")
	assert.Equal(t, fullUpdateScenario.Request.Address, updatedClient.Address(), "Client address should be updated")
	assert.Equal(t, fullUpdateScenario.ExpectedClient.Email, updatedClient.EmailString(), "Client email should remain unchanged")
}

func TestBillingService_UpdateClient_PartialUpdate(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	partialUpdateScenario := scenarios[1] // "Partial Update Request - Name Only"

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	// Create and save the original client
	now := time.Now().UTC()
	originalClient, err := entity.NewClientWithID(
		partialUpdateScenario.ExpectedClient.ID,
		"Alice Johnson", // Original name to be updated
		partialUpdateScenario.ExpectedClient.Email,
		"+1234567890",                        // Original phone (should be cleared)
		"123 Main Street, Anytown, ST 12345", // Original address (should be cleared)
		now,
		now,
	)
	require.NoError(t, err, "Failed to create test client")

	err = clientRepo.Save(originalClient)
	require.NoError(t, err, "Failed to save test client")

	// Test UpdateClient with partial update
	updatedClient, err := billingService.UpdateClient(
		partialUpdateScenario.ExpectedClient.ID,
		partialUpdateScenario.Request,
	)

	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "UpdateClient should succeed for partial update")
	assert.NotNil(t, updatedClient, "Updated client should not be nil")
	assert.Equal(t, partialUpdateScenario.Request.Name, updatedClient.Name(), "Client name should be updated")
	assert.Equal(t, "", updatedClient.PhoneString(), "Client phone should be cleared")
	assert.Equal(t, "", updatedClient.Address(), "Client address should be cleared")
}

func TestBillingService_UpdateClient_NotFound(t *testing.T) {
	// Load test scenarios
	getScenarios := loadGetClientScenarios(t)
	nonExistentID := getScenarios[3].NonExistentIDs[0] // First non-existent ID

	updateScenarios := loadUpdateClientScenarios(t)
	updateRequest := updateScenarios[0].Request // Any valid update request

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	// Test UpdateClient with non-existent ID
	updatedClient, err := billingService.UpdateClient(nonExistentID, updateRequest)

	// Assertions - this should FAIL until implemented
	assert.Error(t, err, "UpdateClient should fail for non-existent ID")
	assert.Nil(t, updatedClient, "Updated client should be nil for non-existent ID")
}

func TestBillingService_UpdateClient_ValidationError(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	invalidRequests := scenarios[3].InvalidRequests // "Invalid Update Requests"

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	// Create and save a test client
	now := time.Now().UTC()
	validClient, err := entity.NewClientWithID(
		"123e4567-e89b-12d3-a456-426614174000",
		"Test Client",
		"test@example.com",
		"+1234567890",
		"Test Address",
		now,
		now,
	)
	require.NoError(t, err)

	err = clientRepo.Save(validClient)
	require.NoError(t, err)

	for _, invalidRequest := range invalidRequests {
		t.Run(invalidRequest.Description, func(t *testing.T) {
			// Test UpdateClient with invalid request
			updatedClient, err := billingService.UpdateClient(
				validClient.ID(),
				invalidRequest.Request,
			)

			// Assertions - this should FAIL until implemented
			assert.Error(t, err, "UpdateClient should fail for invalid request: %s", invalidRequest.Description)
			assert.Nil(t, updatedClient, "Updated client should be nil for invalid request")
		})
	}
}

func TestBillingService_UpdateClient_InvalidUUID(t *testing.T) {
	// Load test scenarios
	getScenarios := loadGetClientScenarios(t)
	invalidIDs := getScenarios[2].InvalidIDs // Invalid UUID scenarios

	updateScenarios := loadUpdateClientScenarios(t)
	updateRequest := updateScenarios[0].Request // Any valid update request

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	for _, invalidID := range invalidIDs {
		t.Run("InvalidID_"+invalidID, func(t *testing.T) {
			// Test UpdateClient with invalid UUID
			updatedClient, err := billingService.UpdateClient(invalidID, updateRequest)

			// Assertions - this should FAIL until implemented
			assert.Error(t, err, "UpdateClient should fail for invalid UUID: %s", invalidID)
			assert.Nil(t, updatedClient, "Updated client should be nil for invalid UUID: %s", invalidID)
		})
	}
}
