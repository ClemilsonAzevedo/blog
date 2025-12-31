package controller

import (
	"encoding/json"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/types"
	"github.com/clemilsonazevedo/blog/internal/services"
)

type Controller[T types.Entity] struct {
	service *services.Service[T]
}

func (c *Controller[T]) Create(w http.ResponseWriter, r *http.Request) {
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.service.Create(entity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *Controller[T]) Update(w http.ResponseWriter, r *http.Request) {
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.service.Update(entity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
