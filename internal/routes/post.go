package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindPostRoutes(cc *controller.PostController, c *chi.Mux) {
	c.Group(func(r chi.Router) {
		r.Post("/post", cc.CreatePost)
		r.Put("/post/{id}", cc.UpdatePost)
		r.Get("/post/{id}", cc.GetPostById)
		r.Delete("/post/{id}", cc.DeletePost)
		r.Get("/posts", cc.GetAllPosts)
	})
}
