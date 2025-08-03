package valueobject

import (
	"strings"
	
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
)

// Phone represents a validated phone number value object
type Phone struct {
	value string
}

// NewPhone creates a new Phone value object with validation
func NewPhone(phone string) (Phone, error) {
	// Normalize the phone number
	normalized := strings.TrimSpace(phone)
	
	// Empty phone is allowed (optional field)
	if normalized == "" {
		return Phone{value: ""}, nil
	}
	
	// Remove common formatting characters for validation
	cleanPhone := strings.ReplaceAll(normalized, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ".", "")
	
	// Length check
	if len(cleanPhone) < 7 || len(cleanPhone) > 15 {
		return Phone{}, errors.NewValidationError("phone", phone, errors.ValidationLength, "phone number must be 7-15 digits")
	}
	
	// Check if starts with valid digit (not 0 for international)
	if strings.HasPrefix(cleanPhone, "0") && strings.HasPrefix(cleanPhone, "+") {
		return Phone{}, errors.NewValidationError("phone", phone, errors.ValidationFormat, "international phone cannot start with 0")
	}
	
	return Phone{value: normalized}, nil
}

// String returns the string representation of the phone
func (p Phone) String() string {
	return p.value
}

// Value returns the underlying phone value
func (p Phone) Value() string {
	return p.value
}

// Equals checks if two phone numbers are equal
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// IsEmpty checks if the phone number is empty
func (p Phone) IsEmpty() bool {
	return p.value == ""
}

// HasCountryCode checks if the phone number has a country code
func (p Phone) HasCountryCode() bool {
	return strings.HasPrefix(p.value, "+")
}

// WithoutCountryCode returns the phone number without the country code
func (p Phone) WithoutCountryCode() string {
	if p.HasCountryCode() {
		return p.value[1:]
	}
	return p.value
}