package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// Health handles GET /health requests
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:  "healthy",
		Service: "billing-service",
		Version: h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
