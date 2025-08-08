package dtos

// PaginationRequest represents pagination parameters from query string
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationResponse represents pagination metadata in response
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Success    bool                `json:"success"`
}

// DefaultPaginationValues defines default pagination parameters
const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
	MinLimit     = 1
	MinPage      = 1
)

// Validate checks if pagination parameters are valid
func (p *PaginationRequest) Validate() error {
	if p.Page < MinPage {
		return NewValidationError("page", "page must be greater than 0")
	}
	if p.Limit < MinLimit || p.Limit > MaxLimit {
		return NewValidationError("limit", "limit must be between 1 and 100")
	}
	return nil
}

// SetDefaults applies default values if not set
func (p *PaginationRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = DefaultPage
	}
	if p.Limit == 0 {
		p.Limit = DefaultLimit
	}
}

// CalculateOffset returns the database offset for pagination
func (p *PaginationRequest) CalculateOffset() int {
	return (p.Page - 1) * p.Limit
}

// CalculateTotalPages calculates total pages from total count
func CalculateTotalPages(totalCount, limit int) int {
	if totalCount == 0 || limit == 0 {
		return 0
	}
	pages := totalCount / limit
	if totalCount%limit > 0 {
		pages++
	}
	return pages
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(field, message string) error {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}
