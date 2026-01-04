package public

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindPublicRoutes(uc *controller.UserController, pc *controller.PostController, c chi.Router) {
	c.Group(func(r chi.Router) {
		r.Post("/register", uc.CreateUser)
		r.Post("/login", uc.LoginUser)

		r.Get("/posts", pc.GetAllPosts)
		r.Get("/post/{id}", pc.GetPostById)
	})
}
