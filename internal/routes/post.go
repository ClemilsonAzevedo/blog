package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindPostRoutes(cc *controller.PostController, c chi.Router) {
	c.Group(func(r chi.Router) {
		r.Post("/post", cc.CreatePost)
		r.Get("/post/{slug}", cc.GetPostBySlug)
		r.Put("/post/{id}", cc.UpdatePost)
		r.Delete("/post/{id}", cc.DeletePost)
		r.Get("/posts", cc.GetPaginatedPosts)
	})
}
