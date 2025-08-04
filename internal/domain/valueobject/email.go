package valueobject

import (
	"encoding/json"
	"strings"
	
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
)

// Email represents a validated email address value object
type Email struct {
	value string `json:"value"`
}

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (Email, error) {
	// Normalize the email
	normalized := strings.ToLower(strings.TrimSpace(email))
	
	// Validate email format
	if normalized == "" {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationRequired, "email is required")
	}
	
	if len(normalized) > 254 {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationLength, "email too long (max 254 characters)")
	}
	
	// Check for @ symbol
	if !strings.Contains(normalized, "@") {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationFormat, "email must contain @ symbol")
	}
	
	// Split and validate parts
	parts := strings.Split(normalized, "@")
	if len(parts) != 2 {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationFormat, "email must have exactly one @ symbol")
	}
	
	localPart := parts[0]
	domain := parts[1]
	
	if localPart == "" {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationFormat, "email missing local part")
	}
	
	if domain == "" {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationFormat, "email missing domain")
	}
	
	// Check for TLD (basic check)
	if !strings.Contains(domain, ".") {
		return Email{}, errors.NewValidationError("email", email, errors.ValidationFormat, "email domain missing TLD")
	}
	
	return Email{value: normalized}, nil
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// Value returns the underlying email value
func (e Email) Value() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty checks if the email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email (before @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// MarshalJSON implements custom JSON marshaling for Email
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value string `json:"value"`
	}{
		Value: e.value,
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Email
func (e *Email) UnmarshalJSON(data []byte) error {
	var temp struct {
		Value string `json:"value"`
	}
	
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	
	e.value = temp.Value
	return nil
}