package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage_ListAll_EmptyStorage(t *testing.T) {
	// Arrange
	storage := NewInMemoryStorage()

	// Act
	result, err := storage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestInMemoryStorage_ListAll_WithMultipleItems(t *testing.T) {
	// Arrange
	storage := NewInMemoryStorage()
	
	testData := map[string]interface{}{
		"key1": "value1",
		"key2": "value2", 
		"key3": "value3",
	}
	
	// Store test data
	for key, value := range testData {
		err := storage.Store(key, value)
		assert.NoError(t, err)
	}

	// Act
	result, err := storage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	
	// Verify all values are present (order may vary)
	expectedValues := []interface{}{"value1", "value2", "value3"}
	for _, value := range result {
		assert.Contains(t, expectedValues, value)
	}
}

func TestInMemoryStorage_ListAll_SingleItem(t *testing.T) {
	// Arrange
	storage := NewInMemoryStorage()
	expectedValue := "single_value"
	
	err := storage.Store("test_key", expectedValue)
	assert.NoError(t, err)

	// Act
	result, err := storage.ListAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedValue, result[0])
}