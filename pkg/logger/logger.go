package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Context keys for trace and request IDs
type ctxKey string

const (
	TraceIDKey   ctxKey = "trace_id"
	RequestIDKey ctxKey = "request_id"
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
	Info(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
	Fatal(ctx context.Context, msg string, args ...interface{})
	InfoJSON(ctx context.Context, payload map[string]interface{})
	Close() error
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Service   string                 `json:"service,omitempty"`
}

// Config holds logger configuration
type Config struct {
	LogFilePath   string
	ServiceName   string
	EnableConsole bool
	EnableFile    bool
	MinLevel      LogLevel
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
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	var writers []io.Writer
	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}
	if config.EnableFile {
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
		service:  config.ServiceName,
		minLevel: config.MinLevel,
	}, nil
}

// JSONLogger implements Logger with JSON formatting and file output
type JSONLogger struct {
	writers  []io.Writer
	service  string
	minLevel LogLevel
	mu       sync.RWMutex
}

// log writes a structured log entry, extracting context fields from ctx
func (l *JSONLogger) log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.shouldLog(level) {
		return
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	contextFields := make(map[string]interface{})
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		contextFields["trace_id"] = traceID
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		contextFields["request_id"] = reqID
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Context:   contextFields,
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
	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *JSONLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, LevelInfo, msg, args...)
}
func (l *JSONLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, LevelError, msg, args...)
}
func (l *JSONLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, LevelWarn, msg, args...)
}
func (l *JSONLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, LevelDebug, msg, args...)
}
func (l *JSONLogger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, LevelFatal, msg, args...)
}

func (l *JSONLogger) InfoJSON(ctx context.Context, payload map[string]interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.shouldLog(LevelInfo) {
		return
	}

	contextFields := make(map[string]interface{})
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		contextFields["trace_id"] = traceID
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		contextFields["request_id"] = reqID
	}
	// Merge payload into context
	for k, v := range payload {
		contextFields[k] = v
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     LevelInfo,
		Message:   "",
		Context:   contextFields,
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
