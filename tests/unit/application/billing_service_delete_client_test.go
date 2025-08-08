package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

func TestBillingService_DeleteClient_Success(t *testing.T) {
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

	// Verify client exists before deletion
	retrievedClient, err := clientRepo.GetByID(validScenario.Client.ID)
	require.NoError(t, err, "Client should exist before deletion")
	require.NotNil(t, retrievedClient, "Client should exist before deletion")

	// Test DeleteClient
	err = billingService.DeleteClient(validScenario.Client.ID)

	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "DeleteClient should succeed for valid ID")

	// Verify client no longer exists
	retrievedClient, err = clientRepo.GetByID(validScenario.Client.ID)
	assert.Error(t, err, "Client should not exist after deletion")
	assert.Nil(t, retrievedClient, "Client should be nil after deletion")
}

func TestBillingService_DeleteClient_NotFound(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	nonExistentID := scenarios[3].NonExistentIDs[0] // First non-existent ID

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	// Test DeleteClient with non-existent ID
	err := billingService.DeleteClient(nonExistentID)

	// Assertions - this should FAIL until implemented
	assert.Error(t, err, "DeleteClient should fail for non-existent ID")
}

func TestBillingService_DeleteClient_InvalidUUID(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	invalidIDs := scenarios[2].InvalidIDs // Invalid UUID scenarios

	// Setup in-memory storage and repository
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)

	for _, invalidID := range invalidIDs {
		t.Run("InvalidID_"+invalidID, func(t *testing.T) {
			// Test DeleteClient with invalid UUID
			err := billingService.DeleteClient(invalidID)

			// Assertions - this should FAIL until implemented
			assert.Error(t, err, "DeleteClient should fail for invalid UUID: %s", invalidID)
		})
	}
}

// GetClientScenario represents test data for get client operations
type GetClientScenario struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Client         ClientTestData `json:"client"`
	InvalidIDs     []string       `json:"invalid_ids"`
	NonExistentIDs []string       `json:"non_existent_ids"`
}

// ClientTestData represents a client in test data
type ClientTestData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadGetClientScenarios(t *testing.T) []GetClientScenario {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")

	// Build path to testdata
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "get_client_scenarios.json")

	// Read and parse JSON
	data, err := os.ReadFile(testDataPath)
	require.NoError(t, err, "Failed to read get client scenarios file")

	var scenarios []GetClientScenario
	err = json.Unmarshal(data, &scenarios)
	require.NoError(t, err, "Failed to parse get client scenarios JSON")

	return scenarios
}
