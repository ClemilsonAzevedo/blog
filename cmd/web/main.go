package main

import (
	"os"
	"net/http"
	"github.com/clemilsonazevedo/blog/internal/bootstrap"
)

func main() {
	srv := bootstrap.InitServer();
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		os.Exit(1)
	}
}