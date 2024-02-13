package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func prepareServer (app *chi.Mux) *chi.Mux {
	chiServer := chi.NewRouter()

	// basic middleware setup
	chiServer.Use(middleware.RequestID)
	chiServer.Use(middleware.RealIP)
	chiServer.Use(middleware.Logger)
	chiServer.Use(middleware.Recoverer)

	// // Set a 60 sec timeout value on api request life
	chiServer.Use(middleware.Timeout(60 * time.Second))

	// register mux
	chiServer.Mount("/", app)

	return chiServer
}

func BuildAndStartServer(app *chi.Mux) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	env := LoadENVVariables()
	server := prepareServer(app)

	// start the server
	log.Println("\n Starting server on port", env)

	http.ListenAndServe(env.Port(), server)

	<-stopChan
	fmt.Println("\n Shutting down")
	time.Sleep(1 * time.Second)
}