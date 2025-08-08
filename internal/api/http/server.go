package http

import (
	"net/http"
	"strings"

	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/middleware"
	"github.com/gjaminon-go-labs/billing-api/internal/application"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	billingService *application.BillingService
	clientHandler  *handlers.ClientHandler
	healthHandler  *handlers.HealthHandler
	errorHandler   *middleware.ErrorHandler
}

// NewServer creates a new HTTP server with dependencies
func NewServer(billingService *application.BillingService) *Server {
	return &Server{
		billingService: billingService,
		clientHandler:  handlers.NewClientHandler(billingService),
		healthHandler:  handlers.NewHealthHandler(),
		errorHandler:   middleware.NewErrorHandler(),
	}
}

// SetupRoutes configures HTTP routes and middleware
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler.Health)

	// API routes
	mux.HandleFunc("/api/v1/clients/", s.handleClientWithIDRoute) // Individual client operations
	mux.HandleFunc("/api/v1/clients", s.handleClientsRoute)       // Collection operations

	// Apply middleware chain
	handler := s.errorHandler.RecoverMiddleware(mux)
	handler = s.errorHandler.LoggingMiddleware(handler)
	handler = s.errorHandler.CORSMiddleware(handler)

	return handler
}

// handleClientsRoute routes requests to the appropriate client handler based on HTTP method
func (s *Server) handleClientsRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.clientHandler.CreateClient(w, r)
	case http.MethodGet:
		s.clientHandler.ListClients(w, r)
	default:
		// Return method not allowed for unsupported methods
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":{"code":"METHOD_NOT_ALLOWED","message":"Method not allowed"},"success":false}`))
	}
}

// handleClientWithIDRoute handles individual client operations (GET, PUT, DELETE /api/v1/clients/{id})
func (s *Server) handleClientWithIDRoute(w http.ResponseWriter, r *http.Request) {
	// Extract client ID from URL path
	clientID := extractClientIDFromPath(r.URL.Path)
	if clientID == "" {
		// Invalid path format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"INVALID_PATH","message":"Invalid client ID in path"},"success":false}`))
		return
	}

	// Route based on HTTP method
	switch r.Method {
	case http.MethodGet:
		s.clientHandler.GetClient(w, r, clientID)
	case http.MethodPut:
		s.clientHandler.UpdateClient(w, r, clientID)
	case http.MethodDelete:
		s.clientHandler.DeleteClient(w, r, clientID)
	default:
		// Return method not allowed for unsupported methods
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":{"code":"METHOD_NOT_ALLOWED","message":"Method not allowed"},"success":false}`))
	}
}

// extractClientIDFromPath extracts the client ID from URL path like /api/v1/clients/{id}
func extractClientIDFromPath(path string) string {
	// Expected path format: /api/v1/clients/{id}
	const prefix = "/api/v1/clients/"

	if !strings.HasPrefix(path, prefix) {
		return ""
	}

	// Extract the ID part after the prefix
	clientID := strings.TrimPrefix(path, prefix)

	// Remove any trailing slash or path segments
	if slashIndex := strings.Index(clientID, "/"); slashIndex != -1 {
		clientID = clientID[:slashIndex]
	}

	// Basic validation - not empty
	if strings.TrimSpace(clientID) == "" {
		return ""
	}

	return clientID
}

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.SetupRoutes()
}
