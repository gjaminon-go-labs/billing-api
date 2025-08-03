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
// - InMemoryStorage (test infrastructure)
package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

type HTTPTestCase struct {
	Description        string                  `json:"description"`
	RequestBody        dtos.CreateClientRequest `json:"request_body"`
	ExpectedStatus     int                     `json:"expected_status"`
	ShouldSucceed      bool                    `json:"should_succeed"`
	ExpectedErrorCode  string                  `json:"expected_error_code,omitempty"`
}

func TestClientHandler_CreateClient(t *testing.T) {
	// Load test data
	testCases := loadHTTPTestCases(t)
	
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

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

			// Call handler
			handler.CreateClient(rr, req)

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

func TestClientHandler_CreateClient_MethodNotAllowed(t *testing.T) {
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	// Test GET method (should be method not allowed)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	handler.CreateClient(rr, req)

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

func TestClientHandler_CreateClient_InvalidJSON(t *testing.T) {
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	// Test invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateClient(rr, req)

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