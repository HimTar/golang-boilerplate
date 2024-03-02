package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/himtar/go-boilerplate/internal/handlers"
	"github.com/himtar/go-boilerplate/libraries/server"
)

func app() *chi.Mux {
	server := chi.NewRouter()

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	server.Mount("/auth", handlers.AuthHandler())

	return server
}

func main() {
	server.BuildAndStartServer(app())
}
