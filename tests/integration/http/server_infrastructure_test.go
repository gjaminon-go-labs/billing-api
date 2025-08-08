// HTTP Server Infrastructure Integration Tests
//
// This file contains integration tests for HTTP server infrastructure components.
// Tests: Health endpoints, CORS middleware, server configuration, cross-cutting concerns
// Scope: Integration tests - Complete HTTP Server infrastructure and middleware
// Use Cases: Infrastructure support for all use cases - Cross-cutting concerns
//
// Test Scenarios:
// - Health check endpoint functionality
// - CORS preflight request handling
// - HTTP middleware behavior
// - Server configuration and routing
// - Infrastructure endpoints and responses
//
// Components Tested:
// - HTTP Server routing and middleware
// - Health check handlers
// - CORS middleware configuration
// - Server infrastructure setup
// - Test helpers (NewInMemoryTestServer)
package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gjaminon-go-labs/billing-api/tests/testhelpers"
	"github.com/stretchr/testify/assert"
)

// BUSINESS_TITLE: System Health Monitoring
// BUSINESS_DESCRIPTION: Operations team and monitoring systems can check if the application is running properly and ready to serve customers
// USER_STORY: As a system administrator, I want to monitor system health so that I can ensure the application is available for users
// BUSINESS_VALUE: Enables proactive monitoring, prevents downtime, supports operational excellence and SLA compliance
// SCENARIOS_TESTED: Health endpoint availability, proper status responses, monitoring system integration
func TestHTTPServer_Integration_HealthCheck(t *testing.T) {
	// Set up complete HTTP server using InMemory test helpers
	server := testhelpers.NewInMemoryTestServer()

	// Create test server
	testServer := httptest.NewServer(server.Handler())
	defer testServer.Close()

	// Make health check request
	resp, err := http.Get(testServer.URL + "/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check content type
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Parse response
	var healthResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	assert.NoError(t, err)

	// Check health response structure
	assert.Equal(t, "healthy", healthResponse["status"])
	assert.Equal(t, "billing-service", healthResponse["service"])
	assert.Contains(t, healthResponse, "version")
}

// BUSINESS_TITLE: Cross-Domain API Access
// BUSINESS_DESCRIPTION: Web applications from different domains can securely access the API, enabling integrations and third-party applications
// USER_STORY: As a developer integrating with the API, I want to make requests from web applications without CORS errors
// BUSINESS_VALUE: Enables API integrations, supports third-party development, allows flexible frontend deployment strategies
// SCENARIOS_TESTED: CORS headers, preflight requests, cross-domain access security, API accessibility
func TestHTTPServer_Integration_CORS(t *testing.T) {
	// Set up complete HTTP server using InMemory test helpers
	server := testhelpers.NewInMemoryTestServer()

	// Create test server
	testServer := httptest.NewServer(server.Handler())
	defer testServer.Close()

	// Make OPTIONS request (preflight)
	req, err := http.NewRequest(http.MethodOptions, testServer.URL+"/api/v1/clients", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check CORS headers
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
