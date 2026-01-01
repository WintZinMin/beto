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

	"github.com/gorilla/mux"
)

const (
	defaultPort = "8080"
	appName     = "Beto Application"
	version     = "1.0.0"
)

// App represents the main application structure
type App struct {
	Router *mux.Router
	Server *http.Server
	Logger *log.Logger
}

// NewApp creates a new application instance
func NewApp() *App {
	app := &App{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stdout, "[BETO] ", log.LstdFlags|log.Lshortfile),
	}

	app.setupRoutes()
	return app
}

// setupRoutes configures all application routes
func (a *App) setupRoutes() {
	// Middleware (must be added before routes)
	a.Router.Use(a.corsMiddleware)
	a.Router.Use(a.loggingMiddleware)

	// Health check endpoint
	a.Router.HandleFunc("/health", a.healthHandler).Methods("GET", "OPTIONS")

	// Version endpoint
	a.Router.HandleFunc("/version", a.versionHandler).Methods("GET", "OPTIONS")

	// Root endpoint
	a.Router.HandleFunc("/", a.rootHandler).Methods("GET", "OPTIONS")

	// API routes
	api := a.Router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", a.statusHandler).Methods("GET", "OPTIONS")
}

// HTTP Handlers
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}

func (a *App) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"name": "%s", "version": "%s"}`, appName, version)
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Welcome to %s API", "version": "%s"}`, appName, version)
}

func (a *App) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"api": "v1", "status": "running", "uptime": "%s"}`, time.Since(startTime))
}

// Middleware
func (a *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		a.Logger.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func (a *App) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start initializes and starts the HTTP server
func (a *App) Start(port string) error {
	a.Server = &http.Server{
		Addr:         ":" + port,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.Logger.Printf("Starting %s on port %s", appName, port)
	return a.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Println("Shutting down server...")
	return a.Server.Shutdown(ctx)
}

var startTime = time.Now()

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create application instance
	app := NewApp()

	// Start server in a goroutine
	go func() {
		if err := app.Start(port); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	app.Logger.Println("Server exited")
}
