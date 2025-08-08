package storage

import (
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/storage"
	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestPostgreSQLStorage_ListAll_EmptyStorage(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	postgresStorage, ok := stack.Storage.(*storage.PostgreSQLStorage)
	assert.True(t, ok, "Expected PostgreSQL storage in integration test")

	// Act
	result, err := postgresStorage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestPostgreSQLStorage_ListAll_WithMultipleItems(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	postgresStorage, ok := stack.Storage.(*storage.PostgreSQLStorage)
	assert.True(t, ok, "Expected PostgreSQL storage in integration test")

	testData := map[string]interface{}{
		"client1": map[string]string{"name": "John Doe", "email": "john@example.com"},
		"client2": map[string]string{"name": "Jane Smith", "email": "jane@example.com"},
		"client3": map[string]string{"name": "Bob Wilson", "email": "bob@example.com"},
	}

	// Store test data
	for key, value := range testData {
		err := postgresStorage.Store(key, value)
		assert.NoError(t, err)
	}

	// Act
	result, err := postgresStorage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)

	// Verify all values are present (order may vary)
	for _, value := range result {
		found := false
		for _, expectedValue := range testData {
			// Compare the nested map structures
			if valueMap, ok := value.(map[string]interface{}); ok {
				if expectedMap, ok := expectedValue.(map[string]string); ok {
					if valueMap["name"] == expectedMap["name"] && valueMap["email"] == expectedMap["email"] {
						found = true
						break
					}
				}
			}
		}
		assert.True(t, found, "Value %v should match one of the expected test data entries", value)
	}
}

func TestPostgreSQLStorage_ListAll_SingleItem(t *testing.T) {
	// Arrange
	stack := testhelpers.NewCleanIntegrationTestStack()
	postgresStorage, ok := stack.Storage.(*storage.PostgreSQLStorage)
	assert.True(t, ok, "Expected PostgreSQL storage in integration test")

	expectedValue := map[string]string{"name": "Test User", "email": "test@example.com"}
	err := postgresStorage.Store("test_key", expectedValue)
	assert.NoError(t, err)

	// Act
	result, err := postgresStorage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	// Verify the returned value matches expected structure
	resultValue := result[0].(map[string]interface{})
	assert.Equal(t, expectedValue["name"], resultValue["name"])
	assert.Equal(t, expectedValue["email"], resultValue["email"])
}
