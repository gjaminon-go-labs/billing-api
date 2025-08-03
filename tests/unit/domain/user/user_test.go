package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
)

func TestNewUser_ValidData_Success(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "john.doe@example.com"
	phone := "+1234567890"
	address := "123 Main St"

	// Act
	user, err := entity.NewUser(name, email, phone, address)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID())
	assert.Equal(t, name, user.Name())
	assert.Equal(t, email, user.EmailString())
	assert.Equal(t, phone, user.PhoneString())
	assert.Equal(t, address, user.Address())
	assert.False(t, user.CreatedAt().IsZero())
	assert.False(t, user.UpdatedAt().IsZero())
}

func TestNewUser_EmptyName_ValidationError(t *testing.T) {
	// Arrange
	name := ""
	email := "john.doe@example.com"
	phone := "+1234567890"
	address := "123 Main St"

	// Act
	user, err := entity.NewUser(name, email, phone, address)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestNewUser_InvalidEmail_ValidationError(t *testing.T) {
	// Arrange
	name := "John Doe"
	email := "invalid-email"
	phone := "+1234567890"
	address := "123 Main St"

	// Act
	user, err := entity.NewUser(name, email, phone, address)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUpdateDetails_ValidData_Success(t *testing.T) {
	// Arrange
	user, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)
	
	newName := "Jane Doe"
	newPhone := "+0987654321"
	newAddress := "456 Oak Ave"

	// Act
	err = user.UpdateDetails(newName, newPhone, newAddress)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, newName, user.Name())
	assert.Equal(t, newPhone, user.PhoneString())
	assert.Equal(t, newAddress, user.Address())
}

func TestUpdateEmail_ValidEmail_Success(t *testing.T) {
	// Arrange
	user, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)
	
	newEmail := "jane.doe@example.com"

	// Act
	err = user.UpdateEmail(newEmail)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, newEmail, user.EmailString())
}

func TestUpdateEmail_InvalidEmail_ValidationError(t *testing.T) {
	// Arrange
	user, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)
	
	invalidEmail := "invalid-email"

	// Act
	err = user.UpdateEmail(invalidEmail)

	// Assert
	assert.Error(t, err)
	// Email should remain unchanged
	assert.Equal(t, "john.doe@example.com", user.EmailString())
}

func TestEquals_SameID_ReturnsTrue(t *testing.T) {
	// Arrange
	user1, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)
	
	user2, err := entity.NewUserWithID(
		user1.ID(),
		"Jane Smith",
		"jane.smith@example.com",
		"+0987654321",
		"456 Oak Ave",
		user1.CreatedAt(),
		user1.UpdatedAt(),
	)
	require.NoError(t, err)

	// Act & Assert
	assert.True(t, user1.Equals(user2))
	assert.True(t, user2.Equals(user1))
}

func TestEquals_DifferentID_ReturnsFalse(t *testing.T) {
	// Arrange
	user1, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)
	
	user2, err := entity.NewUser("Jane Smith", "jane.smith@example.com", "+0987654321", "456 Oak Ave")
	require.NoError(t, err)

	// Act & Assert
	assert.False(t, user1.Equals(user2))
	assert.False(t, user2.Equals(user1))
}

func TestEquals_NilUser_ReturnsFalse(t *testing.T) {
	// Arrange
	user, err := entity.NewUser("John Doe", "john.doe@example.com", "+1234567890", "123 Main St")
	require.NoError(t, err)

	// Act & Assert
	assert.False(t, user.Equals(nil))
}