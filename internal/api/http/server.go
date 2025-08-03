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
	mux.HandleFunc("/api/v1/clients", s.clientHandler.CreateClient)

	// Apply middleware chain
	handler := s.errorHandler.RecoverMiddleware(mux)
	handler = s.errorHandler.LoggingMiddleware(handler)
	handler = s.errorHandler.CORSMiddleware(handler)

	return handler
}

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.SetupRoutes()
}