package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

func TestClientHandler_ListClients_GET_EmptyList(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	handler.ListClients(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	// Check response body contains empty array
	assert.Contains(t, rr.Body.String(), `"data":[]`)
	assert.Contains(t, rr.Body.String(), `"success":true`)
}

func TestClientHandler_ListClients_GET_WithClients(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	// Create test clients
	_, err := billingService.CreateClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)
	
	_, err = billingService.CreateClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	handler.ListClients(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	// Check response contains both clients
	responseBody := rr.Body.String()
	assert.Contains(t, responseBody, `"success":true`)
	assert.Contains(t, responseBody, "john@example.com")
	assert.Contains(t, responseBody, "jane@example.com")
	assert.Contains(t, responseBody, "John Doe")
	assert.Contains(t, responseBody, "Jane Smith")
}

func TestClientHandler_ListClients_MethodNotAllowed(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	handler.ListClients(rr, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	// Check error response
	responseBody := rr.Body.String()
	assert.Contains(t, responseBody, `"success":false`)
	assert.Contains(t, responseBody, "METHOD_NOT_ALLOWED")
}

func TestClientHandler_ListClients_PUT_MethodNotAllowed(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	billingService := application.NewBillingService(clientRepo)
	handler := handlers.NewClientHandler(billingService)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clients", nil)
	rr := httptest.NewRecorder()

	// Act
	handler.ListClients(rr, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}