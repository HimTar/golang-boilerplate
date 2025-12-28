package main

import (
	"context"
	"log"
	"os"

	"github.com/himtar/go-boilerplate/internal/libraries/env"
	"github.com/himtar/go-boilerplate/pkg/logger"
	server "github.com/himtar/go-boilerplate/pkg/server"
)

func main() {
	// Create a single context for server lifecycle logs
	serverContext := context.TODO()

	// Load Config variables
	configVariables, err := env.LoadENVVariables()
	if err != nil {
		log.Printf("ERROR: Failed to load env variables: %v", err)
		log.Println("ERROR: Server will not start due to env variables load failure.")
		os.Exit(1)
	}
	log.Println("ENV variables loaded successfully.")

	// Load logger
	loggerInstance, err := logger.NewDefaultLogger(configVariables.MODULE_NAME())
	if err != nil {
		log.Printf("ERROR: Failed to initialize logger: %v", err)
		log.Println("ERROR: Server will not start due to logger initialization failure.")
		os.Exit(1)
	}
	loggerInstance.Info(serverContext, "Logger initialized successfully.")

	// Create server configuration with defaults
	serverConfig := []server.Middleware{
		server.TraceIDMiddleware(),
		server.RequestIDMiddleware(),
		server.RealIPMiddleware(),
		server.LoggerMiddleware(loggerInstance),
		server.RecovererMiddleware(),
	}

	app := NewApp(configVariables, loggerInstance)

	config, err := server.DefaultServerConfig(
		serverContext,
		app.Router(),
		configVariables.Port(),
		serverConfig,
		&loggerInstance,
		configVariables.READ_TIMEOUT_MS(),
		configVariables.WRITE_TIMEOUT_MS(),
		configVariables.SHUTDOWN_TIMEOUT_MS(),
		configVariables.IDLE_TIMEOUT_MS(),
		configVariables.READ_HEADER_TIMEOUT_MS(),
	)
	if err != nil {
		loggerInstance.Error(serverContext, "ERROR: Failed to load server configurations: %v", err)
		loggerInstance.Error(serverContext, "ERROR: Server will not start due to server configuration failure.")
		os.Exit(1)
	}

	loggerInstance.Info(serverContext, "Server configuration created")

	// Add lifecycle hooks (optional)
	config.OnStartup = func() error {
		loggerInstance.Info(serverContext, "Server startup: Running initialization tasks...")
		// Initialize database connections, caches, etc.
		return nil
	}

	config.OnShutdown = func() error {
		loggerInstance.Info(serverContext, "Server shutdown: Cleaning up resources...")
		// Close database connections, flush logs, etc.
		return nil
	}

	loggerInstance.Info(serverContext, "Starting server with graceful shutdown...")
	if err := server.BuildAndStartServer(serverContext, config); err != nil {
		loggerInstance.Error(serverContext, "Server failed: %v", err)
	}
}
