package public

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

type UserController = controller.UserController
type PostController = controller.PostController
type CommentController = controller.CommentController

func BindPublicRoutes(uc *UserController, pc *PostController, cc *CommentController,
	c chi.Router) {
	c.Group(func(r chi.Router) {
		// Auth
		r.Post("/register", uc.CreateUser)
		r.Post("/login", uc.LoginUser)

		// Posts
		r.Get("/posts", pc.GetPaginatedPosts)
		r.Get("/post", pc.GetPostById)

		// Comments
		r.Get("/comments", cc.GetCommentsByPostID)
	})
}
