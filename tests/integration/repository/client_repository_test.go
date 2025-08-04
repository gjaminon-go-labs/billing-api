package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
)

func TestClientRepository_GetAll_IntegrationTest(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	repo := stack.ClientRepo
	
	// Create and save test clients
	client1, err := entity.NewClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)
	err = repo.Save(client1)
	assert.NoError(t, err)
	
	client2, err := entity.NewClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
	assert.NoError(t, err)
	err = repo.Save(client2)
	assert.NoError(t, err)

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 2)
	
	// Verify clients are present (order may vary)
	expectedEmails := []string{"john@example.com", "jane@example.com"}
	actualEmails := make([]string, len(clients))
	for i, client := range clients {
		actualEmails[i] = client.EmailString()
	}
	
	for _, expectedEmail := range expectedEmails {
		assert.Contains(t, actualEmails, expectedEmail)
	}
}

func TestClientRepository_GetAll_EmptyRepository_IntegrationTest(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	repo := stack.ClientRepo

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Empty(t, clients)
}