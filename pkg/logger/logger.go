package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogFormat represents different logging formats
type LogFormat int

const (
	TextFormat LogFormat = iota
	JSONFormat
)

// Logger represents a structured logger
type Logger struct {
	level      LogLevel
	format     LogFormat
	output     io.Writer
	fields     map[string]interface{}
	callerSkip int
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Caller    string                 `json:"caller,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string
	Output     io.Writer
	CallerSkip int
}

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	logger := &Logger{
		level:      parseLogLevel(config.Level),
		format:     parseLogFormat(config.Format),
		output:     config.Output,
		fields:     make(map[string]interface{}),
		callerSkip: config.CallerSkip,
	}

	if logger.output == nil {
		logger.output = os.Stdout
	}

	return logger
}

// NewDefault creates a logger with default settings
func NewDefault() *Logger {
	return New(Config{
		Level:  "info",
		Format: "json",
		Output: os.Stdout,
	})
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.clone()
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := l.clone()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// WithContext extracts relevant information from context and adds it to logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	newLogger := l.clone()

	// Extract request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		newLogger.fields["request_id"] = requestID
	}

	// Extract user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		newLogger.fields["user_id"] = userID
	}

	return newLogger
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info logs an info level message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Fatal logs a fatal level message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

// log is the internal logging function
func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	// Check if we should log this level
	if level < l.level {
		return
	}

	// Format message with args
	message := msg
	if len(args) > 0 {
		message = fmt.Sprintf(msg, args...)
	}

	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Fields:    l.fields,
	}

	// Add caller information
	if level >= ERROR || l.level == DEBUG {
		if caller := l.getCaller(); caller != "" {
			entry.Caller = caller
		}
	}

	// Output the log entry
	l.output.Write([]byte(l.formatEntry(entry) + "\n"))
}

// formatEntry formats the log entry based on the configured format
func (l *Logger) formatEntry(entry LogEntry) string {
	switch l.format {
	case JSONFormat:
		if data, err := json.Marshal(entry); err == nil {
			return string(data)
		}
		// Fallback to text format if JSON marshaling fails
		fallthrough
	case TextFormat:
		var parts []string
		parts = append(parts, entry.Timestamp)
		parts = append(parts, fmt.Sprintf("[%s]", entry.Level))
		if entry.Caller != "" {
			parts = append(parts, fmt.Sprintf("(%s)", entry.Caller))
		}
		parts = append(parts, entry.Message)

		// Add fields
		if len(entry.Fields) > 0 {
			var fieldParts []string
			for k, v := range entry.Fields {
				fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
			}
			parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, ", ")))
		}

		return strings.Join(parts, " ")
	default:
		return entry.Message
	}
}

// getCaller returns the caller information
func (l *Logger) getCaller() string {
	_, file, line, ok := runtime.Caller(3 + l.callerSkip)
	if !ok {
		return ""
	}

	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// clone creates a copy of the logger
func (l *Logger) clone() *Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}

	return &Logger{
		level:      l.level,
		format:     l.format,
		output:     l.output,
		fields:     newFields,
		callerSkip: l.callerSkip,
	}
}

// parseLogLevel parses a string log level into LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// parseLogFormat parses a string log format into LogFormat
func parseLogFormat(format string) LogFormat {
	switch strings.ToLower(format) {
	case "json":
		return JSONFormat
	case "text", "plain":
		return TextFormat
	default:
		return JSONFormat
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetFormat sets the logging format
func (l *Logger) SetFormat(format LogFormat) {
	l.format = format
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// Global logger instance
var defaultLogger = NewDefault()

// Global logging functions
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}

func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

func WithContext(ctx context.Context) *Logger {
	return defaultLogger.WithContext(ctx)
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	defaultLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return defaultLogger
}

// HTTPLogMiddleware creates a logging middleware for HTTP requests
func (l *Logger) HTTPLogMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(wrapped, r)

			l.WithFields(map[string]interface{}{
				"method":      r.Method,
				"url":         r.URL.String(),
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
				"status_code": wrapped.statusCode,
				"duration":    time.Since(start).String(),
			}).Info("HTTP request")
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Standard library logger adapter
func (l *Logger) StdLogger() *log.Logger {
	return log.New(&loggerWriter{l}, "", 0)
}

// loggerWriter implements io.Writer for standard library logger compatibility
type loggerWriter struct {
	logger *Logger
}

func (w *loggerWriter) Write(p []byte) (int, error) {
	w.logger.Info("%s", string(p))
	return len(p), nil
}
