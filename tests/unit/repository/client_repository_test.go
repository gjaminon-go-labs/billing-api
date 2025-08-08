package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/stretchr/testify/assert"
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

	// Load test fixtures
	fixtures := loadRepositoryTestFixtures(t)

	// Create and save test clients from fixtures
	client1, err := entity.NewClient(fixtures[0].Name, fixtures[0].Email, fixtures[0].Phone, fixtures[0].Address)
	assert.NoError(t, err)
	err = repo.Save(client1)
	assert.NoError(t, err)

	client2, err := entity.NewClient(fixtures[1].Name, fixtures[1].Email, fixtures[1].Phone, fixtures[1].Address)
	assert.NoError(t, err)
	err = repo.Save(client2)
	assert.NoError(t, err)

	client3, err := entity.NewClient(fixtures[2].Name, fixtures[2].Email, fixtures[2].Phone, fixtures[2].Address)
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
	expectedEmails := []string{fixtures[0].Email, fixtures[1].Email, fixtures[2].Email}
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

	// Load test fixture
	fixture := loadSingleClientFixture(t)

	// Create and save test client
	client, err := entity.NewClient(fixture.Name, fixture.Email, fixture.Phone, fixture.Address)
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

type ClientFixture struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func loadRepositoryTestFixtures(t *testing.T) []ClientFixture {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to fixture data
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "repository_test_fixtures.json")

	// Read fixture data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read fixture data file")

	// Parse JSON
	var fixtures []ClientFixture
	err = json.Unmarshal(data, &fixtures)
	assert.NoError(t, err, "Failed to parse fixture data JSON")

	return fixtures
}

func loadSingleClientFixture(t *testing.T) ClientFixture {
	// Get current file directory
	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "Failed to get current file path")

	// Build path to fixture data
	testDataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "client", "single_client_fixture.json")

	// Read fixture data file
	data, err := os.ReadFile(testDataPath)
	assert.NoError(t, err, "Failed to read fixture data file")

	// Parse JSON
	var fixture ClientFixture
	err = json.Unmarshal(data, &fixture)
	assert.NoError(t, err, "Failed to parse fixture data JSON")

	return fixture
}
