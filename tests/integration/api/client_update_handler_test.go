package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

// BUSINESS_TITLE: Update Client Information
// BUSINESS_DESCRIPTION: Sales and service teams can modify client details when contact information changes or corrections are needed
// USER_STORY: As a sales representative, I want to update client contact information so that communication remains accurate and effective
// BUSINESS_VALUE: Maintains data accuracy, improves communication reliability, supports customer relationship management
// SCENARIOS_TESTED: Complete updates, partial updates, validation enforcement, timestamp management
func TestClientHandler_UpdateClient_Success(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	fullUpdateScenario := scenarios[0] // "Full Update Request"

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)

	// Create the original client
	originalClient, err := entity.NewClientWithID(
		fullUpdateScenario.ExpectedClient.ID,
		"Alice Johnson", // Original name
		fullUpdateScenario.ExpectedClient.Email,
		"+1234567890",                        // Original phone
		"123 Main Street, Anytown, ST 12345", // Original address
	)
	require.NoError(t, err)

	err = server.ClientRepository.Save(originalClient)
	require.NoError(t, err)

	// Prepare update request
	requestBody, err := json.Marshal(fullUpdateScenario.Request)
	require.NoError(t, err)

	// Test PUT request
	url := fmt.Sprintf("/api/v1/clients/%s", fullUpdateScenario.ExpectedClient.ID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response dtos.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.True(t, response.Success, "Response should indicate success")

	// Verify updated client data in response
	clientData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response data should be client object")
	assert.Equal(t, fullUpdateScenario.ExpectedClient.ID, clientData["id"], "Client ID should remain unchanged")
	assert.Equal(t, fullUpdateScenario.Request.Name, clientData["name"], "Client name should be updated")
	assert.Equal(t, fullUpdateScenario.Request.Phone, clientData["phone"], "Client phone should be updated")
	assert.Equal(t, fullUpdateScenario.Request.Address, clientData["address"], "Client address should be updated")
	assert.Equal(t, fullUpdateScenario.ExpectedClient.Email, clientData["email"], "Client email should remain unchanged")
}

// BUSINESS_TITLE: Partial Client Information Updates
// BUSINESS_DESCRIPTION: Sales teams can update specific client fields without affecting other information, providing flexibility in data management
// USER_STORY: As a sales representative, I want to update only specific client fields so that I can make targeted corrections without affecting other data
// BUSINESS_VALUE: Enables precise data updates, reduces errors from unnecessary changes, improves data management efficiency
// SCENARIOS_TESTED: Partial field updates, field clearing, unchanged field preservation
func TestClientHandler_UpdateClient_PartialUpdate(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	partialUpdateScenario := scenarios[1] // "Partial Update Request - Name Only"

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)

	// Create the original client
	originalClient, err := entity.NewClientWithID(
		partialUpdateScenario.ExpectedClient.ID,
		"Alice Johnson", // Original name to be updated
		partialUpdateScenario.ExpectedClient.Email,
		"+1234567890",                        // Original phone (should be cleared)
		"123 Main Street, Anytown, ST 12345", // Original address (should be cleared)
	)
	require.NoError(t, err)

	err = server.ClientRepository.Save(originalClient)
	require.NoError(t, err)

	// Prepare partial update request
	requestBody, err := json.Marshal(partialUpdateScenario.Request)
	require.NoError(t, err)

	// Test PUT request
	url := fmt.Sprintf("/api/v1/clients/%s", partialUpdateScenario.ExpectedClient.ID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response dtos.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.True(t, response.Success, "Response should indicate success")

	// Verify updated client data
	clientData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response data should be client object")
	assert.Equal(t, partialUpdateScenario.Request.Name, clientData["name"], "Client name should be updated")
	assert.Equal(t, "", clientData["phone"], "Client phone should be cleared")
	assert.Equal(t, "", clientData["address"], "Client address should be cleared")
}

