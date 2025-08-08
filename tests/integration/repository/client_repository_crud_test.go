package repository_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

// BUSINESS_TITLE: Database Client Retrieval by ID
// BUSINESS_DESCRIPTION: Data access layer properly stores and retrieves individual client information from the database, ensuring data persistence and integrity
// USER_STORY: As a system, I want to reliably retrieve client data by ID so that business operations can access accurate client information
// BUSINESS_VALUE: Validates data persistence layer, ensures no data loss, confirms database operations work correctly for individual client access
// SCENARIOS_TESTED: Successful client retrieval by ID, handling of non-existent client IDs, proper error reporting
func TestClientRepository_GetByID_Success(t *testing.T) {
	// Setup integration test stack
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	repo := stack.ClientRepo

	// Load test scenarios
	scenarios := loadRepositoryTestScenarios(t)
	testClient := scenarios[0] // First client scenario

	// Create and save a client
	now := time.Now().UTC()
	client, err := entity.NewClientWithID(
		testClient.ID,
		testClient.Name,
		testClient.Email,
		testClient.Phone,
		testClient.Address,
		now,
		now,
	)
	require.NoError(t, err)

	err = repo.Save(client)
	require.NoError(t, err)

	// Test GetByID
	retrievedClient, err := repo.GetByID(testClient.ID)

	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "GetByID should succeed for existing client")
	assert.NotNil(t, retrievedClient, "Retrieved client should not be nil")
	assert.Equal(t, testClient.ID, retrievedClient.ID(), "Client ID should match")
	assert.Equal(t, testClient.Name, retrievedClient.Name(), "Client name should match")
	assert.Equal(t, testClient.Email, retrievedClient.EmailString(), "Client email should match")
}

// BUSINESS_TITLE: Database Non-existent Client Handling
// BUSINESS_DESCRIPTION: Data access layer properly handles requests for clients that don't exist in the database, returning appropriate error responses
// USER_STORY: As a system, I want proper error handling when requesting non-existent clients so that business logic can handle missing data appropriately
// BUSINESS_VALUE: Ensures robust error handling, prevents system crashes, enables proper business logic flow for missing data scenarios
// SCENARIOS_TESTED: Non-existent client ID returns proper error, error message indicates client not found
func TestClientRepository_GetByID_NotFound(t *testing.T) {
	// Setup integration test stack
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	repo := stack.ClientRepo

	nonExistentID := "999e4567-e89b-12d3-a456-426614174999"

	// Test GetByID with non-existent ID
	retrievedClient, err := repo.GetByID(nonExistentID)

	// Assertions - this should FAIL until implemented
	assert.Error(t, err, "GetByID should fail for non-existent client")
	assert.Nil(t, retrievedClient, "Retrieved client should be nil for non-existent client")
}

// BUSINESS_TITLE: Database Client Deletion
// BUSINESS_DESCRIPTION: Data access layer properly removes client records from the database, ensuring data cleanup and referential integrity
// USER_STORY: As a system, I want to reliably delete client data so that business operations can remove inactive or invalid client records
// BUSINESS_VALUE: Enables data cleanup operations, supports compliance requirements, maintains database integrity
// SCENARIOS_TESTED: Successful client deletion, verification of removal, handling of non-existent client deletion
func TestClientRepository_Delete_Success(t *testing.T) {
	// Setup integration test stack
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	repo := stack.ClientRepo

	// Load test scenarios
	scenarios := loadRepositoryTestScenarios(t)
	testClient := scenarios[0] // First client scenario

	// Create and save a client
	now := time.Now().UTC()
	client, err := entity.NewClientWithID(
		testClient.ID,
		testClient.Name,
		testClient.Email,
		testClient.Phone,
		testClient.Address,
		now,
		now,
	)
	require.NoError(t, err)

	err = repo.Save(client)
	require.NoError(t, err)

	// Verify client exists before deletion
	existingClient, err := repo.GetByID(testClient.ID)
	require.NoError(t, err, "Client should exist before deletion")
	require.NotNil(t, existingClient, "Client should exist before deletion")

	// Test Delete
	err = repo.Delete(testClient.ID)

	// Assertions - this should FAIL until implemented
	assert.NoError(t, err, "Delete should succeed for existing client")

	// Verify client no longer exists
	deletedClient, err := repo.GetByID(testClient.ID)
	assert.Error(t, err, "GetByID should fail after deletion")
	assert.Nil(t, deletedClient, "Client should be nil after deletion")
}

// BUSINESS_TITLE: Database Non-existent Client Deletion Handling
// BUSINESS_DESCRIPTION: Data access layer properly handles deletion requests for clients that don't exist, providing appropriate error responses
// USER_STORY: As a system, I want proper error handling when attempting to delete non-existent clients so that business logic can handle missing data scenarios
// BUSINESS_VALUE: Ensures robust error handling, prevents unexpected system behavior, enables proper business logic flow
// SCENARIOS_TESTED: Non-existent client deletion returns proper error, error indicates client not found
func TestClientRepository_Delete_NotFound(t *testing.T) {
	// Setup integration test stack
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	repo := stack.ClientRepo

	nonExistentID := "999e4567-e89b-12d3-a456-426614174999"

	// Test Delete with non-existent ID
	err := repo.Delete(nonExistentID)

	// Assertions - this should FAIL until implemented
	assert.Error(t, err, "Delete should fail for non-existent client")
}

// loadRepositoryTestScenarios loads test scenarios from JSON file
func loadRepositoryTestScenarios(t *testing.T) []RepositoryTestClient {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")

	// Build path to testdata
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "repository_test_fixtures.json")

	// Read and parse JSON
	data, err := os.ReadFile(testDataPath)
	require.NoError(t, err, "Failed to read repository test fixtures file")

	var scenarios []RepositoryTestClient
	err = json.Unmarshal(data, &scenarios)
	require.NoError(t, err, "Failed to parse repository test fixtures JSON")

	return scenarios
}

// RepositoryTestClient represents a client in repository test data
type RepositoryTestClient struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
