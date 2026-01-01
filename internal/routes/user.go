package routes

import (
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/go-chi/chi/v5"
)

func BindUserRoutes(uc *controller.UserController, r chi.Router) {
	r.Post("/user", uc.CreateUser)
	r.Get("/user/{id}", uc.GetUserById)
	r.Put("/user/{id}", uc.UpdateUser)
	r.Delete("/user/{id}", uc.DeleteUser)
	r.Get("/users", uc.GetAllUsers)
}
