package middlewares

import (
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/service"
)

func RequireAuthorRole(us service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}

}
