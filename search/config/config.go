package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var WikiURL string
var IndexDir string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	apiURL := GetEnv("API_LAYER_URL", "http://127.0.0.1:2745/v1")
	WikiURL = fmt.Sprintf("%s/wiki", apiURL)
	IndexDir = GetEnv("INDEX_DIR", "../index")
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
