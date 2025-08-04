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

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
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


func TestBillingService_ListClients_IntegrationTest(t *testing.T) {
	// Set up dependencies with PostgreSQL storage
	stack := testhelpers.NewCleanIntegrationTestStack()
	service := stack.BillingService
	
	// Create test clients
	client1, err := service.CreateClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)
	
	client2, err := service.CreateClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
	assert.NoError(t, err)

	// Act
	clients, err := service.ListClients()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 2)
	
	// Verify clients are present (order may vary)
	expectedEmails := []string{"john@example.com", "jane@example.com"}
	actualEmails := make([]string, len(clients))
	for i, client := range clients {
		actualEmails[i] = client.EmailString()
	}
	
	for _, expectedEmail := range expectedEmails {
		assert.Contains(t, actualEmails, expectedEmail)
	}
	
	// Verify that client details are properly deserialized from PostgreSQL
	for _, client := range clients {
		if client.EmailString() == "john@example.com" {
			assert.Equal(t, client1.ID(), client.ID())
			assert.Equal(t, "John Doe", client.Name())
			assert.Equal(t, "+1234567890", client.PhoneString())
			assert.Equal(t, "123 Main St", client.Address())
		} else if client.EmailString() == "jane@example.com" {
			assert.Equal(t, client2.ID(), client.ID())
			assert.Equal(t, "Jane Smith", client.Name())
			assert.Equal(t, "+0987654321", client.PhoneString())
			assert.Equal(t, "456 Oak Ave", client.Address())
		}
	}
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