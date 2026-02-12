package config

import (
	"os"

	"github.com/joho/godotenv"
)

var WikiServiceURL string
var SearchServiceURL string

func init() {
	godotenv.Load()
	WikiServiceURL = GetEnv("WIKI_SERVICE_URL", "http://127.0.0.1:9454")
	SearchServiceURL = GetEnv("SEARCH_SERVICE_URL", "http://127.0.0.1:7724")
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
	}
	return secret
}
