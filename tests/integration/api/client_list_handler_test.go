package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

func TestClientHandler_ListClients_IntegrationTest(t *testing.T) {
	// Arrange - Set up integration test server with PostgreSQL storage
	stack := testhelpers.NewCleanIntegrationTestStack()
	server := stack.HTTPServer
	
	// Create test clients via service
	_, err := stack.BillingService.CreateClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)
	
	_, err = stack.BillingService.CreateClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
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
		if clientMap["email"] == "john@example.com" {
			assert.Equal(t, "John Doe", clientMap["name"])
			assert.Equal(t, "+1234567890", clientMap["phone"])
			assert.Equal(t, "123 Main St", clientMap["address"])
		} else if clientMap["email"] == "jane@example.com" {
			assert.Equal(t, "Jane Smith", clientMap["name"])
			assert.Equal(t, "+0987654321", clientMap["phone"])
			assert.Equal(t, "456 Oak Ave", clientMap["address"])
		}
	}
}

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