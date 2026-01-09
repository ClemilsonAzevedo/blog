package main

import (
	"net/http"
	"os"

	"github.com/clemilsonazevedo/blog/cmd/api"
)

func main() {
	srv := api.InitServer()
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		os.Exit(1)
	}
}
