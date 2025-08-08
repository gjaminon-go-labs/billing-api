package dtos

// CreateClientRequest represents the HTTP request body for creating a client
type CreateClientRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

// UpdateClientRequest represents the HTTP request body for updating a client
// Note: Email is intentionally excluded for security/audit reasons
type UpdateClientRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}
