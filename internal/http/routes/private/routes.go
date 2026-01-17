package private

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/http/middlewares"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
)

type UserController = controller.UserController
type PostController = controller.PostController
type CommentController = controller.CommentController

type UserService = service.UserService

func BindPrivateRoutes(
	pc *PostController,
	uc *UserController,
	cc *CommentController,
	us *UserService,
	c chi.Router,
) {
	c.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth(us))
		// Auth
		r.Post("/logout", uc.Logout)

		// Users
		r.Get("/profiles", uc.Profile)
		r.Put("/profiles", uc.UpdateUser)
		r.Delete("/profiles", uc.DeleteUser)

		// Comments
		r.Post("/comments", cc.CreateComment)
		r.Delete("/comments", cc.DeleteComment)

		// Author Role
		r.Group(func(a chi.Router) {
			a.Use(middlewares.RequireAuthorRole(us))
			a.Post("/posts", pc.CreatePost)
			a.Post("/posts/suggest", pc.CreatePostWithAi)
			a.Put("/posts", pc.UpdatePost)
			a.Delete("/posts", pc.DeletePost)
		})
	})
}
