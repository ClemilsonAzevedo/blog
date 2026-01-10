package api

import (
	"log"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/config/database"
	_ "github.com/clemilsonazevedo/blog/docs"
	"github.com/clemilsonazevedo/blog/internal/cache"
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/http/middlewares"
	"github.com/clemilsonazevedo/blog/internal/http/routes/private"
	"github.com/clemilsonazevedo/blog/internal/http/routes/public"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Blog API
// @version 1.0
// @description REST API for a blog application with posts, comments, and user management.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@blog.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name token

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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	postRepository := repository.NewPostRepository(db)
	postCache := cache.NewPostCache(5 * time.Minute)
	postService := service.NewPostService(postRepository, postCache)
	postController := controller.NewPostController(postService)

	commentRepository := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepository)
	commentController := controller.NewCommentController(commentService)

	// Swagger UI route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

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
			commentController,
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
