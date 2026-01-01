package routes

import (
    "github.com/clemilsonazevedo/blog/internal/controller"
    "github.com/go-chi/chi/v5"
)

func BindUserRoutes(uc *controller.UserController, r chi.Router) {
    r.Route("/user", func(r chi.Router) {
        r.Post("/", uc.CreateUser)          // POST /user
        r.Get("/{id}", uc.GetUserById)      // GET /user/{id}
        r.Put("/{id}", uc.UpdateUser)       // PUT /user/{id}
        r.Delete("/{id}", uc.DeleteUser)    // DELETE /user/{id}
    })
    r.Get("/users", uc.GetAllUsers)         // GET /users
}