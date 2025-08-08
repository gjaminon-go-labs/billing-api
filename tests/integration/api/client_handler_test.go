// API Handler Integration Tests - Client Operations
//
// This file contains integration tests for HTTP API handlers.
// Tests: HTTP request/response handling, DTO conversion, error mapping, handler integration
// Scope: Integration tests - HTTP Handler + Application Service + Repository + InMemory Storage
// Use Cases: UC-B-001 (Create Client) - API presentation layer
//
// Test Scenarios:
// - HTTP request processing and response generation
// - DTO conversion (HTTP JSON ↔ Domain models)
// - HTTP status code mapping from domain errors
// - Content-Type handling and JSON serialization
// - Method validation (POST required)
// - Invalid JSON handling
// - Uses external JSON test data for HTTP request scenarios
//
// Components Tested:
// - ClientHandler (API layer)
// - BillingService (application layer)
// - ClientRepository (repository pattern)
// - PostgreSQL Storage (test infrastructure)
package api

import (
	"bytes"
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

type HTTPTestCase struct {
	Description       string                   `json:"description"`
	RequestBody       dtos.CreateClientRequest `json:"request_body"`
	ExpectedStatus    int                      `json:"expected_status"`
	ShouldSucceed     bool                     `json:"should_succeed"`
	ExpectedErrorCode string                   `json:"expected_error_code,omitempty"`
}

// BUSINESS_TITLE: Create New Client via API
// BUSINESS_DESCRIPTION: Sales representatives and customer service agents can add new clients through the web interface, ensuring all client data is properly validated and stored
// USER_STORY: As a sales representative, I want to create new client records through a web form so that I can organize and track my customer relationships
// BUSINESS_VALUE: Enables customer onboarding, relationship management, and sales tracking. Critical for business growth and customer data organization
// SCENARIOS_TESTED: Valid client creation, data validation (email format, required fields), error handling for invalid data, duplicate prevention
func TestClientHandler_CreateClient(t *testing.T) {
	// Load test data
	testCases := loadHTTPTestCases(t)

	// Set up integration test server with PostgreSQL storage
	server := testhelpers.NewIntegrationTestServer()

	// Test each scenario
	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			// Create HTTP request
			requestBody, err := json.Marshal(testCase.RequestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler through server
			server.Handler().ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, testCase.ExpectedStatus, rr.Code, "Status code mismatch for: %s", testCase.Description)

			// Parse response
			var responseBody map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.NoError(t, err, "Failed to parse response JSON for: %s", testCase.Description)

			if testCase.ShouldSucceed {
				// Check success response structure
				assert.True(t, responseBody["success"].(bool), "Response should indicate success")
				assert.Contains(t, responseBody, "data", "Success response should contain data")

				// Check client data structure
				data := responseBody["data"].(map[string]interface{})
				assert.Contains(t, data, "id", "Client data should contain ID")
				assert.Contains(t, data, "name", "Client data should contain name")
				assert.Contains(t, data, "email", "Client data should contain email")
				assert.Contains(t, data, "created_at", "Client data should contain created_at")
				assert.Contains(t, data, "updated_at", "Client data should contain updated_at")

				// Verify data matches request
				assert.Equal(t, testCase.RequestBody.Name, data["name"])
				assert.Equal(t, testCase.RequestBody.Email, data["email"])
			} else {
				// Check error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.Contains(t, responseBody, "error", "Error response should contain error")

				// Check error structure
				errorDetail := responseBody["error"].(map[string]interface{})
				assert.Contains(t, errorDetail, "code", "Error should contain code")
				assert.Contains(t, errorDetail, "message", "Error should contain message")

				// Check expected error code if specified
				if testCase.ExpectedErrorCode != "" {
					assert.Equal(t, testCase.ExpectedErrorCode, errorDetail["code"])
				}
			}
		})
	}
}

// BUSINESS_TITLE: API Security - Method Validation
// BUSINESS_DESCRIPTION: System prevents unauthorized API calls and ensures only proper HTTP methods are accepted for client creation
// USER_STORY: As a system administrator, I want the API to reject invalid requests so that the application remains secure and stable
// BUSINESS_VALUE: Protects against malicious requests, ensures API consistency, and maintains system security
// SCENARIOS_TESTED: Invalid HTTP methods (PUT instead of POST), proper error messages, security boundaries
func TestClientHandler_CreateClient_MethodNotAllowed(t *testing.T) {
	// Set up integration test server with PostgreSQL storage
	server := testhelpers.NewIntegrationTestServer()

	// Test PUT method (should be method not allowed)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// Parse response
	var responseBody dtos.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Check error response
	assert.False(t, responseBody.Success)
	assert.Equal(t, "METHOD_NOT_ALLOWED", responseBody.Error.Code)
}

// BUSINESS_TITLE: Data Validation - Invalid Request Format
// BUSINESS_DESCRIPTION: System properly handles malformed data submissions and provides clear error messages to users
// USER_STORY: As a user, I want to receive clear error messages when I submit invalid data so that I can correct my input
// BUSINESS_VALUE: Prevents data corruption, improves user experience, reduces support requests
// SCENARIOS_TESTED: Malformed JSON requests, clear error responses, graceful error handling
func TestClientHandler_CreateClient_InvalidJSON(t *testing.T) {
	// Set up integration test server with PostgreSQL storage
	server := testhelpers.NewIntegrationTestServer()

	// Test invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse response
	var responseBody dtos.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Check error response
	assert.False(t, responseBody.Success)
	assert.Equal(t, "INVALID_JSON", responseBody.Error.Code)
}

func loadHTTPTestCases(t *testing.T) []HTTPTestCase {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to HTTP test data
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "http", "create_client_requests.json")

	// Read test data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read HTTP test data file")

	// Parse JSON
	var testCases []HTTPTestCase
	err = json.Unmarshal(data, &testCases)
	assert.NoError(t, err, "Failed to parse HTTP test data JSON")

	return testCases
}
