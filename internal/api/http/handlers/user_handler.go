package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *application.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserData handles GET /api/v1/users/{id} requests
func (h *UserHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL path (simple extraction for now)
	// In production, use a proper router like Gorilla Mux or Chi
	userID := r.URL.Query().Get("id")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "User ID is required", "id")
		return
	}

	// Basic UUID validation
	if len(userID) != 36 {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FORMAT", "User ID must be a valid UUID", "id")
		return
	}

	// TODO: Extract requesting user ID from authentication context
	// For now, we'll use the same ID (user can only access their own data)
	requestingUserID := userID // This would come from JWT token or session

	// Call application service with authorization check
	user, err := h.userService.GetUserByID(userID, requestingUserID)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO (filters sensitive data)
	response := h.toUserResponse(user)

	// Write success response
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// CreateUser handles POST /api/v1/users requests
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse request body
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone,omitempty"`
		Address string `json:"address,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", "")
		return
	}

	// Basic HTTP-level validation
	if req.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "name is required", "name")
		return
	}
	if req.Email == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "email is required", "email")
		return
	}

	// Call application service
	user, err := h.userService.CreateUser(req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO
	response := h.toUserResponse(user)

	// Write success response
	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// UpdateUser handles PUT /api/v1/users/{id} requests
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL
	userID := r.URL.Query().Get("id")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "User ID is required", "id")
		return
	}

	// Basic UUID validation
	if len(userID) != 36 {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FORMAT", "User ID must be a valid UUID", "id")
		return
	}

	// Parse request body
	var req struct {
		Name    string `json:"name"`
		Phone   string `json:"phone,omitempty"`
		Address string `json:"address,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", "")
		return
	}

	// Basic validation
	if req.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "name is required", "name")
		return
	}

	// TODO: Extract requesting user ID from authentication context
	requestingUserID := userID // This would come from JWT token or session

	// Call application service with authorization check
	user, err := h.userService.UpdateUser(userID, requestingUserID, req.Name, req.Phone, req.Address)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entity to response DTO
	response := h.toUserResponse(user)

	// Write success response
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// DeleteUser handles DELETE /api/v1/users/{id} requests
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
	if r.Method != http.MethodDelete {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL
	userID := r.URL.Query().Get("id")
	if userID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_REQUIRED", "User ID is required", "id")
		return
	}

	// Basic UUID validation
	if len(userID) != 36 {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FORMAT", "User ID must be a valid UUID", "id")
		return
	}

	// TODO: Extract requesting user ID from authentication context
	requestingUserID := userID // This would come from JWT token or session

	// Call application service with authorization check
	err := h.userService.DeleteUser(userID, requestingUserID)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Write success response (no content)
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers handles GET /api/v1/users requests (admin operation)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// TODO: Add admin authorization check here
	// For now, this endpoint is available but should be protected

	// Call application service
	users, total, err := h.userService.ListUsers(limit, offset)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	// Convert domain entities to response DTOs
	userResponses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.toUserResponse(user)
	}

	// Create paginated response
	response := struct {
		Users  []dtos.UserResponse `json:"users"`
		Total  int64               `json:"total"`
		Limit  int                 `json:"limit"`
		Offset int                 `json:"offset"`
	}{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	// Write success response
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// handleDomainError converts domain errors to appropriate HTTP responses
func (h *UserHandler) handleDomainError(w http.ResponseWriter, err error) {
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
		
		// Map specific business rule errors to appropriate HTTP status codes
		switch errors.GetErrorCode(err) {
		case errors.BusinessRuleUnauthorized:
			h.writeErrorResponse(w, http.StatusForbidden, code, message, "")
		case errors.BusinessRuleNotFound:
			h.writeErrorResponse(w, http.StatusNotFound, code, message, "")
		case errors.BusinessRuleDuplicate:
			h.writeErrorResponse(w, http.StatusConflict, code, message, "")
		default:
			h.writeErrorResponse(w, http.StatusUnprocessableEntity, code, message, "")
		}
		return
	}

	if errors.IsRepositoryError(err) {
		code := string(errors.GetErrorCode(err))
		message := errors.GetUserMessage(err)
		h.writeErrorResponse(w, http.StatusInternalServerError, code, message, "")
		return
	}

	// Fallback for unknown errors
	h.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred", "")
}

// toUserResponse converts a domain User entity to HTTP response DTO
// This method filters out sensitive data and only returns safe information
func (h *UserHandler) toUserResponse(user *entity.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.EmailString(),
		Phone:     user.PhoneString(),
		Address:   user.Address(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// writeSuccessResponse writes a successful JSON response
func (h *UserHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := dtos.SuccessResponse{
		Data:    data,
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error JSON response
func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, code, message, field string) {
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