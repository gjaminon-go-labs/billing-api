package dtos

import "time"

// ClientResponse represents the HTTP response body for a client
type ClientResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   ErrorDetail `json:"error"`
	Success bool        `json:"success"`
}

// ErrorDetail contains specific error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}