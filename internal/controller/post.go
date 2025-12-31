package controller

import (
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/services"
	"github.com/go-chi/chi/v5"
)

func FindPostBySlug(w http.ResponseWriter, r *http.Request) {
	slugParam := chi.URLParam(r, "slug")
	posts := services.FindPostBySlug(slugParam)

	w.Write([]byte(posts.Content))
}
