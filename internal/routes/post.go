package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindPostRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Get("/posts/{slug}", controller.FindPostBySlug)
	})
}
