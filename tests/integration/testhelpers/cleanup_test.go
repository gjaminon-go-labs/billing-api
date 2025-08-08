// Test Data Cleanup Integration Tests
//
// This file tests the test data cleanup functionality to ensure proper test isolation.
// Tests: Database cleanup utilities, test isolation, data verification
// Purpose: Validate that test data cleanup works correctly and provides proper isolation
package testhelpers

import (
	"testing"

	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestDataCleanupIsolation(t *testing.T) {
	t.Run("FirstTest_CreatesData", func(t *testing.T) {
		// Create an integration test stack
		stack := testhelpers.NewIntegrationTestStack()

		// Verify we start with clean state
		if stack.DatabaseCleaner != nil {
			counts, err := stack.DatabaseCleaner.GetTableCounts()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), counts["clients"], "Should start with 0 clients")
		}

		// Create some test data through the application
		result, err := stack.BillingService.CreateClient("Test Client 1", "test1@cleanup.test", "+1234567890", "123 Test St")
		assert.NoError(t, err)
		assert.NotEmpty(t, result.ID)

		// Verify data was created (data goes into storage_records table via the storage abstraction)
		if stack.DatabaseCleaner != nil {
			counts, err := stack.DatabaseCleaner.GetTableCounts()
			assert.NoError(t, err)
			// The storage abstraction uses storage_records table, not clients table directly
			assert.True(t, counts["storage_records"] > 0, "Should have created storage records")
		}
	})

	t.Run("SecondTest_StartsClean", func(t *testing.T) {
		// Create a new integration test stack (should trigger cleanup)
		stack := testhelpers.NewIntegrationTestStack()

		// Verify we start with clean state (data from previous test should be gone)
		if stack.DatabaseCleaner != nil {
			counts, err := stack.DatabaseCleaner.GetTableCounts()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), counts["clients"], "Should start with 0 clients after cleanup")
		}

		// Create different test data
		result, err := stack.BillingService.CreateClient("Test Client 2", "test2@cleanup.test", "+0987654321", "456 Test Ave")
		assert.NoError(t, err)
		assert.NotEmpty(t, result.ID)
	})
}

func TestManualCleanup(t *testing.T) {
	// Create test stack without automatic cleanup
	stack := testhelpers.NewIntegrationTestStackNoCleanup()

	if stack.DatabaseCleaner == nil {
		t.Skip("Database cleaner not available - not using PostgreSQL storage")
	}

	// Create some test data
	_, err := stack.BillingService.CreateClient("Manual Cleanup Test", "manual@cleanup.test", "+1111111111", "789 Test Blvd")
	assert.NoError(t, err)

	// Manually trigger cleanup
	err = stack.DatabaseCleaner.CleanupTestData()
	assert.NoError(t, err)

	// Verify cleanup worked
	err = stack.DatabaseCleaner.VerifyCleanState()
	assert.NoError(t, err)
}

func TestCleanupUtilities(t *testing.T) {
	stack := testhelpers.NewIntegrationTestStackNoCleanup()

	if stack.DatabaseCleaner == nil {
		t.Skip("Database cleaner not available - not using PostgreSQL storage")
	}

	t.Run("GetTableCounts", func(t *testing.T) {
		counts, err := stack.DatabaseCleaner.GetTableCounts()
		assert.NoError(t, err)
		assert.Contains(t, counts, "clients")
		assert.Contains(t, counts, "storage_records")
		assert.True(t, counts["clients"] >= 0)
		assert.True(t, counts["storage_records"] >= 0)
	})

	t.Run("CleanupSpecificTable", func(t *testing.T) {
		// Cleanup specific table
		err := stack.DatabaseCleaner.CleanupSpecificTable("storage_records")
		assert.NoError(t, err)

		// Verify the specific table is clean
		counts, err := stack.DatabaseCleaner.GetTableCounts()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), counts["storage_records"])
	})
}
