package dtos

import "fmt"

// PaginationRequest represents pagination parameters from the client
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Default pagination values
const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// SetDefaults sets default values for pagination if not provided
func (p *PaginationRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = DefaultPage
	}
	if p.Limit == 0 {
		p.Limit = DefaultLimit
	}
}

// Validate validates the pagination parameters
func (p *PaginationRequest) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if p.Limit < 1 || p.Limit > MaxLimit {
		return fmt.Errorf("limit must be between 1 and %d", MaxLimit)
	}
	return nil
}

// CalculateOffset calculates the database offset for pagination
func (p *PaginationRequest) CalculateOffset() int {
	return (p.Page - 1) * p.Limit
}

// PaginationResponse represents pagination metadata in the response
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// CalculateTotalPages calculates the total number of pages
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

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
	Success    bool                `json:"success"`
}
