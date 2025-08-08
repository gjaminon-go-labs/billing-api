// Transaction Basic Test - Validates transaction isolation fundamentals
package integration_test

import (
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionBasic_SingleTransaction verifies basic transaction functionality
func TestTransactionBasic_SingleTransaction(t *testing.T) {
	// Use transaction isolation
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	
	// Create a client
	client, err := entity.NewClient(
		"Transaction Test Client",
		"transaction@example.com",
		"+15559990000",
		"123 Transaction St",
	)
	require.NoError(t, err, "Failed to create client entity")
	
	// Save the client
	err = stack.ClientRepo.Save(client)
	require.NoError(t, err, "Failed to save client")
	
	// Verify we can retrieve it in the same transaction
	retrieved, err := stack.ClientRepo.GetByID(client.ID())
	require.NoError(t, err, "Failed to retrieve client")
	assert.Equal(t, client.Name(), retrieved.Name())
	
	// After cleanup, the transaction will be rolled back
	// and the client won't exist in the database
}

// TestTransactionBasic_VerifyRollback ensures data doesn't persist after rollback
func TestTransactionBasic_VerifyRollback(t *testing.T) {
	// First, create a client in a transaction
	var clientID string
	t.Run("CreateInTransaction", func(t *testing.T) {
		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()
		
		client, err := entity.NewClient(
			"Rollback Test Client",
			"rollback@example.com",
			"+15558880000",
			"456 Rollback Ave",
		)
		require.NoError(t, err)
		
		err = stack.ClientRepo.Save(client)
		require.NoError(t, err)
		
		// Save the ID for later verification
		clientID = client.ID()
		
		// Client exists in this transaction
		_, err = stack.ClientRepo.GetByID(clientID)
		require.NoError(t, err, "Client should exist in transaction")
	})
	
	// Now verify it was rolled back
	t.Run("VerifyRollback", func(t *testing.T) {
		stack, cleanup := testhelpers.WithTransaction(t)
		defer cleanup()
		
		// Try to get the client - should not exist
		_, err := stack.ClientRepo.GetByID(clientID)
		assert.Error(t, err, "Client should not exist after rollback")
		
		// Count all clients - should not include the rolled back one
		allClients, err := stack.ClientRepo.GetAll()
		require.NoError(t, err)
		
		for _, c := range allClients {
			assert.NotEqual(t, "Rollback Test Client", c.Name(), 
				"Found rolled back client - transaction isolation failed")
			assert.NotEqual(t, "rollback@example.com", c.EmailString(),
				"Found rolled back client email - transaction isolation failed")
		}
	})
}

// TestTransactionBasic_MultipleOperations tests multiple operations in a transaction
func TestTransactionBasic_MultipleOperations(t *testing.T) {
	stack, cleanup := testhelpers.WithTransaction(t)
	defer cleanup()
	
	// Create multiple clients
	clients := []*entity.Client{}
	for i := 0; i < 3; i++ {
		client, err := entity.NewClient(
			"Multi Op Client",
			"multi@example.com",
			"+15557770000",
			"789 Multi St",
		)
		require.NoError(t, err)
		
		// Note: In a real scenario, we'd use unique emails
		// But in isolated transactions, duplicates are OK
		err = stack.ClientRepo.Save(client)
		require.NoError(t, err)
		
		clients = append(clients, client)
	}
	
	// Verify all were created
	for _, client := range clients {
		retrieved, err := stack.ClientRepo.GetByID(client.ID())
		require.NoError(t, err)
		assert.Equal(t, "Multi Op Client", retrieved.Name())
	}
	
	// Update one
	clients[0].UpdateDetails("Updated Multi Op", "+15551112222", "999 Updated St")
	err := stack.ClientRepo.Save(clients[0])
	require.NoError(t, err)
	
	// Delete another
	err = stack.ClientRepo.Delete(clients[1].ID())
	require.NoError(t, err)
	
	// Verify the changes
	updated, err := stack.ClientRepo.GetByID(clients[0].ID())
	require.NoError(t, err)
	assert.Equal(t, "Updated Multi Op", updated.Name())
	
	_, err = stack.ClientRepo.GetByID(clients[1].ID())
	assert.Error(t, err, "Deleted client should not be found")
	
	// All changes will be rolled back after cleanup
}