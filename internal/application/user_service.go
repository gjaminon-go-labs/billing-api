package application

import (
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
)

// UserService handles user-related business operations
// Following the Application Service pattern from Clean Architecture
type UserService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new user service with dependencies
func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user in the system
func (s *UserService) CreateUser(name, email, phone, address string) (*entity.User, error) {
	// 1. Check if user with this email already exists
	existingUser, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err // Repository error
	}
	
	if existingUser != nil {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleDuplicate,
			"A user with this email already exists",
			"email",
		)
	}
	
	// 2. Create new user entity (includes domain validation)
	user, err := entity.NewUser(name, email, phone, address)
	if err != nil {
		return nil, err // Validation error
	}
	
	// 3. Store user in repository
	if err := s.userRepository.Create(user); err != nil {
		return nil, err // Repository error
	}
	
	return user, nil
}

// GetUserByID retrieves a user by their ID with authorization check
func (s *UserService) GetUserByID(userID, requestingUserID string) (*entity.User, error) {
	// Authorization: Users can only access their own data
	// In production, this would integrate with a proper authentication system
	if userID != requestingUserID {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleUnauthorized,
			"You can only access your own user data",
			"user_id",
		)
	}
	
	// Retrieve user from repository
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, err // Repository error
	}
	
	if user == nil {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleNotFound,
			"User not found",
			"user_id",
		)
	}
	
	return user, nil
}

// GetUserByEmail retrieves a user by their email (admin operation)
func (s *UserService) GetUserByEmail(email string) (*entity.User, error) {
	// This method would typically require admin privileges
	// For now, it's available but should be protected by middleware
	
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err // Repository error
	}
	
	if user == nil {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleNotFound,
			"User not found",
			"email",
		)
	}
	
	return user, nil
}

// UpdateUser updates user information with authorization check
func (s *UserService) UpdateUser(userID, requestingUserID, name, phone, address string) (*entity.User, error) {
	// Authorization: Users can only update their own data
	if userID != requestingUserID {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleUnauthorized,
			"You can only update your own user data",
			"user_id",
		)
	}
	
	// Retrieve existing user
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, err // Repository error
	}
	
	if user == nil {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleNotFound,
			"User not found",
			"user_id",
		)
	}
	
	// Update user details (includes domain validation)
	if err := user.UpdateDetails(name, phone, address); err != nil {
		return nil, err // Validation error
	}
	
	// Save updated user
	if err := s.userRepository.Update(user); err != nil {
		return nil, err // Repository error
	}
	
	return user, nil
}

// UpdateUserEmail updates user email with authorization check
func (s *UserService) UpdateUserEmail(userID, requestingUserID, newEmail string) (*entity.User, error) {
	// Authorization: Users can only update their own data
	if userID != requestingUserID {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleUnauthorized,
			"You can only update your own email",
			"user_id",
		)
	}
	
	// Check if new email is already in use
	existingUser, err := s.userRepository.GetByEmail(newEmail)
	if err != nil {
		return nil, err // Repository error
	}
	
	if existingUser != nil && existingUser.ID() != userID {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleDuplicate,
			"A user with this email already exists",
			"email",
		)
	}
	
	// Retrieve current user
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, err // Repository error
	}
	
	if user == nil {
		return nil, errors.NewBusinessRuleError(
			errors.BusinessRuleNotFound,
			"User not found",
			"user_id",
		)
	}
	
	// Update email (includes domain validation)
	if err := user.UpdateEmail(newEmail); err != nil {
		return nil, err // Validation error
	}
	
	// Save updated user
	if err := s.userRepository.Update(user); err != nil {
		return nil, err // Repository error
	}
	
	return user, nil
}

// DeleteUser removes a user from the system with authorization check
func (s *UserService) DeleteUser(userID, requestingUserID string) error {
	// Authorization: Users can only delete their own account
	// Admins could delete any account (would need admin check here)
	if userID != requestingUserID {
		return errors.NewBusinessRuleError(
			errors.BusinessRuleUnauthorized,
			"You can only delete your own account",
			"user_id",
		)
	}
	
	// Check if user exists before deletion
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return err // Repository error
	}
	
	if user == nil {
		return errors.NewBusinessRuleError(
			errors.BusinessRuleNotFound,
			"User not found",
			"user_id",
		)
	}
	
	// Delete user from repository
	return s.userRepository.Delete(userID)
}

// ListUsers retrieves a paginated list of users (admin operation)
func (s *UserService) ListUsers(limit, offset int) ([]*entity.User, int64, error) {
	// This method would typically require admin privileges
	// For now, it's available but should be protected by middleware
	
	// Get users with pagination
	users, err := s.userRepository.List(limit, offset)
	if err != nil {
		return nil, 0, err // Repository error
	}
	
	// Get total count for pagination metadata
	total, err := s.userRepository.Count()
	if err != nil {
		return nil, 0, err // Repository error
	}
	
	return users, total, nil
}