package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindCommentRoutes(cc *controller.CommentController, c chi.Router) {
	c.Post("/comment", cc.CreateComment)
	c.Get("/comment/{id}", cc.GetCommentById)
	c.Put("/comment/{id}", cc.UpdateComment)
	c.Delete("/comment/{id}", cc.DeleteComment)
}
