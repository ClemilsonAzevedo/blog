package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	//Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))

	//Routes
	BindPostRoutes(r)

	return r
}
