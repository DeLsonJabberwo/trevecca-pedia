package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	wikierrors "wiki/errors"

	"github.com/joho/godotenv"
)

// init loads the .env file when the package is imported
func init() {
	// Try to load .env file, but don't fail if it doesn't exist
	// This allows environment variables to be set via other means (Docker, shell, etc.)
	_ = godotenv.Load(".env")
}

func GetDatabase() (*sql.DB, error) {
	// Get connection parameters from environment variables with defaults
	host := getEnv("WIKI_DB_HOST", "localhost")
	port := getEnv("WIKI_DB_PORT", "5432")
	dbname := getEnv("WIKI_DB_NAME", "wiki")
	user := getEnv("WIKI_DB_USER", "wiki_user")
	password := getEnv("WIKI_DB_PASSWORD", "myatt")

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	return db, nil
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDataDir() string {
	return filepath.Join("..", "wiki-fs")

}
