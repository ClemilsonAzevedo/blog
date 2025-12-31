package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/go-chi/chi/v5"
)

func BindUserRoutes(uc *controller.Controller[entities.User], c *chi.Mux) {
	c.Group(func(r chi.Router) {
		r.Post("/user", uc.Create)
		r.Get("/user/{id}", uc.GetByID)
	})
}
