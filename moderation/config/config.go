package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AuthUrl string

func GetJWTSecret() string {
	secret := GetEnv("JWT_SECRET", "")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return secret
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	apiUrl := GetEnv("API_LAYER_URL", "http://127.0.0.1:2745/v1")
	AuthUrl = fmt.Sprintf("%s/auth", apiUrl)

}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
