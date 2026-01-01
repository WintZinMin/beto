package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	app := NewApp()

	assert.NotNil(t, app)
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.Logger)
	assert.Nil(t, app.Server)
}

func TestHealthHandler(t *testing.T) {
	app := NewApp()
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.NotEmpty(t, response["timestamp"])
}

func TestVersionHandler(t *testing.T) {
	app := NewApp()
	req, err := http.NewRequest("GET", "/version", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, appName, response["name"])
	assert.Equal(t, version, response["version"])
}

func TestRootHandler(t *testing.T) {
	app := NewApp()
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], appName)
	assert.Equal(t, version, response["version"])
}

func TestStatusHandler(t *testing.T) {
	app := NewApp()
	req, err := http.NewRequest("GET", "/api/v1/status", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "v1", response["api"])
	assert.Equal(t, "running", response["status"])
	assert.NotEmpty(t, response["uptime"])
}

func TestLoggingMiddleware(t *testing.T) {
	app := NewApp()

	// Capture log output
	var buf bytes.Buffer
	app.Logger.SetOutput(&buf)

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that logging middleware logged the request
	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/health")
}

func TestCORSMiddleware(t *testing.T) {
	app := NewApp()

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/health", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkHeaders {
				assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
				assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "GET")
				assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
			}
		})
	}
}

func TestAppStartAndShutdown(t *testing.T) {
	app := NewApp()

	// Start server in background
	go func() {
		err := app.Start("0") // Use port 0 to get any available port
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestEnvironmentVariables(t *testing.T) {
	// Test with custom port
	os.Setenv("PORT", "9999")
	defer os.Unsetenv("PORT")

	port := os.Getenv("PORT")
	assert.Equal(t, "9999", port)

	// Test default port when not set
	os.Unsetenv("PORT")
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	assert.Equal(t, defaultPort, port)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "8080", defaultPort)
	assert.Equal(t, "Beto Application", appName)
	assert.Equal(t, "1.0.0", version)
}

func TestInvalidRoutes(t *testing.T) {
	app := NewApp()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Invalid path",
			method:         "GET",
			path:           "/invalid",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid method on valid path",
			method:         "POST",
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func BenchmarkHealthHandler(b *testing.B) {
	app := NewApp()
	req, _ := http.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
	}
}

func BenchmarkVersionHandler(b *testing.B) {
	app := NewApp()
	req, _ := http.NewRequest("GET", "/version", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
	}
}

func TestJSONResponseFormat(t *testing.T) {
	app := NewApp()

	endpoints := []string{"/health", "/version", "/", "/api/v1/status"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, err := http.NewRequest("GET", endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			// Verify it's valid JSON
			var jsonData interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &jsonData)
			assert.NoError(t, err, "Response should be valid JSON")
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	app := NewApp()

	const numRequests = 100
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/health", nil)
			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}
}
