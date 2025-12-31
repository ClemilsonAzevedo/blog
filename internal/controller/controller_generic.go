package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Controller[T types.Entity] struct {
	service *service.Service[T]
}

func NewController[T types.Entity](service *service.Service[T]) *Controller[T] {
	return &Controller[T]{
		service: service,
	}
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

func (c *Controller[T]) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parseID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := c.service.GetByID(parseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (c *Controller[T]) Update(w http.ResponseWriter, r *http.Request) {
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error decoding JSON:", err)
		return
	}
	if err := c.service.Update(entity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error decoding JSON:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
