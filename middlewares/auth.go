package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/config/secret"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func JWTAuth(us service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "This Route Missing a token", http.StatusUnauthorized)
				return
			}

			tokenStr = cookie.Value
			if tokenStr == "" {
				http.Error(w, "This Route Missing a token", http.StatusUnauthorized)
				return
			}

			_, claim, err := ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, "This Token is Invalid", http.StatusUnauthorized)
				return
			}

			email, _ := claim["email"].(string)
			if email == "" {
				http.Error(w, "This Token is Invalid", http.StatusUnauthorized)
			}

			user, err := us.GetUserByEmail(email)
			if err != nil {
				http.Error(w, "This Token is Invalid", http.StatusUnauthorized)
			}

			fmt.Fprintf(w, "%v", user)
			next.ServeHTTP(w, r)
		})
	}
}

func GenerateJWT(user entities.User, td time.Duration) (string, int64, error) {
	now := time.Now()
	exp := now.Add(td).Unix()
	claims := jwt.MapClaims{
		"id":        user.ID,
		"UserName":  user.UserName,
		"Email":     user.Email,
		"Role":      user.Role,
		"CreatedAt": user.CreatedAt,
		"iat":       now.Unix(),
		"exp":       exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret.GetJWTSecret()))

	return signed, exp, err
}

func ValidateJWT(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (any, error) {
			return []byte(secret.GetJWTSecret()), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, jwt.ErrInvalidKeyType
	}

	return token, claims, nil
}

func HashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func CheckPassword(hash string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
