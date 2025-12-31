package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/go-chi/chi/v5"
)

func BindCommentRoutes(uc *controller.Controller[entities.Comment], c *chi.Mux) {
	c.Group(func(r chi.Router) {
		r.Post("/comment", uc.Create)
		r.Get("/comment/{id}", uc.GetByID)
	})
}
