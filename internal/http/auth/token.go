package auth

import (
	"time"

	"github.com/clemilsonazevedo/blog/config/secret"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

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
