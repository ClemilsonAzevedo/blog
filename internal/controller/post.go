package controller

import (
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/services"
	"github.com/go-chi/chi/v5"
)

func FindPostBySlug(w http.ResponseWriter, r *http.Request) {
	slugParam := chi.URLParam(r, "slug")
	posts := services.FindPostBySlug(slugParam)

	if posts == nil {
		w.WriteHeader(404)
		w.Write([]byte("article not found"))
		return
	}

	w.Write([]byte(posts.Content))
}
