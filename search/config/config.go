package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var WikiURL string

func init() {
	if err := godotenv.Load(); err != nil {

		log.Println("Warning: .env file not found, using defaults")
	}

	WikiURL = GetEnv("API_LAYER_URL", "http://127.0.0.1:2745/v1/wiki")
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
