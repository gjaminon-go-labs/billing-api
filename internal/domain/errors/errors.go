package errors

import (
	"fmt"
	"strings"
)

// ErrorCode represents a specific error code for programmatic handling
type ErrorCode string

const (
	// Validation error codes
	ValidationRequired    ErrorCode = "VALIDATION_REQUIRED"
	ValidationFormat      ErrorCode = "VALIDATION_FORMAT"
	ValidationLength      ErrorCode = "VALIDATION_LENGTH"
	ValidationRange       ErrorCode = "VALIDATION_RANGE"
	
	// Business rule error codes
	BusinessRuleViolation ErrorCode = "BUSINESS_RULE_VIOLATION"
	BusinessRuleConflict  ErrorCode = "BUSINESS_RULE_CONFLICT"
	
	// Repository error codes
	RepositoryNotFound    ErrorCode = "REPOSITORY_NOT_FOUND"
	RepositoryConnection  ErrorCode = "REPOSITORY_CONNECTION"
	RepositoryConstraint  ErrorCode = "REPOSITORY_CONSTRAINT"
	RepositoryInternal    ErrorCode = "REPOSITORY_INTERNAL"
)

// ValidationError represents input validation failures
type ValidationError struct {
	Field   string
	Value   interface{}
	Code    ErrorCode
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation failed: %s", e.Message)
}

func (e ValidationError) ErrorCode() ErrorCode {
	return e.Code
}

func (e ValidationError) UserMessage() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, code ErrorCode, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Code:    code,
		Message: message,
	}
}

// BusinessRuleError represents domain business rule violations
type BusinessRuleError struct {
	Rule    string
	Code    ErrorCode
	Message string
	Context map[string]interface{}
}

func (e BusinessRuleError) Error() string {
	if e.Rule != "" {
		return fmt.Sprintf("business rule violation '%s': %s", e.Rule, e.Message)
	}
	return fmt.Sprintf("business rule violation: %s", e.Message)
}

func (e BusinessRuleError) ErrorCode() ErrorCode {
	return e.Code
}

func (e BusinessRuleError) UserMessage() string {
	return e.Message
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(rule string, code ErrorCode, message string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:    rule,
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// RepositoryError represents data persistence issues
type RepositoryError struct {
	Operation string
	Code      ErrorCode
	Message   string
	Cause     error
}

func (e RepositoryError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("repository error during '%s': %s", e.Operation, e.Message)
	}
	return fmt.Sprintf("repository error: %s", e.Message)
}

func (e RepositoryError) ErrorCode() ErrorCode {
	return e.Code
}

func (e RepositoryError) UserMessage() string {
	// Don't expose internal details to users
	switch e.Code {
	case RepositoryNotFound:
		return "The requested resource was not found"
	case RepositoryConnection:
		return "Service temporarily unavailable"
	case RepositoryConstraint:
		return "The operation conflicts with existing data"
	default:
		return "An internal error occurred"
	}
}

func (e RepositoryError) Unwrap() error {
	return e.Cause
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(operation string, code ErrorCode, message string, cause error) *RepositoryError {
	return &RepositoryError{
		Operation: operation,
		Code:      code,
		Message:   message,
		Cause:     cause,
	}
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

func (e ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(messages, "; "))
}

func (e ValidationErrors) UserMessage() string {
	if len(e.Errors) == 0 {
		return "Please check your input"
	}
	
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.UserMessage())
	}
	return strings.Join(messages, "; ")
}

// Add adds a validation error to the collection
func (e *ValidationErrors) Add(field string, value interface{}, code ErrorCode, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Value:   value,
		Code:    code,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// NewValidationErrors creates a new validation errors collection
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}