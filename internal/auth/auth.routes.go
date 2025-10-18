package auth

import (
	"github.com/go-chi/chi"
)

func Router() chi.Router {
	router := chi.NewRouter()

	router.Post("/login", loginHandler)
	router.Post("/register", registerHandler)

	return router
}
