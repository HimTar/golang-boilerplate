package main

import (
	"net/http"

	"github.com/himtar/go-boilerplate/internal/auth"
	"github.com/himtar/go-boilerplate/internal/libraries/env"
	"github.com/himtar/go-boilerplate/pkg/logger"
	server "github.com/himtar/go-boilerplate/pkg/server"
)

type App struct {
	router *server.HTTPRouter
	Config *env.Variables
	Logger logger.Logger
	// Add other shared dependencies here
}

// NewApp constructs the App, sets up routes, and injects dependencies.
func NewApp(config *env.Variables, log logger.Logger) *App {
	r := server.New()

	// Application layer: register endpoints
	r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Context(), "Ping endpoint called")
		w.Write([]byte("pong"))
	})

	// Auth context (domain module)
	authModule := auth.NewAuthModule(config, log)
	r.Mount("/auth", authModule.Router())

	return &App{
		router: r,
		Config: config,
		Logger: log,
	}
}

// Router exposes the configured HTTP router for the server.
func (a *App) Router() *server.HTTPRouter {
	return a.router
}
