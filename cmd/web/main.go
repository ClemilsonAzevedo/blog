package main

import (
	"log"
	"net/http"
	"os"

	"github.com/clemilsonazevedo/blog/config/db"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/internal/routes"
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
	//Router
	route := routes.InitRouter()

	//Database
	db, err := db.NewPostgresConfig()
	if err != nil {
		log.Println("ERROR INITIALIZING DATABASE")
		os.Exit(1)
	}

	repository.AutoMigrate(db)
	return route
}
