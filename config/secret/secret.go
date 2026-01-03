package secret

import "os"

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "in-production-change-this-please"
	}
	return secret
}
