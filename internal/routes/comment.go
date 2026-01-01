package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindCommentRoutes(cc *controller.CommentController, c *chi.Mux) {
	c.Group(func(r chi.Router) {
		r.Post("/comment", cc.CreateComment)
		r.Get("/comment/{id}", cc.GetCommentById)
		r.Put("/comment/{id}", cc.UpdateComment)
		r.Delete("/comment/{id}", cc.DeleteComment)
	})
}
