package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BindPostRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		router.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello world"))
		})
	})

}
