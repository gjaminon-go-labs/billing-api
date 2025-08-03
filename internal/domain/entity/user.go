package entity

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/valueobject"
)

// User represents a system user aggregate root
type User struct {
	id        string
	name      string `validate:"required,min=2,max=100"`
	email     valueobject.Email
	phone     valueobject.Phone
	address   string `validate:"omitempty,max=500"`
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User with validation
func NewUser(name, email, phone, address string) (*User, error) {
	// Validate and create value objects
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}
	
	phoneVO, err := valueobject.NewPhone(phone)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}
	
	// Normalize primitive fields (validation handled by struct tags)
	normalizedName := strings.TrimSpace(name)
	normalizedAddress := strings.TrimSpace(address)
	
	// Create user instance
	user := &User{
		id:        uuid.New().String(),
		name:      normalizedName,
		email:     emailVO,
		phone:     phoneVO,
		address:   normalizedAddress,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}
	
	// Validate the complete user using hybrid approach
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate performs hybrid validation: value objects + struct tags + custom business rules
func (u *User) Validate() error {
	// 1. Value objects are already validated during creation
	// 2. Run declarative validation on primitive fields (struct tags)
	if err := validator.New().Struct(u); err != nil {
		return u.convertValidatorErrors(err)
	}
	
	// 3. Run any additional custom business validation
	return u.validateBusinessRules()
}

// convertValidatorErrors converts validator library errors to structured ValidationErrors
func (u *User) convertValidatorErrors(err error) error {
	validationErrors := errors.NewValidationErrors()
	
	if validatorErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validatorErrs {
			field := strings.ToLower(fieldErr.Field())
			var code errors.ErrorCode
			var message string
			
			switch fieldErr.Tag() {
			case "required":
				code = errors.ValidationRequired
				message = field + " is required"
			case "min":
				code = errors.ValidationLength
				message = field + " must be at least " + fieldErr.Param() + " characters"
			case "max":
				code = errors.ValidationLength
				message = field + " must be at most " + fieldErr.Param() + " characters"
			default:
				code = errors.ValidationFormat
				message = field + " validation failed"
			}
			
			validationErrors.Add(field, fieldErr.Value(), code, message)
		}
	}
	
	if validationErrors.HasErrors() {
		return validationErrors
	}
	
	return err // Return original error if we couldn't convert it
}

// validateBusinessRules performs custom business validation beyond struct tags and value objects
func (u *User) validateBusinessRules() error {
	// Future business rules can be added here:
	// - Email uniqueness (requires repository)
	// - Complex cross-field validation
	// - Context-specific rules
	// - Domain-specific constraints
	
	return nil
}

// NewUserWithID creates a user with a specific ID (for repository loading)
func NewUserWithID(id, name, email, phone, address string, createdAt, updatedAt time.Time) (*User, error) {
	// Validate and create value objects
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}
	
	phoneVO, err := valueobject.NewPhone(phone)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}
	
	// Create user instance
	user := &User{
		id:        id,
		name:      strings.TrimSpace(name),
		email:     emailVO,
		phone:     phoneVO,
		address:   strings.TrimSpace(address),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	
	// Validate the complete user using hybrid approach
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	return user, nil
}

// UpdateDetails updates user details with validation
func (u *User) UpdateDetails(name, phone, address string) error {
	// Create new phone value object
	phoneVO, err := valueobject.NewPhone(phone)
	if err != nil {
		return err // ValidationError already properly structured
	}
	
	// Update fields (normalization + validation via struct tags)
	u.name = strings.TrimSpace(name)
	u.phone = phoneVO
	u.address = strings.TrimSpace(address)
	u.updatedAt = time.Now().UTC()
	
	// Validate the updated user using hybrid approach
	return u.Validate()
}

// UpdateEmail updates the user's email address
func (u *User) UpdateEmail(email string) error {
	// Validate and create email value object
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return err // ValidationError already properly structured
	}
	
	u.email = emailVO
	u.updatedAt = time.Now().UTC()
	
	return nil
}

// Getters
func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Phone() valueobject.Phone {
	return u.phone
}

func (u *User) Address() string {
	return u.address
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// EmailString returns the email as string (convenience method for compatibility)
func (u *User) EmailString() string {
	return u.email.String()
}

// PhoneString returns the phone as string (convenience method for compatibility)
func (u *User) PhoneString() string {
	return u.phone.String()
}

// Equals checks if two users are equal (by ID)
func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id == other.id
}

// String returns a string representation of the user
func (u *User) String() string {
	return "User{ID: " + u.id + ", Name: " + u.name + ", Email: " + u.email.String() + "}"
}