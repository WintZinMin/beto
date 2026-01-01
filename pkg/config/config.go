package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the application
type Config struct {
	// Application settings
	Port        string
	AppName     string
	AppVersion  string
	Environment string

	// Database settings
	Database DatabaseConfig

	// Redis settings
	Redis RedisConfig

	// JWT settings
	JWT JWTConfig

	// Server settings
	Server ServerConfig

	// CORS settings
	CORS CORSConfig

	// Rate limiting
	RateLimit RateLimitConfig

	// Logging
	Logging LoggingConfig

	// External APIs
	ExternalAPIs ExternalAPIConfig

	// File upload
	FileUpload FileUploadConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	GracefulTimeout time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerWindow int
	WindowDuration    time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// ExternalAPIConfig holds external API configuration
type ExternalAPIConfig struct {
	APIKey             string
	ExternalServiceURL string
}

// FileUploadConfig holds file upload configuration
type FileUploadConfig struct {
	MaxFileSize string
	UploadPath  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		AppName:     getEnv("APP_NAME", "Beto Application"),
		AppVersion:  getEnv("APP_VERSION", "1.0.0"),
		Environment: getEnv("APP_ENV", "development"),

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "beto_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-change-me"),
			Expiry: getEnvAsDuration("JWT_EXPIRY", "24h"),
		},

		Server: ServerConfig{
			ReadTimeout:     getEnvAsDuration("READ_TIMEOUT", "15s"),
			WriteTimeout:    getEnvAsDuration("WRITE_TIMEOUT", "15s"),
			IdleTimeout:     getEnvAsDuration("IDLE_TIMEOUT", "60s"),
			GracefulTimeout: getEnvAsDuration("GRACEFUL_TIMEOUT", "30s"),
		},

		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},

		RateLimit: RateLimitConfig{
			RequestsPerWindow: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			WindowDuration:    getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),
		},

		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},

		ExternalAPIs: ExternalAPIConfig{
			APIKey:             getEnv("API_KEY", ""),
			ExternalServiceURL: getEnv("EXTERNAL_SERVICE_URL", "https://api.example.com"),
		},

		FileUpload: FileUploadConfig{
			MaxFileSize: getEnv("MAX_FILE_SIZE", "10MB"),
			UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		},
	}

	return config, nil
}

// DatabaseURL returns the database connection string
func (d DatabaseConfig) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// RedisURL returns the Redis connection string
func (r RedisConfig) RedisURL() string {
	if r.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d", r.Password, r.Host, r.Port, r.DB)
	}
	return fmt.Sprintf("redis://%s:%s/%d", r.Host, r.Port, r.DB)
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsTest returns true if running in test environment
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		var result []string
		for _, item := range splitAndTrim(value, ",") {
			if item != "" {
				result = append(result, item)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, item := range splitString(s, sep) {
		trimmed := trimSpace(item)
		result = append(result, trimmed)
	}
	return result
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading spaces
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Trim trailing spaces
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
