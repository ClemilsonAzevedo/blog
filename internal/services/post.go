package services

import (
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/repository"
)

func FindPostBySlug(slug string) response.PostResponse {
	post, err := repository.FindPostBySlug(slug)
	if err != nil {
		exceptions.BadRequestException(err.Error())
	}

	return post
}
