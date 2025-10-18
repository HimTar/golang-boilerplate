package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/himtar/go-boilerplate/internal/auth"
	server "github.com/himtar/go-boilerplate/libraries"
)

func app() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Mount("/auth", auth.Router())

	return router
}

func main() {
	server.BuildAndStartServer(app())
}
