package middlewares

import (
	"context"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/http/auth"
	"github.com/clemilsonazevedo/blog/internal/service"
)

func RequireAuth(us *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenStr := ""
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "This route is missing a token!", http.StatusUnauthorized)
				return
			}

			tokenStr = cookie.Value
			if tokenStr == "" {
				http.Error(w, "This route is missing a token!", http.StatusUnauthorized)
				return
			}

			_, claim, err := auth.ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, "This token is invalid!", http.StatusUnauthorized)
				return
			}

			email, _ := claim["Email"].(string)
			if email == "" {
				http.Error(w, "This token is invalid!", http.StatusUnauthorized)
			}

			user, err := us.GetUserByEmail(email)
			if err != nil {
				http.Error(w, "User not found!", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
