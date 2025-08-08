package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

// BUSINESS_TITLE: Retrieve Individual Client Information
// BUSINESS_DESCRIPTION: Customer service agents and sales representatives can look up specific client details using client ID for support and relationship management
// USER_STORY: As a customer service agent, I want to retrieve a client's complete information so that I can provide personalized support
// BUSINESS_VALUE: Enables efficient customer support, improves service quality, supports relationship management workflows
// SCENARIOS_TESTED: Valid client lookup, handling of non-existent clients, invalid ID format protection
func TestClientHandler_GetClient_Success(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	validScenario := scenarios[0] // "Valid Get Client Scenario"

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	// Create a client first
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
	require.NoError(t, err)

	err = server.ClientRepository.Save(client)
	require.NoError(t, err)

	// Test GET request
	url := fmt.Sprintf("/api/v1/clients/%s", validScenario.Client.ID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response dtos.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.True(t, response.Success, "Response should indicate success")

	// Verify client data in response
	clientData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response data should be client object")
	assert.Equal(t, validScenario.Client.ID, clientData["id"], "Client ID should match")
	assert.Equal(t, validScenario.Client.Name, clientData["name"], "Client name should match")
	assert.Equal(t, validScenario.Client.Email, clientData["email"], "Client email should match")
}

// BUSINESS_TITLE: Handle Non-existent Client Lookup
// BUSINESS_DESCRIPTION: System properly handles requests for clients that don't exist, providing clear feedback to customer service agents
// USER_STORY: As a customer service agent, I want clear feedback when a client ID doesn't exist so that I can inform customers appropriately
// BUSINESS_VALUE: Prevents confusion, improves customer service quality, reduces support escalations
// SCENARIOS_TESTED: Non-existent client ID returns 404, clear error message provided
func TestClientHandler_GetClient_NotFound(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	nonExistentID := scenarios[3].NonExistentIDs[0] // First non-existent ID

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	// Test GET request with non-existent ID
	url := fmt.Sprintf("/api/v1/clients/%s", nonExistentID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 Not Found")

	var response dtos.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.False(t, response.Success, "Response should indicate failure")
	assert.Equal(t, "REPOSITORY_NOT_FOUND", response.Error.Code, "Error code should be REPOSITORY_NOT_FOUND")
}

// BUSINESS_TITLE: Invalid Client ID Format Protection
// BUSINESS_DESCRIPTION: System protects against malformed client ID requests, preventing system errors and providing clear feedback
// USER_STORY: As a customer service agent, I want clear error messages when I enter an invalid client ID format so that I can correct my input
// BUSINESS_VALUE: Prevents system errors, improves user experience, guides users toward correct input format
// SCENARIOS_TESTED: Various invalid UUID formats return 400 Bad Request with helpful error messages
func TestClientHandler_GetClient_InvalidUUID(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	invalidIDs := scenarios[2].InvalidIDs // Invalid UUID scenarios

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	for _, invalidID := range invalidIDs {
		t.Run("InvalidID_"+invalidID, func(t *testing.T) {
			// Test GET request with invalid ID
			url := fmt.Sprintf("/api/v1/clients/%s", invalidID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.RemoteAddr = "192.0.2.1:1234"

			w := httptest.NewRecorder()
			server.HTTPHandler.ServeHTTP(w, req)

			// Assertions - this should FAIL until implemented
			assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request for invalid ID: %s", invalidID)

			var response dtos.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.False(t, response.Success, "Response should indicate failure")
			assert.Contains(t, response.Error.Code, "VALIDATION", "Error code should indicate validation error")
		})
	}
}

// loadGetClientScenarios loads test scenarios from JSON file
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
