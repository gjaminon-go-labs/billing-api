// Client Domain Unit Tests
//
// This file contains unit tests for Client domain entity validation.
// Tests: Client creation with email validation rules
// Scope: Pure unit tests - single component (Client entity) with no external dependencies
// Use Cases: UC-B-001 (Create Client) - Domain validation layer
//
// Test Scenarios:
// - Valid email addresses (various formats)
// - Invalid email addresses (missing @, domain, TLD, etc.)
// - Edge cases and malformed inputs
// - Uses external JSON test data for comprehensive coverage
package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
)

type ClientTestCase struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	ShouldFail  bool   `json:"should_fail"`
	Description string `json:"description"`
}

func TestClient_RequiresValidEmail(t *testing.T) {
	// Load test data
	testCases := loadClientTestCases(t)

	// Test each scenario
	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			// Attempt to create client with test data
			_, err := entity.NewClient(testCase.Name, testCase.Email, testCase.Phone, testCase.Address)
			
			if testCase.ShouldFail {
				// Should fail with validation error
				assert.Error(t, err, "Client creation should fail for: %s", testCase.Description)
			} else {
				// Should succeed
				assert.NoError(t, err, "Client creation should succeed for: %s", testCase.Description)
			}
		})
	}
}

func loadClientTestCases(t *testing.T) []ClientTestCase {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")
	
	// Build path to shared test data at tests root
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "testdata", "client", "client_test_cases.json")
	
	// Read test data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read test data file")
	
	// Parse JSON
	var testCases []ClientTestCase
	err = json.Unmarshal(data, &testCases)
	assert.NoError(t, err, "Failed to parse test data JSON")
	
	return testCases
}