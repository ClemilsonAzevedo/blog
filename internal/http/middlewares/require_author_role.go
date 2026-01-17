package middlewares

import (
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/service"
)

func RequireAuthorRole(us *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*entities.User)
			if !ok {
				exceptions.Unauthorized(w, "You need login to this route")
				return
			}

			if user.Role != enums.Author {
				exceptions.Unauthorized(w, "Author Role Required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