// BUSINESS_TITLE: Handle Non-existent Client Updates
// BUSINESS_DESCRIPTION: System properly handles update requests for clients that don't exist, providing clear feedback
// USER_STORY: As a sales representative, I want clear feedback when trying to update a client that doesn't exist so that I can verify the client ID
// BUSINESS_VALUE: Prevents confusion, provides clear error messages, improves user experience
// SCENARIOS_TESTED: Non-existent client ID update returns 404 with clear error message
func TestClientHandler_UpdateClient_NotFound(t *testing.T) {
	// Load test scenarios
	getScenarios := loadGetClientScenarios(t)
	nonExistentID := getScenarios[3].NonExistentIDs[0] // First non-existent ID

	updateScenarios := loadUpdateClientScenarios(t)
	updateRequest := updateScenarios[0].Request // Any valid update request

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)

	// Prepare update request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	// Test PUT request with non-existent ID
	url := fmt.Sprintf("/api/v1/clients/%s", nonExistentID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 Not Found")

	var response dtos.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.False(t, response.Success, "Response should indicate failure")
	assert.Equal(t, "REPOSITORY_NOT_FOUND", response.Error.Code, "Error code should be REPOSITORY_NOT_FOUND")
}

// BUSINESS_TITLE: Client Update Validation Protection
// BUSINESS_DESCRIPTION: System validates client update data to prevent invalid information from corrupting client records
// USER_STORY: As a sales representative, I want clear validation error messages when I submit invalid client updates so that I can correct the data
// BUSINESS_VALUE: Maintains data quality, prevents system errors, guides users toward correct input
// SCENARIOS_TESTED: Various validation errors with specific field feedback
func TestClientHandler_UpdateClient_ValidationError(t *testing.T) {
	// Load test scenarios
	scenarios := loadUpdateClientScenarios(t)
	invalidRequests := scenarios[3].InvalidRequests // "Invalid Update Requests"

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)

	// Create a test client
	validClient, err := entity.NewClientWithID(
		"123e4567-e89b-12d3-a456-426614174000",
		"Test Client",
		"test@example.com",
		"+1234567890",
		"Test Address",
	)
	require.NoError(t, err)

	err = server.ClientRepository.Save(validClient)
	require.NoError(t, err)

	for _, invalidRequest := range invalidRequests {
		t.Run(invalidRequest.Description, func(t *testing.T) {
			// Prepare invalid update request
			requestBody, err := json.Marshal(invalidRequest.Request)
			require.NoError(t, err)

			// Test PUT request with invalid data
			url := fmt.Sprintf("/api/v1/clients/%s", validClient.ID())
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.RemoteAddr = "192.0.2.1:1234"

			w := httptest.NewRecorder()
			server.HTTPHandler.ServeHTTP(w, req)

			// Assertions - this should FAIL until implemented
			assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request for: %s", invalidRequest.Description)

			var response dtos.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
			assert.False(t, response.Success, "Response should indicate failure")
			assert.Contains(t, response.Error.Code, "VALIDATION", "Error code should indicate validation error")
		})
	}
}

// BUSINESS_TITLE: Invalid JSON Request Handling for Updates
// BUSINESS_DESCRIPTION: System properly handles malformed JSON in client update requests, providing clear error feedback
// USER_STORY: As a sales representative, I want clear error messages when I submit malformed update data so that I can correct the format
// BUSINESS_VALUE: Prevents system errors, improves user experience, guides users toward correct JSON format
// SCENARIOS_TESTED: Malformed JSON returns 400 Bad Request with appropriate error message
func TestClientHandler_UpdateClient_InvalidJSON(t *testing.T) {
	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)

	// Create a test client
	validClient, err := entity.NewClientWithID(
		"123e4567-e89b-12d3-a456-426614174000",
		"Test Client",
		"test@example.com",
		"+1234567890",
		"Test Address",
	)
	require.NoError(t, err)

	err = server.ClientRepository.Save(validClient)
	require.NoError(t, err)

	// Test PUT request with invalid JSON
	invalidJSON := `{"name": "Test", "phone": "+123456789"` // Missing closing brace
	url := fmt.Sprintf("/api/v1/clients/%s", validClient.ID())
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request")

	var response dtos.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.False(t, response.Success, "Response should indicate failure")
	assert.Equal(t, "INVALID_JSON", response.Error.Code, "Error code should be INVALID_JSON")
}

// loadUpdateClientScenarios loads test scenarios from JSON file
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
