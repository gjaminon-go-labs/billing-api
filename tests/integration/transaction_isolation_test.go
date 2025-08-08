// Transaction Isolation Proof of Concept Test
//
// This file demonstrates that transaction-based test isolation works correctly
// allowing parallel test execution without data conflicts.
package integration_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionIsolation_ParallelExecution verifies that tests can run in parallel
// without data conflicts when using transaction-based isolation
func TestTransactionIsolation_ParallelExecution(t *testing.T) {
	// Mark test as parallelizable
	t.Parallel()

	// Number of parallel tests to run
	const numParallel = 5

	// Run multiple tests in parallel
	for i := 0; i < numParallel; i++ {
		testNum := i
		t.Run(fmt.Sprintf("ParallelTest_%d", testNum), func(t *testing.T) {
			// Mark this subtest as parallel
			t.Parallel()

			// Use transaction isolation
			stack, cleanup := testhelpers.WithTransaction(t)
			defer cleanup()

			// Create a unique client in this transaction
			client, err := entity.NewClient(
				fmt.Sprintf("Test Client %d", testNum),
				fmt.Sprintf("test%d@example.com", testNum),
				fmt.Sprintf("+1555000%04d", testNum),
				"123 Test St",
			)
			require.NoError(t, err, "Failed to create client entity")

			// Save client
			err = stack.ClientRepo.Save(client)
			require.NoError(t, err, "Failed to save client in transaction")

			// Verify client exists in this transaction
			retrieved, err := stack.ClientRepo.GetByID(client.ID())
			require.NoError(t, err, "Failed to retrieve client in same transaction")
			assert.Equal(t, client.Name(), retrieved.Name())
			assert.Equal(t, client.EmailString(), retrieved.EmailString())

			// Count clients - in an isolated transaction, we should only see:
			// 1. Any data that was committed before our transaction started
			// 2. Our own uncommitted changes
			allClients, err := stack.ClientRepo.GetAll()
			require.NoError(t, err, "Failed to get all clients")

			// Check that our client is among the results
			found := false
			for _, c := range allClients {
				if c.ID() == client.ID() {
					found = true
					break
				}
			}
			assert.True(t, found, "Our client should be visible in our transaction")

			// The key test: other parallel tests' uncommitted data should NOT be visible
			// If we see emails from other parallel tests, isolation has failed
			for _, c := range allClients {
				// Check if this is a client from another parallel test (not ours)
				if c.ID() != client.ID() && c.EmailString() != client.EmailString() {
					// Check if it matches the pattern of our parallel tests
					for j := 0; j < numParallel; j++ {
						if j != testNum && c.EmailString() == fmt.Sprintf("test%d@example.com", j) {
							assert.Fail(t, "Transaction isolation failed",
								"Found uncommitted data from parallel test %d in test %d", j, testNum)
						}
					}
				}
			}
		})
	}
}

// TestTransactionIsolation_NoDataPersistence verifies that transaction rollback
// prevents any data from persisting after test completion
func TestTransactionIsolation_NoDataPersistence(t *testing.T) {
	// First test - create data in transaction
	t.Run("CreateDataInTransaction", func(t *testing.T) {
		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()

		// Create a client
		client, err := entity.NewClient(
			"Temporary Test Client",
			"temp@example.com",
			"+15550001234",
			"456 Temp Ave",
		)
		require.NoError(t, err, "Failed to create client entity")

		err = stack.ClientRepo.Save(client)
		require.NoError(t, err, "Failed to save client")

		// Verify it exists in this transaction
		retrieved, err := stack.ClientRepo.GetByID(client.ID())
		require.NoError(t, err, "Client should exist in transaction")
		assert.Equal(t, client.Name(), retrieved.Name())
	})

	// Second test - verify data was rolled back
	t.Run("VerifyDataRolledBack", func(t *testing.T) {
		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()

		// Get all clients - should be empty (or only contain data from this transaction)
		allClients, err := stack.ClientRepo.GetAll()
		require.NoError(t, err, "Failed to get all clients")

		// Should not see the "Temporary Test Client" from previous test
		for _, client := range allClients {
			assert.NotEqual(t, "Temporary Test Client", client.Name(),
				"Found data from previous test - transaction rollback failed")
			assert.NotEqual(t, "temp@example.com", client.EmailString(),
				"Found data from previous test - transaction rollback failed")
		}
	})
}

