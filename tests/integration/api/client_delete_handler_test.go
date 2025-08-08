package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

// BUSINESS_TITLE: Remove Client Records
// BUSINESS_DESCRIPTION: Administrative users can remove client records when accounts are closed or data needs to be purged per privacy regulations
// USER_STORY: As an administrator, I want to delete client records so that inactive accounts don't clutter the system and privacy requirements are met
// BUSINESS_VALUE: Maintains clean data, supports compliance with privacy regulations, improves system performance
// SCENARIOS_TESTED: Successful deletion, non-existent client handling, proper cleanup verification
func TestClientHandler_DeleteClient_Success(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	validScenario := scenarios[0] // "Valid Get Client Scenario"

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	// Create a client first
	client, err := entity.NewClientWithID(
		validScenario.Client.ID,
		validScenario.Client.Name,
		validScenario.Client.Email,
		validScenario.Client.Phone,
		validScenario.Client.Address,
	)
	require.NoError(t, err)

	err = server.ClientRepository.Save(client)
	require.NoError(t, err)

	// Verify client exists before deletion
	existingClient, err := server.ClientRepository.GetByID(validScenario.Client.ID)
	require.NoError(t, err, "Client should exist before deletion")
	require.NotNil(t, existingClient, "Client should exist before deletion")

	// Test DELETE request
	url := fmt.Sprintf("/api/v1/clients/%s", validScenario.Client.ID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.RemoteAddr = "192.0.2.1:1234"

	w := httptest.NewRecorder()
	server.HTTPHandler.ServeHTTP(w, req)

	// Assertions - this should FAIL until implemented
	assert.Equal(t, http.StatusNoContent, w.Code, "Should return 204 No Content")
	assert.Empty(t, w.Body.String(), "Response body should be empty for 204")

	// Verify client no longer exists
	deletedClient, err := server.ClientRepository.GetByID(validScenario.Client.ID)
	assert.Error(t, err, "Client should not exist after deletion")
	assert.Nil(t, deletedClient, "Client should be nil after deletion")
}

// BUSINESS_TITLE: Handle Non-existent Client Deletion
// BUSINESS_DESCRIPTION: System properly handles deletion requests for clients that don't exist, providing appropriate feedback
// USER_STORY: As an administrator, I want clear feedback when trying to delete a client that doesn't exist so that I can verify the operation
// BUSINESS_VALUE: Prevents confusion, provides clear system feedback, maintains operation consistency
// SCENARIOS_TESTED: Non-existent client ID deletion returns 404 with clear error message
func TestClientHandler_DeleteClient_NotFound(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	nonExistentID := scenarios[3].NonExistentIDs[0] // First non-existent ID

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	// Test DELETE request with non-existent ID
	url := fmt.Sprintf("/api/v1/clients/%s", nonExistentID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
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

// BUSINESS_TITLE: Invalid Client ID Format Protection for Deletion
// BUSINESS_DESCRIPTION: System protects against malformed client ID deletion requests, preventing system errors
// USER_STORY: As an administrator, I want clear error messages when I enter an invalid client ID format for deletion so that I can correct my input
// BUSINESS_VALUE: Prevents system errors, improves user experience, guides users toward correct input format
// SCENARIOS_TESTED: Various invalid UUID formats for deletion return 400 Bad Request
func TestClientHandler_DeleteClient_InvalidUUID(t *testing.T) {
	// Load test scenarios
	scenarios := loadGetClientScenarios(t)
	invalidIDs := scenarios[2].InvalidIDs // Invalid UUID scenarios

	// Setup integration test server
	server := testhelpers.NewIntegrationTestServer(t)
	defer server.Close()

	for _, invalidID := range invalidIDs {
		t.Run("InvalidID_"+invalidID, func(t *testing.T) {
			// Test DELETE request with invalid ID
			url := fmt.Sprintf("/api/v1/clients/%s", invalidID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
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
