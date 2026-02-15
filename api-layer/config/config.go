package config

import (
	"github.com/joho/godotenv"
	"os"
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