// TestTransactionIsolation_ConcurrentCreates verifies that concurrent creates
// in different transactions don't conflict
func TestTransactionIsolation_ConcurrentCreates(t *testing.T) {
	t.Parallel()

	// Use channels to coordinate concurrent operations
	start := make(chan struct{})
	done := make(chan error, 5)

	// Start 5 concurrent goroutines
	for i := 0; i < 5; i++ {
		go func(num int) {
			// Wait for start signal
			<-start

			// Each goroutine uses its own transaction
			stack, cleanup := testhelpers.WithTransaction(t)
			defer cleanup()

			// Create a client with the same email (would conflict without isolation)
			client, err := entity.NewClient(
				fmt.Sprintf("Concurrent Client %d", num),
				"concurrent@example.com", // Same email in all transactions
				fmt.Sprintf("+1555999%04d", num),
				"789 Concurrent Blvd",
			)
			if err != nil {
				done <- err
				return
			}

			// This should succeed in each isolated transaction
			err = stack.ClientRepo.Save(client)
			done <- err
		}(i)
	}

	// Start all goroutines simultaneously
	close(start)

	// Collect results
	for i := 0; i < 5; i++ {
		err := <-done
		assert.NoError(t, err, "Concurrent create failed - transactions not properly isolated")
	}
}

// TestTransactionIsolation_UpdateIsolation verifies that updates in one transaction
// don't affect reads in another transaction
func TestTransactionIsolation_UpdateIsolation(t *testing.T) {
	t.Parallel()

	// First, create a client that will exist before our test transactions
	baseStack := testhelpers.NewIntegrationTestStack()
	defer func() {
		if baseStack.DatabaseCleaner != nil {
			baseStack.DatabaseCleaner.CleanupTestData()
		}
	}()

	// Create a base client
	baseClient, err := entity.NewClient(
		"Original Name",
		"original@example.com",
		"+15550000000",
		"111 Original St",
	)
	require.NoError(t, err, "Failed to create base client entity")

	err = baseStack.ClientRepo.Save(baseClient)
	require.NoError(t, err, "Failed to save base client")

	// Now test isolation between transactions
	var wg sync.WaitGroup
	wg.Add(2)

	// Transaction 1: Update the client
	go func() {
		defer wg.Done()

		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()

		// Get and update the client
		client, err := stack.ClientRepo.GetByID(baseClient.ID())
		if err != nil {
			t.Logf("Transaction 1: Failed to get client: %v", err)
			return
		}

		// Update the client details
		err = client.UpdateDetails("Updated Name in Transaction 1", "+15551111111", "222 Updated Ave")
		if err != nil {
			t.Logf("Transaction 1: Failed to update client details: %v", err)
			return
		}

		err = stack.ClientRepo.Save(client)
		if err != nil {
			t.Logf("Transaction 1: Failed to save updated client: %v", err)
			return
		}

		// Verify update in this transaction
		updated, err := stack.ClientRepo.GetByID(baseClient.ID())
		if err != nil {
			t.Logf("Transaction 1: Failed to get updated client: %v", err)
			return
		}

		assert.Equal(t, "Updated Name in Transaction 1", updated.Name())
	}()

	// Transaction 2: Read the client (should see original, not update)
	go func() {
		defer wg.Done()

		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()

		// Get the client - should see original name
		client, err := stack.ClientRepo.GetByID(baseClient.ID())
		if err != nil {
			t.Logf("Transaction 2: Failed to get client: %v", err)
			return
		}

		// Should still see original name (not the update from Transaction 1)
		assert.Equal(t, "Original Name", client.Name(),
			"Transaction 2 saw updates from Transaction 1 - isolation failed")
	}()

	wg.Wait()
}
