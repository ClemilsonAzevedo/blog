package public

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindPublicRoutes(uc *controller.UserController, pc *controller.PostController, c chi.Router) {
	c.Group(func(r chi.Router) {
		r.Post("/login", uc.LoginUser)
		r.Post("/register", uc.CreateUser)
		r.Get("/posts", pc.GetAllPosts)
	})
}
