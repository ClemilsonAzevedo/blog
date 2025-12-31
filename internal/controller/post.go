package controller

import (
	"fmt"

	"github.com/clemilsonazevedo/blog/internal/services"
	"github.com/go-chi/chi/v5"
)

func FindPostBySlug(c chi.Mux) error {
	post := services.FindPostBySlug("clema")

	fmt.Printf("%v", post)
	return nil
}
