package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/modfile"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// Logger interface for dependency injection
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	// InfoJSON logs a structured JSON payload at info level
	InfoJSON(payload map[string]interface{})

	// WithContext returns a new logger instance with additional context fields
	WithContext(fields map[string]interface{}) Logger

	// WithTraceID returns a new logger instance with trace ID
	WithTraceID(traceID string) Logger

	// WithField returns a new logger instance with a single field
	WithField(key string, value interface{}) Logger

	// Close closes the logger and flushes any buffered data
	Close() error
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Service   string                 `json:"service,omitempty"`
}

// Config holds logger configuration
type Config struct {
	// LogFilePath is the path to the log file
	LogFilePath string

	// ServiceName identifies the service in logs
	ServiceName string

	// EnableConsole writes logs to stdout
	EnableConsole bool

	// EnableFile writes logs to file
	EnableFile bool

	// MinLevel is the minimum log level to output (default: INFO)
	MinLevel LogLevel
}

func getModuleNameFromGoMod() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(modFile.Module.Mod.Path)
}

// DefaultConfig returns default logger configuration
func DefaultConfig(serviceName string) *Config {

	return &Config{
		LogFilePath:   "logs/app.log",
		ServiceName:   serviceName,
		EnableConsole: true,
		EnableFile:    true,
		MinLevel:      LevelInfo,
	}
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if !c.EnableConsole && !c.EnableFile {
		return fmt.Errorf("at least one output (console or file) must be enabled")
	}
	if c.EnableFile && c.LogFilePath == "" {
		return fmt.Errorf("log file path required when file output is enabled")
	}
	return nil
}

// NewJSONLogger creates a new JSON logger instance
func NewJSONLogger(config *Config) (Logger, error) {
	if config == nil {
		return nil, fmt.Errorf("Config not found")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	var writers []io.Writer

	// Add console output if enabled
	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// Add file output if enabled
	if config.EnableFile {
		// Ensure log directory exists
		if err := os.MkdirAll("logs", 0755); err != nil {
			return nil, fmt.Errorf("failed to create logs directory: %w", err)
		}

		file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, file)
	}

	return &JSONLogger{
		writers:  writers,
		context:  make(map[string]interface{}),
		service:  config.ServiceName,
		minLevel: config.MinLevel,
	}, nil
}

// JSONLogger implements Logger with JSON formatting and file output
type JSONLogger struct {
	writers  []io.Writer
	context  map[string]interface{}
	traceID  string
	service  string
	minLevel LogLevel
	mu       sync.RWMutex
}

// log writes a structured log entry
func (l *JSONLogger) log(level LogLevel, msg string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Check if we should log this level
	if !l.shouldLog(level) {
		return
	}

	// Format message with args if provided
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Context:   l.copyContext(),
		TraceID:   l.traceID,
		Service:   l.service,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to plain text if JSON marshaling fails
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	// Add newline
	jsonBytes = append(jsonBytes, '\n')

	// Write to all configured writers
	for _, w := range l.writers {
		if _, err := w.Write(jsonBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		}
	}

	// Exit if fatal
	if level == LevelFatal {
		os.Exit(1)
	}
}

// shouldLog checks if a log level should be logged based on minimum level
func (l *JSONLogger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
		LevelFatal: 4,
	}
	return levels[level] >= levels[l.minLevel]
}

// copyContext creates a copy of the context map for thread safety
func (l *JSONLogger) copyContext() map[string]interface{} {
	if len(l.context) == 0 {
		return nil
	}

	copied := make(map[string]interface{}, len(l.context))
	for k, v := range l.context {
		copied[k] = v
	}
	return copied
}

// clone creates a new logger instance with copied configuration
func (l *JSONLogger) clone() *JSONLogger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return &JSONLogger{
		writers:  l.writers, // Writers are shared (safe for concurrent use)
		context:  l.copyContext(),
		traceID:  l.traceID,
		service:  l.service,
		minLevel: l.minLevel,
	}
}

// Info logs an informational message
func (l *JSONLogger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l *JSONLogger) InfoJSON(payload map[string]interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.shouldLog(LevelInfo) {
		return
	}

	// Merge payload into context for this log entry
	context := l.copyContext()
	if context == nil {
		context = make(map[string]interface{})
	}
	for k, v := range payload {
		context[k] = v
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     LevelInfo,
		Message:   "", // Message is optional for InfoJSON
		Context:   context,
		TraceID:   l.traceID,
		Service:   l.service,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}
	jsonBytes = append(jsonBytes, '\n')

	for _, w := range l.writers {
		if _, err := w.Write(jsonBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		}
	}
}

// Error logs an error message
func (l *JSONLogger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}

// Warn logs a warning message
func (l *JSONLogger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, msg, args...)
}

// Debug logs a debug message
func (l *JSONLogger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, msg, args...)
}

// Fatal logs a fatal message and exits
func (l *JSONLogger) Fatal(msg string, args ...interface{}) {
	l.log(LevelFatal, msg, args...)
}

// WithContext returns a new logger instance with additional context fields
// This creates a clone, so it's safe for concurrent use
func (l *JSONLogger) WithContext(fields map[string]interface{}) Logger {
	newLogger := l.clone()

	// Merge new fields into context
	for k, v := range fields {
		newLogger.context[k] = v
	}

	return newLogger
}

// WithTraceID returns a new logger instance with trace ID
func (l *JSONLogger) WithTraceID(traceID string) Logger {
	newLogger := l.clone()
	newLogger.traceID = traceID
	return newLogger
}

// WithField returns a new logger instance with a single field
func (l *JSONLogger) WithField(key string, value interface{}) Logger {
	newLogger := l.clone()
	newLogger.context[key] = value
	return newLogger
}

// Close closes any open file handles
func (l *JSONLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, w := range l.writers {
		if closer, ok := w.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return fmt.Errorf("failed to close writer: %w", err)
			}
		}
	}
	return nil
}

// NewDefaultLogger creates a default logger with fallback to simple logger
func NewDefaultLogger(serviceName string) (Logger, error) {
	jsonLogger, err := NewJSONLogger(DefaultConfig(serviceName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create JSON logger, using simple logger: %v\n", err)
		return nil, err
	}
	return jsonLogger, nil
}
