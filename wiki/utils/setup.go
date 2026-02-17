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
	// Only load .env file in local development (not on fly.io)
	// Fly.io sets FLY_APP_NAME environment variable
	if os.Getenv("FLY_APP_NAME") == "" {
		// Try to load .env file, but don't fail if it doesn't exist
		// This allows environment variables to be set via other means (Docker, shell, etc.)
		_ = godotenv.Load(".env")
	}
}

func GetDatabase() (*sql.DB, error) {
	// Check if running on fly.io (production)
	isProduction := os.Getenv("FLY_APP_NAME") != ""

	// First check for DATABASE_URL (set by fly postgres attach)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, wikierrors.DatabaseError(err)
		}
		return db, nil
	}

	// Fall back to individual WIKI_DB_* environment variables
	host := os.Getenv("WIKI_DB_HOST")
	port := os.Getenv("WIKI_DB_PORT")
	dbname := os.Getenv("WIKI_DB_NAME")
	user := os.Getenv("WIKI_DB_USER")
	password := os.Getenv("WIKI_DB_PASSWORD")

	// In production (fly.io), fail fast if required secrets are missing
	if isProduction {
		var missing []string
		if host == "" {
			missing = append(missing, "WIKI_DB_HOST")
		}
		if port == "" {
			missing = append(missing, "WIKI_DB_PORT")
		}
		if password == "" {
			missing = append(missing, "WIKI_DB_PASSWORD")
		}
		if dbname == "" {
			missing = append(missing, "WIKI_DB_NAME")
		}
		if user == "" {
			missing = append(missing, "WIKI_DB_USER")
		}

		if len(missing) > 0 {
			return nil, fmt.Errorf("missing required database secrets on fly.io: %v. Set either DATABASE_URL (recommended: fly postgres attach) or individual WIKI_DB_* secrets",
				missing)
		}
	} else {
		// Local development: use defaults if not set
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if dbname == "" {
			dbname = "wiki"
		}
		if user == "" {
			user = "wiki_user"
		}
		if password == "" {
			password = "myatt"
		}
	}

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
	if dataDir := os.Getenv("WIKI_DATA_DIR"); dataDir != "" {
		return dataDir
	}
	return filepath.Join("..", "wiki-fs")
}
