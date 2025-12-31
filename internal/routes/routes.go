package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	BindPostRoutes(r)

	return r
}
