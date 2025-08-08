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
	httpserver "github.com/gjaminon-go-labs/billing-api/internal/api/http"
	"github.com/gjaminon-go-labs/billing-api/internal/di"
	"github.com/gjaminon-go-labs/billing-api/tests/testdata"
)

func TestClientHandler_ListClients_Pagination_IntegrationTest(t *testing.T) {
	// Setup test container with PostgreSQL
	container := di.NewContainer(di.IntegrationTestConfig())
	
	// Get services
	billingService, err := container.GetBillingService()
	require.NoError(t, err)
	
	// Create HTTP server
	server := httpserver.NewServer(billingService)
	handler := server.Handler()
	
	// Setup: Create 25 test clients
	for i := 1; i <= 25; i++ {
		clientData := fmt.Sprintf(`{
			"name": "Test Client %d",
			"email": "client%d@test.com",
			"phone": "+123456789%d",
			"address": "Address %d"
		}`, i, i, i%10, i)
		
		req := httptest.NewRequest("POST", "/api/v1/clients", testdata.JSONReader(clientData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// Test cases
	tests := []struct {
		name               string
		queryParams        string
		expectedStatus     int
		expectedPage       int
		expectedLimit      int
		expectedDataCount  int
		expectedTotalCount int
		expectedTotalPages int
	}{
		{
			name:               "First page with default limit",
			queryParams:        "?page=1&limit=10",
			expectedStatus:     http.StatusOK,
			expectedPage:       1,
			expectedLimit:      10,
			expectedDataCount:  10,
			expectedTotalCount: 25,
			expectedTotalPages: 3,
		},
		{
			name:               "Second page",
			queryParams:        "?page=2&limit=10",
			expectedStatus:     http.StatusOK,
			expectedPage:       2,
			expectedLimit:      10,
			expectedDataCount:  10,
			expectedTotalCount: 25,
			expectedTotalPages: 3,
		},
		{
			name:               "Last page with partial results",
			queryParams:        "?page=3&limit=10",
			expectedStatus:     http.StatusOK,
			expectedPage:       3,
			expectedLimit:      10,
			expectedDataCount:  5,
			expectedTotalCount: 25,
			expectedTotalPages: 3,
		},
		{
			name:               "Large limit gets all results",
			queryParams:        "?page=1&limit=50",
			expectedStatus:     http.StatusOK,
			expectedPage:       1,
			expectedLimit:      50,
			expectedDataCount:  25,
			expectedTotalCount: 25,
			expectedTotalPages: 1,
		},
		{
			name:               "Page beyond available data",
			queryParams:        "?page=10&limit=10",
			expectedStatus:     http.StatusOK,
			expectedPage:       10,
			expectedLimit:      10,
			expectedDataCount:  0,
			expectedTotalCount: 25,
			expectedTotalPages: 3,
		},
		{
			name:           "Invalid page parameter",
			queryParams:    "?page=0&limit=10",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Limit exceeds maximum",
			queryParams:    "?page=1&limit=101",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:               "No parameters uses defaults",
			queryParams:        "",
			expectedStatus:     http.StatusOK,
			expectedPage:       1,
			expectedLimit:      20,
			expectedDataCount:  20,
			expectedTotalCount: 25,
			expectedTotalPages: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/v1/clients"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			
			// Execute
			handler.ServeHTTP(rec, req)
			
			// Assert status
			assert.Equal(t, tt.expectedStatus, rec.Code)
			
			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var response struct {
					Data       []dtos.ClientResponse `json:"data"`
					Pagination struct {
						Page       int `json:"page"`
						Limit      int `json:"limit"`
						TotalCount int `json:"total_count"`
						TotalPages int `json:"total_pages"`
					} `json:"pagination"`
					Success bool `json:"success"`
				}
				
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				
				// Assert pagination metadata
				assert.Equal(t, tt.expectedPage, response.Pagination.Page)
				assert.Equal(t, tt.expectedLimit, response.Pagination.Limit)
				assert.Equal(t, tt.expectedTotalCount, response.Pagination.TotalCount)
				assert.Equal(t, tt.expectedTotalPages, response.Pagination.TotalPages)
				
				// Assert data count
				assert.Len(t, response.Data, tt.expectedDataCount)
				
				// Verify data ordering (should be consistent)
				if len(response.Data) > 1 {
					// Check that results are ordered
					for i := 1; i < len(response.Data); i++ {
						prevName := response.Data[i-1].Name
						currName := response.Data[i].Name
						assert.True(t, prevName < currName, 
							"Results should be ordered: %s should come before %s", prevName, currName)
					}
				}
			}
		})
	}
	
	// Cleanup
	container.Reset()
}

func TestClientHandler_Pagination_Consistency_IntegrationTest(t *testing.T) {
	// Setup test container
	container := di.NewContainer(di.IntegrationTestConfig())
	billingService, err := container.GetBillingService()
	require.NoError(t, err)
	
	server := httpserver.NewServer(billingService)
	handler := server.Handler()
	
	// Create exactly 15 clients
	clientIDs := make([]string, 15)
	for i := 0; i < 15; i++ {
		clientData := fmt.Sprintf(`{
			"name": "Client %02d",
			"email": "test%d@example.com",
			"phone": "+1234567890",
			"address": "Test Address"
		}`, i, i)
		
		req := httptest.NewRequest("POST", "/api/v1/clients", testdata.JSONReader(clientData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		
		require.Equal(t, http.StatusCreated, rec.Code)
		
		var response struct {
			Data dtos.ClientResponse `json:"data"`
		}
		json.NewDecoder(rec.Body).Decode(&response)
		clientIDs[i] = response.Data.ID
	}
	
	// Test: Fetch all pages and verify no duplicates and no missing items
	allClients := make(map[string]bool)
	pageSize := 5
	
	for page := 1; page <= 3; page++ {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/clients?page=%d&limit=%d", page, pageSize), nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var response struct {
			Data []dtos.ClientResponse `json:"data"`
			Pagination struct {
				Page       int `json:"page"`
				Limit      int `json:"limit"`
				TotalCount int `json:"total_count"`
				TotalPages int `json:"total_pages"`
			} `json:"pagination"`
		}
		
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		
		// Verify pagination metadata
		assert.Equal(t, page, response.Pagination.Page)
		assert.Equal(t, pageSize, response.Pagination.Limit)
		assert.Equal(t, 15, response.Pagination.TotalCount)
		assert.Equal(t, 3, response.Pagination.TotalPages)
		
		// Verify page size
		expectedCount := pageSize
		if page == 3 {
			expectedCount = 5 // Last page has only 5 items
		}
		assert.Len(t, response.Data, expectedCount)
		
		// Track all clients to check for duplicates
		for _, client := range response.Data {
			if allClients[client.ID] {
				t.Errorf("Duplicate client found: %s", client.ID)
			}
			allClients[client.ID] = true
		}
	}
	
	// Verify we got all 15 clients exactly once
	assert.Len(t, allClients, 15, "Should have received all 15 clients across pages")
	
	// Cleanup
	container.Reset()
}