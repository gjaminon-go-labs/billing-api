// HTTP Integration Test Shared Utilities
//
// This file contains shared utilities and helpers for HTTP integration tests.
// Provides: Common test data types, test data loading functions, shared test infrastructure
// Scope: Test utilities - Support for all HTTP integration tests
// Use Cases: Supporting all HTTP test scenarios
//
// Contents:
// - HTTPIntegrationTestCase struct for external test data
// - loadHTTPIntegrationTestCases() function with file path resolution
// - Common test data structures and utilities
// - Shared test helper functions for HTTP testing
//
// Used By:
// - client_integration_test.go (client use case tests)
// - Future invoice_integration_test.go (invoice use case tests)
// - Any HTTP integration tests requiring external test data
package http

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
)

// HTTPIntegrationTestCase represents a test case for HTTP integration tests
type HTTPIntegrationTestCase struct {
	Description        string                  `json:"description"`
	RequestBody        dtos.CreateClientRequest `json:"request_body"`
	ExpectedStatus     int                     `json:"expected_status"`
	ShouldSucceed      bool                    `json:"should_succeed"`
	ExpectedErrorCode  string                  `json:"expected_error_code,omitempty"`
}

// loadHTTPIntegrationTestCases loads test cases from external JSON file
func loadHTTPIntegrationTestCases(t *testing.T) []HTTPIntegrationTestCase {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")
	
	// Build path to HTTP test data
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "http", "create_client_requests.json")
	
	// Read test data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read HTTP integration test data file")
	
	// Parse JSON
	var testCases []HTTPIntegrationTestCase
	err = json.Unmarshal(data, &testCases)
	assert.NoError(t, err, "Failed to parse HTTP integration test data JSON")
	
	return testCases
}