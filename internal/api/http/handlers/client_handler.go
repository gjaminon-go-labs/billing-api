package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
)

// ClientHandler handles HTTP requests for client operations
type ClientHandler struct {
	billingService *application.BillingService
}

// NewClientHandler creates a new client handler
func NewClientHandler(billingService *application.BillingService) *ClientHandler {
	return &ClientHandler{
		billingService: billingService,
	}
}

// CreateClient handles POST /clients requests
func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse request body
	var req dtos.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", "")
		return
	}

	// Validate required fields (basic HTTP-level validation)
	if req.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "name is required", "name")
		return
	}
	if req.Email == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "email is required", "email")
		return
	}

	// Call application service
	client, err := h.billingService.CreateClientLegacy(req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO
	response := h.toClientResponse(client)

	// Write success response
	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// ListClients handles GET /clients requests
func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	
	// Always use pagination (with defaults if not specified)
	{
		// Parse and validate pagination
		paginationReq := dtos.PaginationRequest{}
		
		if pageStr != "" {
			page := 0
			_, err := fmt.Sscanf(pageStr, "%d", &page)
			if err != nil {
				h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "invalid page parameter", "")
				return
			}
			paginationReq.Page = page
		}
		
		if limitStr != "" {
			limit := 0
			_, err := fmt.Sscanf(limitStr, "%d", &limit)
			if err != nil {
				h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_PARAMETER", "invalid limit parameter", "")
				return
			}
			paginationReq.Limit = limit
		}
		
		// Validate before setting defaults (to catch invalid values like 0 or negative)
		if pageStr != "" && paginationReq.Page <= 0 {
			h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "page must be greater than 0", "")
			return
		}
		if limitStr != "" && (paginationReq.Limit <= 0 || paginationReq.Limit > dtos.MaxLimit) {
			h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "limit must be between 1 and 100", "")
			return
		}
		
		// Set defaults
		paginationReq.SetDefaults()
		
		// Final validation
		if err := paginationReq.Validate(); err != nil {
			h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), "")
			return
		}
		
		// Call paginated service method
		result, err := h.billingService.ListClientsWithPagination(paginationReq.Page, paginationReq.Limit)
		if err != nil {
			h.handleDomainError(w, err)
			return
		}
		
		// Convert domain entities to response DTOs
		clientResponses := make([]dtos.ClientResponse, len(result.Clients))
		for i, client := range result.Clients {
			clientResponses[i] = h.toClientResponse(client)
		}
		
		// Create paginated response
		paginationResponse := &dtos.PaginationResponse{
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			TotalCount: result.Pagination.TotalCount,
			TotalPages: result.Pagination.TotalPages,
		}
		
		// Write paginated response
		h.writePaginatedResponse(w, http.StatusOK, clientResponses, paginationResponse)
	}
}

// handleDomainError converts domain errors to appropriate HTTP responses
func (h *ClientHandler) handleDomainError(w http.ResponseWriter, err error) {
	// Check error type and map to HTTP status code
	if errors.IsValidationError(err) || errors.IsValidationErrors(err) {
		code := string(errors.GetErrorCode(err))
		message := errors.GetUserMessage(err)

		// Try to extract field information from validation error
		var field string
		if validationErr, ok := err.(*errors.ValidationError); ok {
			field = validationErr.Field
		}

		h.writeErrorResponse(w, http.StatusBadRequest, code, message, field)
		return
	}

	if errors.IsBusinessRuleError(err) {
		code := string(errors.GetErrorCode(err))
		message := errors.GetUserMessage(err)
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, code, message, "")
		return
	}

	if errors.IsRepositoryError(err) {
		code := errors.GetErrorCode(err)
		message := errors.GetUserMessage(err)

		// Map specific repository error codes to appropriate HTTP status codes
		var statusCode int
		switch code {
		case errors.RepositoryNotFound:
			statusCode = http.StatusNotFound
		case errors.RepositoryConstraint:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}

		h.writeErrorResponse(w, statusCode, string(code), message, "")
		return
	}

	// Fallback for unknown errors
	h.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred", "")
}

// toClientResponse converts a domain Client entity to HTTP response DTO
func (h *ClientHandler) toClientResponse(client *entity.Client) dtos.ClientResponse {
	return dtos.ClientResponse{
		ID:        client.ID(),
		Name:      client.Name(),
		Email:     client.EmailString(),
		Phone:     client.PhoneString(),
		Address:   client.Address(),
		CreatedAt: client.CreatedAt(),
		UpdatedAt: client.UpdatedAt(),
	}
}

// writeSuccessResponse writes a successful JSON response
func (h *ClientHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := dtos.SuccessResponse{
		Data:    data,
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error JSON response
func (h *ClientHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, code, message, field string) {
	errorDetail := dtos.ErrorDetail{
		Code:    code,
		Message: message,
	}
	if field != "" {
		errorDetail.Field = field
	}

	response := dtos.ErrorResponse{
		Error:   errorDetail,
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// GetClient handles GET /clients/{id} requests
func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request, clientID string) {
	// Get client from service
	client, err := h.billingService.GetClientByID(clientID)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO
	response := h.toClientResponse(client)

	// Write success response
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// UpdateClient handles PUT /clients/{id} requests
func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request, clientID string) {
	// Parse request body
	var req dtos.UpdateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", "")
		return
	}

	// Update client via service
	client, err := h.billingService.UpdateClient(clientID, req)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO
	response := h.toClientResponse(client)

	// Write success response
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// DeleteClient handles DELETE /clients/{id} requests
func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request, clientID string) {
	// Delete client via service
	err := h.billingService.DeleteClient(clientID)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Write success response with no content
	w.WriteHeader(http.StatusNoContent)
}

// writePaginatedResponse writes a paginated response with metadata
func (h *ClientHandler) writePaginatedResponse(w http.ResponseWriter, statusCode int, data interface{}, pagination *dtos.PaginationResponse) {
	response := dtos.PaginatedResponse{
		Data:       data,
		Pagination: pagination,
		Success:    true,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
