package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/dtos"
)

// ErrorHandler provides middleware for handling panics and errors
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler middleware
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// RecoverMiddleware recovers from panics and returns a proper error response
func (e *ErrorHandler) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				
				// Write internal server error response
				e.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred")
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func (e *ErrorHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s - %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers for development
func (e *ErrorHandler) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// writeErrorResponse writes a structured error response
func (e *ErrorHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	errorDetail := dtos.ErrorDetail{
		Code:    code,
		Message: message,
	}

	response := dtos.ErrorResponse{
		Error:   errorDetail,
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}