// Full HTTP Server Integration Tests - Client Operations
//
// This file contains end-to-end HTTP integration tests for client use cases.
// Tests: Complete HTTP request/response cycle, full server stack, real network calls
// Scope: Integration tests - Complete HTTP Server + All components + Real HTTP requests
// Use Cases: UC-B-001 (Create Client) - End-to-end HTTP workflow
//
// Test Scenarios:
// - End-to-end HTTP POST requests with real network calls
// - Complete server routing and middleware stack
// - Request/response JSON processing 
// - Multi-request persistence across HTTP calls
// - Success and failure response structure validation
// - Uses external JSON test data for comprehensive scenarios
//
// Components Tested:
// - Complete HTTP Server (with routing)
// - Middleware (CORS, logging, error handling)
// - ClientHandler (API layer)
// - BillingService (application layer)
// - ClientRepository (repository pattern) 
// - InMemoryStorage (test infrastructure)
// - Test helpers (NewInMemoryTestServer)
package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

func TestHTTPServer_Integration_CreateClient(t *testing.T) {
	// Load test data
	testCases := loadHTTPIntegrationTestCases(t)
	
	// Set up complete HTTP server using InMemory test helpers
	server := testhelpers.NewInMemoryTestServer()
	
	// Create test server
	testServer := httptest.NewServer(server.Handler())
	defer testServer.Close()

	// Test each scenario
	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			// Create HTTP request
			requestBody, err := json.Marshal(testCase.RequestBody)
			assert.NoError(t, err)

			// Make actual HTTP request to test server
			resp, err := http.Post(testServer.URL+"/api/v1/clients", "application/json", bytes.NewReader(requestBody))
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, testCase.ExpectedStatus, resp.StatusCode, "Status code mismatch for: %s", testCase.Description)

			// Check content type
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			// Parse response
			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err, "Failed to parse response JSON for: %s", testCase.Description)

			if testCase.ShouldSucceed {
				// Check success response structure
				assert.True(t, responseBody["success"].(bool), "Response should indicate success")
				assert.Contains(t, responseBody, "data", "Success response should contain data")
				
				// Check client data structure
				data := responseBody["data"].(map[string]interface{})
				assert.Contains(t, data, "id", "Client data should contain ID")
				assert.NotEmpty(t, data["id"], "Client ID should not be empty")
				
				// Verify data matches request
				assert.Equal(t, testCase.RequestBody.Name, data["name"])
				assert.Equal(t, testCase.RequestBody.Email, data["email"])
				if testCase.RequestBody.Phone != "" {
					assert.Equal(t, testCase.RequestBody.Phone, data["phone"])
				}
				if testCase.RequestBody.Address != "" {
					assert.Equal(t, testCase.RequestBody.Address, data["address"])
				}
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

func TestHTTPServer_Integration_PersistenceAcrossRequests(t *testing.T) {
	// Set up complete HTTP server using InMemory test helpers (shared storage)
	server := testhelpers.NewInMemoryTestServer()
	
	// Create test server
	testServer := httptest.NewServer(server.Handler())
	defer testServer.Close()

	// Create first client
	firstClient := dtos.CreateClientRequest{
		Name:  "First Client",
		Email: "first@example.com",
		Phone: "+1234567890",
	}
	requestBody, _ := json.Marshal(firstClient)
	
	resp1, err := http.Post(testServer.URL+"/api/v1/clients", "application/json", bytes.NewReader(requestBody))
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Create second client  
	secondClient := dtos.CreateClientRequest{
		Name:  "Second Client",
		Email: "second@example.com",
		Phone: "+0987654321",
	}
	requestBody, _ = json.Marshal(secondClient)
	
	resp2, err := http.Post(testServer.URL+"/api/v1/clients", "application/json", bytes.NewReader(requestBody))
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Verify both clients got different IDs (stored separately)
	var response1, response2 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&response1)
	json.NewDecoder(resp2.Body).Decode(&response2)
	
	data1 := response1["data"].(map[string]interface{})
	data2 := response2["data"].(map[string]interface{})
	
	assert.NotEqual(t, data1["id"], data2["id"], "Different clients should have different IDs")
	assert.Equal(t, "First Client", data1["name"])
	assert.Equal(t, "Second Client", data2["name"])
}