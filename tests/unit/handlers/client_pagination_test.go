package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

func TestClientHandler_ListClients_WithPagination(t *testing.T) {
	tests := []struct {
		name               string
		queryParams        string
		setupClients       int
		expectedStatus     int
		expectedPage       int
		expectedLimit      int
		expectedDataCount  int
		expectedTotalCount int
		expectedTotalPages int
		expectedError      string
	}{
		{
			name:               "Valid pagination - first page",
			queryParams:        "?page=1&limit=5",
			setupClients:       12,
			expectedStatus:     http.StatusOK,
			expectedPage:       1,
			expectedLimit:      5,
			expectedDataCount:  5,
			expectedTotalCount: 12,
			expectedTotalPages: 3,
		},
		{
			name:               "Valid pagination - second page",
			queryParams:        "?page=2&limit=5",
			setupClients:       12,
			expectedStatus:     http.StatusOK,
			expectedPage:       2,
			expectedLimit:      5,
			expectedDataCount:  5,
			expectedTotalCount: 12,
			expectedTotalPages: 3,
		},
		{
			name:               "Valid pagination - last page with partial results",
			queryParams:        "?page=3&limit=5",
			setupClients:       12,
			expectedStatus:     http.StatusOK,
			expectedPage:       3,
			expectedLimit:      5,
			expectedDataCount:  2,
			expectedTotalCount: 12,
			expectedTotalPages: 3,
		},
		{
			name:               "Default pagination when no params",
			queryParams:        "",
			setupClients:       25,
			expectedStatus:     http.StatusOK,
			expectedPage:       1,
			expectedLimit:      20,
			expectedDataCount:  20,
			expectedTotalCount: 25,
			expectedTotalPages: 2,
		},
		{
			name:           "Invalid page number - zero",
			queryParams:    "?page=0&limit=10",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "page must be greater than 0",
		},
		{
			name:           "Invalid page number - negative",
			queryParams:    "?page=-1&limit=10",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "page must be greater than 0",
		},
		{
			name:           "Invalid limit - zero",
			queryParams:    "?page=1&limit=0",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "limit must be between 1 and 100",
		},
		{
			name:           "Limit exceeds maximum",
			queryParams:    "?page=1&limit=101",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "limit must be between 1 and 100",
		},
		{
			name:           "Invalid page format",
			queryParams:    "?page=abc&limit=10",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid page parameter",
		},
		{
			name:           "Invalid limit format",
			queryParams:    "?page=1&limit=xyz",
			setupClients:   5,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid limit parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := infrastructure.NewInMemoryStorage()
			clientRepo := repository.NewClientRepository(storage)
			billingService := application.NewBillingService(clientRepo)
			handler := handlers.NewClientHandler(billingService)

			// Create test clients
			for i := 0; i < tt.setupClients; i++ {
				createReq := dtos.CreateClientRequest{
					Name:    fmt.Sprintf("Client %02d", i),
					Email:   fmt.Sprintf("client%d@test.com", i),
					Phone:   "+1234567890",
					Address: fmt.Sprintf("Address %d", i),
				}
				_, err := billingService.CreateClient(createReq)
				require.NoError(t, err)
			}

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/clients"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			// Execute
			handler.ListClients(rec, req)

			// Assert status
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}

			if tt.expectedStatus == http.StatusOK {
				// Parse response with pagination
				var response struct {
					Data       []dtos.ClientResponse `json:"data"`
					Pagination *struct {
						Page       int `json:"page"`
						Limit      int `json:"limit"`
						TotalCount int `json:"total_count"`
						TotalPages int `json:"total_pages"`
					} `json:"pagination,omitempty"`
					Success bool `json:"success"`
				}

				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)

				// Check if pagination is present
				if tt.queryParams != "" || tt.setupClients > 20 {
					// Should have pagination metadata
					require.NotNil(t, response.Pagination, "Pagination metadata should be present")
					assert.Equal(t, tt.expectedPage, response.Pagination.Page)
					assert.Equal(t, tt.expectedLimit, response.Pagination.Limit)
					assert.Equal(t, tt.expectedTotalCount, response.Pagination.TotalCount)
					assert.Equal(t, tt.expectedTotalPages, response.Pagination.TotalPages)
				}

				// Verify data count
				assert.Len(t, response.Data, tt.expectedDataCount)
			}
		})
	}
}

func TestBillingService_ListClientsWithPagination(t *testing.T) {
	tests := []struct {
		name               string
		page               int
		limit              int
		totalClients       int
		expectedCount      int
		expectedTotalPages int
	}{
		{
			name:               "First page of results",
			page:               1,
			limit:              10,
			totalClients:       25,
			expectedCount:      10,
			expectedTotalPages: 3,
		},
		{
			name:               "Second page of results",
			page:               2,
			limit:              10,
			totalClients:       25,
			expectedCount:      10,
			expectedTotalPages: 3,
		},
		{
			name:               "Last page with partial results",
			page:               3,
			limit:              10,
			totalClients:       25,
			expectedCount:      5,
			expectedTotalPages: 3,
		},
		{
			name:               "Page beyond available data",
			page:               5,
			limit:              10,
			totalClients:       25,
			expectedCount:      0,
			expectedTotalPages: 3,
		},
		{
			name:               "Single page with all results",
			page:               1,
			limit:              50,
			totalClients:       25,
			expectedCount:      25,
			expectedTotalPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := infrastructure.NewInMemoryStorage()
			clientRepo := repository.NewClientRepository(storage)
			service := application.NewBillingService(clientRepo)

			// Create test clients
			for i := 0; i < tt.totalClients; i++ {
				createReq := dtos.CreateClientRequest{
					Name:    fmt.Sprintf("Client %03d", i),
					Email:   fmt.Sprintf("client%d@test.com", i),
					Phone:   "+1234567890",
					Address: fmt.Sprintf("Address %d", i),
				}
				_, err := service.CreateClient(createReq)
				require.NoError(t, err)
			}

			// Execute
			result, err := service.ListClientsWithPagination(tt.page, tt.limit)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.page, result.Pagination.Page)
			assert.Equal(t, tt.limit, result.Pagination.Limit)
			assert.Equal(t, tt.totalClients, result.Pagination.TotalCount)
			assert.Equal(t, tt.expectedTotalPages, result.Pagination.TotalPages)
			assert.Len(t, result.Clients, tt.expectedCount)
		})
	}
}