package errors

import (
	"errors"
)

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsValidationErrors checks if an error is a ValidationErrors (multiple)
func IsValidationErrors(err error) bool {
	var validationErrs *ValidationErrors
	return errors.As(err, &validationErrs)
}

// IsBusinessRuleError checks if an error is a BusinessRuleError
func IsBusinessRuleError(err error) bool {
	var businessErr *BusinessRuleError
	return errors.As(err, &businessErr)
}

// IsRepositoryError checks if an error is a RepositoryError
func IsRepositoryError(err error) bool {
	var repoErr *RepositoryError
	return errors.As(err, &repoErr)
}

// GetErrorCode extracts the error code from structured errors
func GetErrorCode(err error) ErrorCode {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.ErrorCode()
	}

	var businessErr *BusinessRuleError
	if errors.As(err, &businessErr) {
		return businessErr.ErrorCode()
	}

	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return repoErr.ErrorCode()
	}

	return ""
}

// GetUserMessage extracts a user-friendly message from structured errors
func GetUserMessage(err error) string {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.UserMessage()
	}

	var validationErrs *ValidationErrors
	if errors.As(err, &validationErrs) {
		return validationErrs.UserMessage()
	}

	var businessErr *BusinessRuleError
	if errors.As(err, &businessErr) {
		return businessErr.UserMessage()
	}

	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return repoErr.UserMessage()
	}

	// Fallback for unstructured errors
	return "An error occurred"
}

// IsClientError checks if the error is a client-side error (validation, business rules)
func IsClientError(err error) bool {
	return IsValidationError(err) || IsValidationErrors(err) || IsBusinessRuleError(err)
}

// IsServerError checks if the error is a server-side error (repository, infrastructure)
func IsServerError(err error) bool {
	return IsRepositoryError(err)
}
