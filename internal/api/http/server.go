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
	userService    *application.UserService
	clientHandler  *handlers.ClientHandler
	userHandler    *handlers.UserHandler
	healthHandler  *handlers.HealthHandler
	errorHandler   *middleware.ErrorHandler
}

// NewServer creates a new HTTP server with dependencies
func NewServer(billingService *application.BillingService, userService *application.UserService) *Server {
	return &Server{
		billingService: billingService,
		userService:    userService,
		clientHandler:  handlers.NewClientHandler(billingService),
		userHandler:    handlers.NewUserHandler(userService),
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
	
	// User API routes
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				s.userHandler.GetUserData(w, r)
			} else {
				s.userHandler.ListUsers(w, r)
			}
		case http.MethodPost:
			s.userHandler.CreateUser(w, r)
		case http.MethodPut:
			s.userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			s.userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

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