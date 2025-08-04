package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
)

func TestClientRepository_GetAll_EmptyRepository(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	repo := repository.NewClientRepository(storage)

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Empty(t, clients)
}

func TestClientRepository_GetAll_WithMultipleClients(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	repo := repository.NewClientRepository(storage)
	
	// Create and save test clients
	client1, err := entity.NewClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)
	err = repo.Save(client1)
	assert.NoError(t, err)
	
	client2, err := entity.NewClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
	assert.NoError(t, err)
	err = repo.Save(client2)
	assert.NoError(t, err)
	
	client3, err := entity.NewClient("Bob Wilson", "bob@example.com", "", "")
	assert.NoError(t, err)
	err = repo.Save(client3)
	assert.NoError(t, err)

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 3)
	
	// Verify all clients are present (order may vary)
	expectedEmails := []string{"john@example.com", "jane@example.com", "bob@example.com"}
	actualEmails := make([]string, len(clients))
	for i, client := range clients {
		actualEmails[i] = client.EmailString()
	}
	
	for _, expectedEmail := range expectedEmails {
		assert.Contains(t, actualEmails, expectedEmail)
	}
}

func TestClientRepository_GetAll_SingleClient(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	repo := repository.NewClientRepository(storage)
	
	// Create and save test client
	client, err := entity.NewClient("Test User", "test@example.com", "+1111111111", "Test Address")
	assert.NoError(t, err)
	err = repo.Save(client)
	assert.NoError(t, err)

	// Act
	clients, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 1)
	
	// Verify client data
	retrievedClient := clients[0]
	assert.Equal(t, client.ID(), retrievedClient.ID())
	assert.Equal(t, client.Name(), retrievedClient.Name())
	assert.Equal(t, client.EmailString(), retrievedClient.EmailString())
	assert.Equal(t, client.PhoneString(), retrievedClient.PhoneString())
	assert.Equal(t, client.Address(), retrievedClient.Address())
}