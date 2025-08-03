// Production Main Entry Point
//
// This is the production entry point for the billing service.
// Provides: Kubernetes-ready server with graceful shutdown, signal handling, DI integration
// Features: Configuration loading, HTTP server lifecycle, graceful termination
// Deployment: Designed for Kubernetes with proper SIGTERM/SIGINT handling
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gjaminon-go-labs/billing-api/internal/config"
)

func main() {
	// Initialize application
	if err := run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

// run contains the main application logic
func run() error {
	log.Println("🚀 Starting Billing Service...")

	// 1. Load configuration
	environment := config.GetEnvironment()
	log.Printf("📋 Environment: %s", environment)

	appConfig, err := config.LoadConfig(environment)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	log.Printf("✅ Configuration loaded for %s environment", environment)

	// 2. Create DI container
	container, err := config.NewProductionContainerFromEnvironment(environment)
	if err != nil {
		return fmt.Errorf("failed to create DI container: %w", err)
	}
	log.Println("✅ Dependency injection container initialized")

	// 3. Get HTTP server from DI container
	httpServer, err := container.GetHTTPServer()
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}
	log.Println("✅ HTTP server created")

	// 4. Configure and start HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		Handler:      httpServer.Handler(),
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
		IdleTimeout:  appConfig.Server.IdleTimeout,
	}

	// 5. Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🌐 HTTP server starting on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// 6. Set up signal handling for Kubernetes
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, 
		syscall.SIGTERM, // Kubernetes graceful shutdown signal
		syscall.SIGINT,  // Ctrl+C for local development
	)

	// 7. Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		log.Println("✅ Server stopped")

	case sig := <-signals:
		log.Printf("🛑 Received signal: %s, starting graceful shutdown...", sig)

		// 8. Graceful shutdown sequence
		if err := gracefulShutdown(server, appConfig.Server.ShutdownTimeout); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
	}

	log.Println("✅ Billing Service stopped gracefully")
	return nil
}

// gracefulShutdown performs graceful shutdown of the HTTP server
func gracefulShutdown(server *http.Server, timeout time.Duration) error {
	log.Printf("⏳ Starting graceful shutdown (timeout: %s)...", timeout)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Phase 1: Stop accepting new requests (0-5 seconds)
	log.Println("📤 Stopping acceptance of new requests...")
	
	// Phase 2: Shutdown server with connection draining (5-25 seconds)
	log.Println("🔄 Draining existing connections...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Force closing server due to timeout: %v", err)
		
		// Phase 3: Force close if timeout exceeded (25-30 seconds)
		log.Println("🔨 Force closing remaining connections...")
		return server.Close()
	}

	log.Println("✅ All connections drained successfully")
	return nil
}

// Development notes:
// 
// This main.go is designed for Kubernetes deployment with:
// 
// 1. **SIGTERM Handling**: Kubernetes sends SIGTERM for graceful shutdown
// 2. **SIGINT Handling**: For local development (Ctrl+C)
// 3. **SIGKILL**: Cannot be caught - Kubernetes sends after grace period
// 
// 4. **Graceful Shutdown Timeline** (default 30s Kubernetes grace period):
//    - 0-5s:   Stop accepting new requests  
//    - 5-25s:  Drain existing connections
//    - 25-30s: Force close remaining connections
//    - 30s+:   Kubernetes sends SIGKILL (force termination)
//
// 5. **Configuration Sources** (priority order):
//    - Environment variables (Kubernetes secrets/configmaps)
//    - Environment-specific YAML (development.yaml, production.yaml)
//    - Base YAML (base.yaml)
//
// 6. **Health Checks**: 
//    - Readiness: /health (available after successful startup)
//    - Liveness: /health (available during normal operation)
//
// 7. **DI Container**: Uses optimized dependency injection with:
//    - Singleton services (memory efficient)
//    - Lazy initialization (performance optimized)
//    - Thread-safe operations (concurrent safe)
//
// 8. **Current Limitations**:
//    - Uses InMemoryStorage (PostgreSQL implementation pending)
//    - Basic health checks (database health check pending)
//    - No metrics endpoint (can be added via configuration)