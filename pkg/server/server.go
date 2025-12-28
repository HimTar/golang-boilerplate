package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/himtar/go-boilerplate/pkg/logger"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	// App is the main router/handler
	App *HTTPRouter

	// Port to listen on (e.g., ":8080")
	Port string

	// Middlewares to apply globally
	Middlewares []Middleware

	// Server timeouts
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration

	// Shutdown timeout - how long to wait for connections to close
	ShutdownTimeout time.Duration

	// Lifecycle hooks
	OnStartup  func() error // Called after server starts
	OnShutdown func() error // Called before server shuts down

	// Logger for structured logging (nil = use default log package)
	Logger logger.Logger
}

// DefaultServerConfig returns sensible production defaults
func DefaultServerConfig(
	app *HTTPRouter,
	port string,
	middlewares []Middleware,
	logger *logger.Logger,
	readTimeout, writeTimeout, shutdownTimeout, idleTimeout, readHeaderTimeout *time.Duration,
) (*ServerConfig, error) {

	if logger == nil {
		return nil, fmt.Errorf("ERROR: Logger is required")
	}

	const (
		defaultReadTimeout       = 15 * time.Second
		defaultWriteTimeout      = 15 * time.Second
		defaultShutdownTimeout   = 10 * time.Second
		defaultIdleTimeout       = 60 * time.Second
		defaultReadHeaderTimeout = 10 * time.Second
	)
	rt := defaultReadTimeout
	if readTimeout != nil {
		rt = *readTimeout
	}
	wt := defaultWriteTimeout
	if writeTimeout != nil {
		wt = *writeTimeout
	}
	st := defaultShutdownTimeout
	if shutdownTimeout != nil {
		st = *shutdownTimeout
	}
	it := defaultIdleTimeout
	if idleTimeout != nil {
		it = *idleTimeout
	}

	rht := defaultReadHeaderTimeout
	if readHeaderTimeout != nil {
		rht = *readHeaderTimeout
	}

	(*logger).Info("Initializing ServerConfig for port %s", port)
	return &ServerConfig{
		App:               app,
		Port:              port,
		Middlewares:       middlewares,
		ReadTimeout:       rt,
		WriteTimeout:      wt,
		ShutdownTimeout:   st,
		IdleTimeout:       it,
		ReadHeaderTimeout: rht,
		Logger:            *logger,
	}, nil
}

// Validate checks if configuration is valid
func (c *ServerConfig) Validate() error {
	c.Logger.Info("Validating server configuration...")
	if c.App == nil {
		c.Logger.Error("App cannot be nil")
		return errors.New("app cannot be nil")
	}
	if c.Port == "" {
		c.Logger.Error("Port cannot be empty")
		return errors.New("port cannot be empty")
	}
	if c.ShutdownTimeout <= 0 {
		c.Logger.Error("Shutdown timeout must be positive")
		return errors.New("shutdown timeout must be positive")
	}
	if c.Logger == nil {
		c.Logger.Error("Logger cannot be empty")
		return errors.New("Logger cannot be empty")
	}
	c.Logger.Info("Server configuration validated successfully")
	return nil
}

// prepareRouter applies middleware to the application router
func prepareRouter(app *HTTPRouter, middlewares []Middleware, logger logger.Logger) *HTTPRouter {
	server := New()

	// register middlewares
	if len(middlewares) > 0 {
		logger.Info("Preparing router with %d middleware(s)...", len(middlewares))
		server.Use(middlewares...)
	} else {
		logger.Info("Preparing router with no middleware")
	}

	// mounting app on the server
	server.Mount("/", app.Handler())

	return server
}

// BuildAndStartServer starts an HTTP server with graceful shutdown
// Returns error only if server fails to start or configuration is invalid
func BuildAndStartServer(config *ServerConfig) error {
	config.Logger.Info("Starting BuildAndStartServer...")
	// Validate configuration
	if err := config.Validate(); err != nil {
		config.Logger.Error("Invalid server configuration: %v", err)
		return fmt.Errorf("invalid server configuration: %w", err)
	}

	// Prepare router with middleware
	router := prepareRouter(config.App, config.Middlewares, config.Logger)

	// Create HTTP server with timeouts
	config.Logger.Info("Creating HTTP server on port %s", config.Port)
	srv := &http.Server{
		Addr:              config.Port,
		Handler:           router.Handler(),
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// Channel to capture server errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		config.Logger.Info("HTTP server starting on %s", config.Port)

		// Call startup hook if defined
		if config.OnStartup != nil {
			config.Logger.Info("Executing startup hook...")
			if err := config.OnStartup(); err != nil {
				config.Logger.Error("Startup hook failed: %v", err)
				serverErrors <- fmt.Errorf("startup hook failed: %w", err)
				return
			}
		}

		config.Logger.Info("Server is ready to accept connections on %s", config.Port)
		config.Logger.Info("ListenAndServe called")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			config.Logger.Error("Server failed: %v", err)
			serverErrors <- fmt.Errorf("server failed: %w", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	config.Logger.Info("Setting up signal handling for graceful shutdown...")
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		config.Logger.Error("Server error received: %v", err)
		return err
	case sig := <-stopChan:
		config.Logger.Info("Received signal: %v, initiating graceful shutdown...", sig)
	}

	// Create shutdown context with timeout
	config.Logger.Info("Creating shutdown context with timeout: %v", config.ShutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	// Call shutdown hook if defined
	if config.OnShutdown != nil {
		config.Logger.Info("Executing shutdown hooks...")
		if err := config.OnShutdown(); err != nil {
			config.Logger.Error("Shutdown hook failed: %v", err)
		}
	}

	// Attempt graceful shutdown
	config.Logger.Info("Shutting down server (timeout: %v)...", config.ShutdownTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		// Force close if graceful shutdown fails
		config.Logger.Error("Graceful shutdown failed, forcing close: %v", err)
		if closeErr := srv.Close(); closeErr != nil {
			config.Logger.Error("Failed to close server: %v", closeErr)
			return fmt.Errorf("failed to close server: %w", closeErr)
		}
		return fmt.Errorf("forced shutdown due to: %w", err)
	}

	config.Logger.Info("Server stopped gracefully")
	return nil
}
