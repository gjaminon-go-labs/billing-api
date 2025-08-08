package application

import (
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestBillingService_ListClients_EmptyService(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Act
	clients, err := service.ListClients()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, clients)
	assert.Empty(t, clients)
}

func TestBillingService_ListClients_WithMultipleClients(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Create clients via service
	_, err := service.CreateClient("John Doe", "john@example.com", "+1234567890", "123 Main St")
	assert.NoError(t, err)

	_, err = service.CreateClient("Jane Smith", "jane@example.com", "+0987654321", "456 Oak Ave")
	assert.NoError(t, err)

	_, err = service.CreateClient("Bob Wilson", "bob@example.com", "", "")
	assert.NoError(t, err)

	// Act
	clients, err := service.ListClients()

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

func TestBillingService_ListClients_SingleClient(t *testing.T) {
	// Arrange
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Create client via service
	client, err := service.CreateClient("Test User", "test@example.com", "+1111111111", "Test Address")
	assert.NoError(t, err)

	// Act
	clients, err := service.ListClients()

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
