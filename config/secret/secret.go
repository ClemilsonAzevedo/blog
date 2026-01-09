package secret

import "os"

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "in-production-change-this-please"
	}
	return secret
}

func GetOpenAiKey() string {
	api_key := os.Getenv("OPENAI_API_KEY")
	if api_key == "" {
		api_key = "your-openai-api-key-here"
	}
	return api_key
}
