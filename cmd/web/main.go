package main

import (
	"log"
	"net/http"
	"os"

	"github.com/clemilsonazevedo/blog/config/db_config"
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/internal/routes"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	srv := initSrv()
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		os.Exit(1)
	}
}

func initSrv() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))

	db, err := db_config.NewPostgresConfig()
	if err != nil {
		log.Println("ERROR INITIALIZING DATABASE")
		os.Exit(1)
	}

	db_config.AutoMigrate(db)

	postRepo := repository.NewRepository[entities.Post](db)
	postService := service.NewService[entities.Post](postRepo)
	postController := controller.NewController[entities.Post](postService)
	routes.BindPostRoutes(postController, r)

	userRepo := repository.NewRepository[entities.User](db)
	userService := service.NewService[entities.User](userRepo)
	userController := controller.NewController[entities.User](userService)
	routes.BindUserRoutes(userController, r)

	commentRepo := repository.NewRepository[entities.Comment](db)
	commentService := service.NewService[entities.Comment](commentRepo)
	commentController := controller.NewController[entities.Comment](commentService)
	routes.BindCommentRoutes(commentController, r)

	return r
}
