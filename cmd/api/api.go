package api

import (
	"log"
	"net/http"

	"github.com/clemilsonazevedo/blog/config/database"
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/internal/routes/private"
	"github.com/clemilsonazevedo/blog/internal/routes/public"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User = entities.User
type Post = entities.Post
type Comment = entities.Comment

func InitServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))

	db, err := database.NewPostgresConfig()
	if err != nil {
		log.Fatal("ERROR INITIALIZING DATABASE")
	}

	database.AutoMigrate(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postController := controller.NewPostController(postService)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentController := controller.NewCommentController(commentService)

	//api routes
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(middlewares.SetVersionHeader("v1.0"))

		v1.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"version":"v1.0","status":"ok"}`))
		})

		public.BindPublicRoutes(
			userController,
			postController,
			v1,
		)
		private.BindPrivateRoutes(
			postController,
			userController,
			commentController,
			userService,
			v1,
		)
	})

	return r
}
