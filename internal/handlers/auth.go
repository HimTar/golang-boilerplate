package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login Route called")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Register Route called")
}

func AuthHandler() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", loginHandler)
	r.Get("/register", registerHandler)

	return r
}