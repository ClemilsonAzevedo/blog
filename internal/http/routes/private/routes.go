package private

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/http/middlewares"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
)

func BindPrivateRoutes(
	pc *controller.PostController,
	uc *controller.UserController,
	cc *controller.CommentController,
	us *service.UserService,
	c chi.Router,
) {
	c.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth(us))
		// Users
		r.Get("/logout", uc.Logout)

		r.Get("/profile", uc.GetUserById)
		r.Put("/profile", uc.UpdateUser)
		r.Delete("/profile", uc.DeleteUser)

		// Comments
		r.Post("/comment", cc.CreateComment)
		r.Get("/comments", cc.GetCommentById)

		// Author Role
		r.Group(func(a chi.Router) {
			a.Use(middlewares.RequireAuthorRole(us))
			a.Post("/post", pc.CreatePost)
			a.Put("/post/{id}", pc.UpdatePost)
			a.Delete("/post/{id}", pc.DeletePost)

			a.Get("/users", uc.GetAllUsers)
		})
	})
}
