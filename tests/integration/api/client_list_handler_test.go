package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

// BUSINESS_TITLE: View All Clients
// BUSINESS_DESCRIPTION: Sales representatives and managers can view a complete list of all clients in the system for planning and relationship management
// USER_STORY: As a sales manager, I want to see all clients in the system so that I can track our customer base and assign accounts to team members
// BUSINESS_VALUE: Provides visibility into customer portfolio, enables territory management, supports sales planning and customer relationship oversight
// SCENARIOS_TESTED: Retrieving client lists with data, proper client information display, handling multiple clients
func TestClientHandler_ListClients_IntegrationTest(t *testing.T) {
	// Arrange - Set up integration test server with PostgreSQL storage
	stack := testhelpers.NewCleanIntegrationTestStack()
	server := stack.HTTPServer

	// Load test fixtures
	fixtures := loadAPITestFixtures(t)

	// Create test clients via service from fixtures
	_, err := stack.BillingService.CreateClient(fixtures[0].Name, fixtures[0].Email, fixtures[0].Phone, fixtures[0].Address)
	assert.NoError(t, err)

	_, err = stack.BillingService.CreateClient(fixtures[1].Name, fixtures[1].Email, fixtures[1].Phone, fixtures[1].Address)
	assert.NoError(t, err)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	server.Handler().ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse response
	var response dtos.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check success response structure
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Check that data is an array of clients
	clientsData, ok := response.Data.([]interface{})
	assert.True(t, ok, "Data should be an array")
	assert.Len(t, clientsData, 2)

	// Verify client data structure
	for _, clientData := range clientsData {
		clientMap, ok := clientData.(map[string]interface{})
		assert.True(t, ok, "Each client should be a map")

		// Check required fields
		assert.Contains(t, clientMap, "id")
		assert.Contains(t, clientMap, "name")
		assert.Contains(t, clientMap, "email")
		assert.Contains(t, clientMap, "created_at")
		assert.Contains(t, clientMap, "updated_at")

		// Check one of our test clients
		if clientMap["email"] == fixtures[0].Email {
			assert.Equal(t, fixtures[0].Name, clientMap["name"])
			assert.Equal(t, fixtures[0].Phone, clientMap["phone"])
			assert.Equal(t, fixtures[0].Address, clientMap["address"])
		} else if clientMap["email"] == fixtures[1].Email {
			assert.Equal(t, fixtures[1].Name, clientMap["name"])
			assert.Equal(t, fixtures[1].Phone, clientMap["phone"])
			assert.Equal(t, fixtures[1].Address, clientMap["address"])
		}
	}
}

// BUSINESS_TITLE: Empty Client List Handling
// BUSINESS_DESCRIPTION: System gracefully handles scenarios where no clients exist yet, providing clear messaging for new users or empty databases
// USER_STORY: As a new user setting up the system, I want to see an appropriate message when no clients exist so that I understand the system is working correctly
// BUSINESS_VALUE: Improves user experience for new system deployments, prevents confusion about empty states, guides users toward their first actions
// SCENARIOS_TESTED: Empty database states, proper empty list responses, user-friendly empty state handling
func TestClientHandler_ListClients_EmptyList_IntegrationTest(t *testing.T) {
	// Arrange - Set up integration test server with clean PostgreSQL storage
	stack := testhelpers.NewCleanIntegrationTestStack()
	server := stack.HTTPServer

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	server.Handler().ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response dtos.SuccessResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check success response with empty array
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Check that data is an empty array
	clientsData, ok := response.Data.([]interface{})
	assert.True(t, ok, "Data should be an array")
	assert.Empty(t, clientsData)
}

type ClientFixture struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadAPITestFixtures(t *testing.T) []ClientFixture {
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
