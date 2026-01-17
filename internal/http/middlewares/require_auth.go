package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
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
				exceptions.Unauthorized(w, "Token is missing!")
				return
			}

			tokenStr = cookie.Value
			if tokenStr == "" {
				exceptions.BadRequest(w, errors.New("Request Error"), "The token is empty.", nil)
				return
			}

			_, claim, err := auth.ValidateJWT(tokenStr)
			if err != nil {
				exceptions.BadRequest(w, errors.New("Request Error"), "The token is invalid.", nil)
				return
			}

			email, _ := claim["Email"].(string)
			if email == "" {
				exceptions.Unauthorized(w, "Unauthorized account.")
			}

			user, err := us.GetUserByEmail(email)
			if err != nil {
				exceptions.BadRequest(w, errors.New("Request Error"), "User not found.", nil)
				return
			}

			ctx = context.WithValue(ctx, "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
