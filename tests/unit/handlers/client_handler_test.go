package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/stretchr/testify/assert"
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

	// Load test fixtures
	fixtures := loadHandlerTestFixtures(t)

	// Create test clients from fixtures
	_, err := billingService.CreateClientLegacy(fixtures[0].Name, fixtures[0].Email, fixtures[0].Phone, fixtures[0].Address)
	assert.NoError(t, err)

	_, err = billingService.CreateClientLegacy(fixtures[1].Name, fixtures[1].Email, fixtures[1].Phone, fixtures[1].Address)
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
	assert.Contains(t, responseBody, fixtures[0].Email)
	assert.Contains(t, responseBody, fixtures[1].Email)
	assert.Contains(t, responseBody, fixtures[0].Name)
	assert.Contains(t, responseBody, fixtures[1].Name)
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

type ClientFixture struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadHandlerTestFixtures(t *testing.T) []ClientFixture {
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
