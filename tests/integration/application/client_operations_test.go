// Application Layer Integration Tests - Client Operations
//
// This file contains integration tests for the BillingService application layer.
// Tests: Business orchestration, service coordination, repository integration
// Scope: Integration tests - Application Service + Repository + InMemory Storage
// Use Cases: UC-B-001 (Create Client) - Application orchestration layer
//
// Test Scenarios:
// - Service orchestration of client creation workflow
// - Domain validation enforcement through application layer
// - Repository integration and data persistence
// - Error propagation from domain through application layer
// - Uses external JSON test data shared with domain tests
//
// Components Tested:
// - BillingService (application layer)
// - ClientRepository (repository pattern)
// - InMemoryStorage (test infrastructure)
package billing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

type ClientTestCase struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	ShouldFail  bool   `json:"should_fail"`
	Description string `json:"description"`
}

func TestBillingService_CreateClient(t *testing.T) {
	// Load shared test data
	testCases := loadClientTestCases(t)

	// Set up dependencies with in-memory storage
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Test each scenario
	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			// Act - attempt to create client via billing service orchestration
			client, err := service.CreateClient(testCase.Name, testCase.Email, testCase.Phone, testCase.Address)

			if testCase.ShouldFail {
				// Should fail with validation error from domain layer
				assert.Error(t, err, "Client creation should fail for: %s", testCase.Description)
				assert.Nil(t, client, "Client should be nil when creation fails")
			} else {
				// Should succeed through proper orchestration
				assert.NoError(t, err, "Client creation should succeed for: %s", testCase.Description)
				assert.NotNil(t, client, "Client should not be nil when creation succeeds")

				// Verify client properties if creation succeeded
				if client != nil {
					assert.Equal(t, testCase.Name, client.Name())
					assert.Equal(t, testCase.Email, client.EmailString())
					assert.Equal(t, testCase.Phone, client.PhoneString())
					assert.Equal(t, testCase.Address, client.Address())
					assert.NotEmpty(t, client.ID())
				}
			}
		})
	}
}

// BUSINESS_TITLE: Client List Business Logic
// BUSINESS_DESCRIPTION: Business service layer properly orchestrates client retrieval, ensuring data consistency and business rule enforcement
// USER_STORY: As a business analyst, I want to ensure the system correctly processes client list requests through all business layers
// BUSINESS_VALUE: Validates that business logic layer works correctly, ensures data integrity, confirms proper service orchestration
// SCENARIOS_TESTED: Service layer coordination, business rule application, data persistence validation
func TestBillingService_ListClients_IntegrationTest(t *testing.T) {
	// Set up dependencies with PostgreSQL storage using transaction isolation
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	service := stack.BillingService

	// Load test client fixtures
	fixtures := loadClientFixtures(t)

	// Create test clients from fixtures
	client1, err := service.CreateClient(fixtures[0].Name, fixtures[0].Email, fixtures[0].Phone, fixtures[0].Address)
	assert.NoError(t, err)

	client2, err := service.CreateClient(fixtures[1].Name, fixtures[1].Email, fixtures[1].Phone, fixtures[1].Address)
	assert.NoError(t, err)

	// Act
	clients, err := service.ListClients()

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

	// Verify that client details are properly deserialized from PostgreSQL
	for _, client := range clients {
		if client.EmailString() == fixtures[0].Email {
			assert.Equal(t, client1.ID(), client.ID())
			assert.Equal(t, fixtures[0].Name, client.Name())
			assert.Equal(t, fixtures[0].Phone, client.PhoneString())
			assert.Equal(t, fixtures[0].Address, client.Address())
		} else if client.EmailString() == fixtures[1].Email {
			assert.Equal(t, client2.ID(), client.ID())
			assert.Equal(t, fixtures[1].Name, client.Name())
			assert.Equal(t, fixtures[1].Phone, client.PhoneString())
			assert.Equal(t, fixtures[1].Address, client.Address())
		}
	}
}

type ClientFixture struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadClientFixtures(t *testing.T) []ClientFixture {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to fixture data at tests root
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

func loadClientTestCases(t *testing.T) []ClientTestCase {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to shared test data at tests root
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "client_test_cases.json")

	// Read test data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read test data file")

	// Parse JSON
	var testCases []ClientTestCase
	err = json.Unmarshal(data, &testCases)
	assert.NoError(t, err, "Failed to parse test data JSON")

	return testCases
}
