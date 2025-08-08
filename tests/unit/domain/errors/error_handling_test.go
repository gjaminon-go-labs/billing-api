// Domain Error Handling Unit Tests
//
// This file contains unit tests for domain error types and error handling.
// Tests: ValidationError structure, error interface implementation, error formatting
// Scope: Pure unit tests - single component (domain errors) with no external dependencies
// Use Cases: All use cases - Error handling is cross-cutting concern
//
// Test Scenarios:
// - ValidationError creation and Error() method
// - Error message formatting and structure
// - Error interface compliance
// - Successful entity creation (no errors)
package errors

import (
	"testing"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/infrastructure/repository"
	"github.com/gjaminon-go-labs/billing-api/tests/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandling_ValidationError(t *testing.T) {
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Act - create client with invalid email
	client, err := service.CreateClient("John Doe", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, client)

	// Check if it's a ValidationError
	assert.True(t, errors.IsValidationError(err))

	// Check error code
	assert.Equal(t, errors.ValidationRequired, errors.GetErrorCode(err))

	// Check user message
	userMsg := errors.GetUserMessage(err)
	assert.Contains(t, userMsg, "email is required")

	// Check if it's a client error (not server error)
	assert.True(t, errors.IsClientError(err))
	assert.False(t, errors.IsServerError(err))
}

func TestErrorHandling_FormatValidationError(t *testing.T) {
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Act - create client with invalid email format
	client, err := service.CreateClient("John Doe", "invalid-email", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, client)

	// Check if it's a ValidationError
	assert.True(t, errors.IsValidationError(err))

	// Check error code
	assert.Equal(t, errors.ValidationFormat, errors.GetErrorCode(err))

	// Check user message
	userMsg := errors.GetUserMessage(err)
	assert.Contains(t, userMsg, "email must contain @ symbol")
}

func TestErrorHandling_SuccessfulCreation(t *testing.T) {
	// Set up dependencies
	storage := infrastructure.NewInMemoryStorage()
	clientRepo := repository.NewClientRepository(storage)
	service := application.NewBillingService(clientRepo)

	// Act - create valid client
	client, err := service.CreateClient("John Doe", "john@example.com", "+1234567890", "123 Main St")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name())
	assert.Equal(t, "john@example.com", client.EmailString())
}
