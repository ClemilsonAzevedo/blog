package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/clemilsonazevedo/blog/config/database"
	"github.com/clemilsonazevedo/blog/internal/cache"
	"github.com/clemilsonazevedo/blog/internal/controller"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/internal/http/middlewares"
	"github.com/clemilsonazevedo/blog/internal/http/routes/private"
	"github.com/clemilsonazevedo/blog/internal/http/routes/public"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
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

	database.MigrateRoleEnums(db)
	database.AutoMigrate(db)
	if err := CreateAuthor(db); err != nil {
		log.Printf("It is not possible to create the author: %v", err)
	} else {
		log.Println("Author created with success!")
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	postRepo := repository.NewPostRepository(db)
	postCache := cache.NewPostCache(5 * time.Minute)
	postService := service.NewPostService(postRepo, postCache)
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

func CreateAuthor(db *gorm.DB) error {
	authorName := os.Getenv("AUTHOR_NAME")
	authorEmail := os.Getenv("AUTHOR_EMAIL")
	authorPassword := os.Getenv("AUTHOR_PASSWORD")
	if authorEmail == "" || authorPassword == "" || authorName == "" {
		log.Println("UserName or Email or Password not configured!")
		return nil
	}

	var count int64
	db.Model(entities.User{}).Where("email = ?", authorEmail).Count(&count)

	if count > 0 {
		log.Println("Author already exists")
		return nil
	}
	//TODO: Usar a função de hash do Cle
	author := entities.User{
		UserName: authorName,
		Email:    authorEmail,
		Password: authorPassword,
		Role:     enums.Author,
	}

	if err := db.Create(&author).Error; err != nil {
		return err
	}

	log.Printf("Author created: %s", authorEmail)
	return nil
}
