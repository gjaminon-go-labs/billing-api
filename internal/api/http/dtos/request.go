package dtos

// CreateClientRequest represents the HTTP request body for creating a client
type CreateClientRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

// GetUserRequest represents the request for getting user data
type GetUserRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}