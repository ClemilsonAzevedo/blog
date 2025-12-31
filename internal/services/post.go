package services

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/repository"
)

func FindPostBySlug(slug string) *entities.Post {
	post, err := repository.FindPostBySlug(slug)
	if err != nil {
		exceptions.BadRequestException(err.Error())
	}

	return post
}
