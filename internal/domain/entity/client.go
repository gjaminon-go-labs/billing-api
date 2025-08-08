package entity

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/valueobject"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Client represents a billing client aggregate root
type Client struct {
	id        string `validate:"required,min=2,max=100"`
	name      string `validate:"required,min=2,max=100"`
	email     valueobject.Email
	phone     valueobject.Phone
	address   string `validate:"omitempty,max=500"`
	createdAt time.Time
	updatedAt time.Time
}

// NewClient creates a new Client with validation
func NewClient(name, email, phone, address string) (*Client, error) {
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

	// Create client instance
	client := &Client{
		id:        uuid.New().String(),
		name:      normalizedName,
		email:     emailVO,
		phone:     phoneVO,
		address:   normalizedAddress,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}

	// Validate the complete client using hybrid approach
	if err := client.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// Validate performs hybrid validation: value objects + struct tags + custom business rules
func (c *Client) Validate() error {
	// 1. Value objects are already validated during creation
	// 2. Run declarative validation on primitive fields (struct tags)
	if err := validator.New().Struct(c); err != nil {
		return c.convertValidatorErrors(err)
	}

	// 3. Run any additional custom business validation
	return c.validateBusinessRules()
}

// convertValidatorErrors converts validator library errors to structured ValidationErrors
func (c *Client) convertValidatorErrors(err error) error {
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
func (c *Client) validateBusinessRules() error {
	// Future business rules can be added here:
	// - Email uniqueness (requires repository)
	// - Complex cross-field validation
	// - Context-specific rules
	// - Domain-specific constraints

	return nil
}

// NewClientWithID creates a client with a specific ID (for repository loading)
func NewClientWithID(id, name, email, phone, address string, createdAt, updatedAt time.Time) (*Client, error) {
	// Validate and create value objects
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}

	phoneVO, err := valueobject.NewPhone(phone)
	if err != nil {
		return nil, err // ValidationError already properly structured
	}

	// Create client instance
	client := &Client{
		id:        id,
		name:      strings.TrimSpace(name),
		email:     emailVO,
		phone:     phoneVO,
		address:   strings.TrimSpace(address),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	// Validate the complete client using hybrid approach
	if err := client.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// UpdateDetails updates client details with validation
func (c *Client) UpdateDetails(name, phone, address string) error {
	// Create new phone value object
	phoneVO, err := valueobject.NewPhone(phone)
	if err != nil {
		return err // ValidationError already properly structured
	}

	// Update fields (normalization + validation via struct tags)
	c.name = strings.TrimSpace(name)
	c.phone = phoneVO
	c.address = strings.TrimSpace(address)
	c.updatedAt = time.Now().UTC()

	// Validate the updated client using hybrid approach
	return c.Validate()
}

// UpdateEmail updates the client's email address
func (c *Client) UpdateEmail(email string) error {
	// Validate and create email value object
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return err // ValidationError already properly structured
	}

	c.email = emailVO
	c.updatedAt = time.Now().UTC()

	return nil
}

// Getters
func (c *Client) ID() string {
	return c.id
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Email() valueobject.Email {
	return c.email
}

func (c *Client) Phone() valueobject.Phone {
	return c.phone
}

func (c *Client) Address() string {
	return c.address
}

func (c *Client) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Client) UpdatedAt() time.Time {
	return c.updatedAt
}

// EmailString returns the email as string (convenience method for compatibility)
func (c *Client) EmailString() string {
	return c.email.String()
}

// PhoneString returns the phone as string (convenience method for compatibility)
func (c *Client) PhoneString() string {
	return c.phone.String()
}

// Equals checks if two clients are equal (by ID)
func (c *Client) Equals(other *Client) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}

// String returns a string representation of the client
func (c *Client) String() string {
	return "Client{ID: " + c.id + ", Name: " + c.name + ", Email: " + c.email.String() + "}"
}

// MarshalJSON implements custom JSON marshaling for Client
func (c *Client) MarshalJSON() ([]byte, error) {
	// Create a struct with public fields for JSON marshaling
	jsonClient := struct {
		ID        string            `json:"id"`
		Name      string            `json:"name"`
		Email     valueobject.Email `json:"email"`
		Phone     valueobject.Phone `json:"phone"`
		Address   string            `json:"address"`
		CreatedAt time.Time         `json:"createdAt"`
		UpdatedAt time.Time         `json:"updatedAt"`
	}{
		ID:        c.id,
		Name:      c.name,
		Email:     c.email,
		Phone:     c.phone,
		Address:   c.address,
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
	}

	return json.Marshal(jsonClient)
}

// UnmarshalJSON implements custom JSON unmarshaling for Client
func (c *Client) UnmarshalJSON(data []byte) error {
	// Create a struct with public fields for JSON unmarshaling
	var jsonClient struct {
		ID        string            `json:"id"`
		Name      string            `json:"name"`
		Email     valueobject.Email `json:"email"`
		Phone     valueobject.Phone `json:"phone"`
		Address   string            `json:"address"`
		CreatedAt time.Time         `json:"createdAt"`
		UpdatedAt time.Time         `json:"updatedAt"`
	}

	if err := json.Unmarshal(data, &jsonClient); err != nil {
		return err
	}

	// Assign to private fields
	c.id = jsonClient.ID
	c.name = jsonClient.Name
	c.email = jsonClient.Email
	c.phone = jsonClient.Phone
	c.address = jsonClient.Address
	c.createdAt = jsonClient.CreatedAt
	c.updatedAt = jsonClient.UpdatedAt

	return nil
}
