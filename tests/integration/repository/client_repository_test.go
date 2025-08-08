package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

// BUSINESS_TITLE: Database Client Retrieval
// BUSINESS_DESCRIPTION: Data access layer properly stores and retrieves client information from the database, ensuring data persistence and integrity
// USER_STORY: As a system architect, I want to ensure client data is properly persisted and retrieved from the database
// BUSINESS_VALUE: Validates data persistence layer, ensures no data loss, confirms database operations work correctly
// SCENARIOS_TESTED: Database storage, data retrieval, persistence validation, multi-client handling
func TestClientRepository_GetAll_IntegrationTest(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	repo := stack.ClientRepo

	// Load test fixtures
	fixtures := loadIntegrationRepositoryFixtures(t)

	// Create and save test clients from fixtures
	client1, err := entity.NewClient(fixtures[0].Name, fixtures[0].Email, fixtures[0].Phone, fixtures[0].Address)
	assert.NoError(t, err)
	err = repo.Save(client1)
	assert.NoError(t, err)

	client2, err := entity.NewClient(fixtures[1].Name, fixtures[1].Email, fixtures[1].Phone, fixtures[1].Address)
	assert.NoError(t, err)
	err = repo.Save(client2)
	assert.NoError(t, err)

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 2)

	// Verify clients are present (order may vary)
	expectedEmails := []string{fixtures[0].Email, fixtures[1].Email}
	actualEmails := make([]string, len(clients))
	for i, client := range clients {
		actualEmails[i] = client.EmailString()
	}

	for _, expectedEmail := range expectedEmails {
		assert.Contains(t, actualEmails, expectedEmail)
	}
}

// BUSINESS_TITLE: Empty Database State Handling
// BUSINESS_DESCRIPTION: System properly handles empty database scenarios, ensuring reliable behavior when no client data exists
// USER_STORY: As a system administrator, I want the system to handle empty database states gracefully without errors
// BUSINESS_VALUE: Ensures system stability in edge cases, prevents crashes on empty databases, supports clean deployments
// SCENARIOS_TESTED: Empty database queries, null result handling, system stability with no data
func TestClientRepository_GetAll_EmptyRepository_IntegrationTest(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	repo := stack.ClientRepo

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Empty(t, clients)
}

type ClientFixture struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadIntegrationRepositoryFixtures(t *testing.T) []ClientFixture {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to fixture data
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "client_fixtures.json")

	// Read fixture data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read fixture data file")

	// Parse JSON
	var fixtures []ClientFixture
	err = json.Unmarshal(data, &fixtures)
	assert.NoError(t, err, "Failed to parse fixture data JSON")

	return fixtures
}
