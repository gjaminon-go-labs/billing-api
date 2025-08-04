package http

import (
	"net/http"

	"github.com/gjaminon-go-labs/billing-api/internal/application"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/handlers"
	"github.com/gjaminon-go-labs/billing-api/internal/api/http/middleware"
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
	mux.HandleFunc("/api/v1/clients", s.handleClientsRoute)

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

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.SetupRoutes()
}