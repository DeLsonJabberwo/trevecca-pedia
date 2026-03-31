package utils

import (
	"database/sql"
	"fmt"
	"log"
	"moderation/config"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

func GetDatabase() (*sql.DB, error) {
	// Get connection parameters from environment variables with defaults

	var host, port, dbname, user, password string
	databaseUrl := config.GetEnv("DATABASE_URL", "")
	if databaseUrl == "" {
		host = config.GetEnv("MOD_DB_HOST", "localhost")
		port = config.GetEnv("MOD_DB_PORT", "5434")
		dbname = config.GetEnv("MOD_DB_NAME", "mod")
		user = config.GetEnv("MOD_DB_USER", "mod_user")
		password = config.GetEnv("MOD_DB_PASSWORD", "modpass")
	} else {
		host, port, dbname, user, password = parseDatabaseURL(databaseUrl)
	}

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %s\n", err)
		return nil, err
	}
	//log.Printf("Connected to database: %s\n", connStr)
	return db, nil
}

func parseDatabaseURL(databaseURL string) (host, port, dbname, user, password string) {
	// Set defaults
	host = "localhost"
	port = "5434"

	u, err := url.Parse(databaseURL)
	if err != nil {
		log.Printf("Failed to parse DATABASE_URL: %s\n", err)
		return
	}

	if u.Hostname() != "" {
		host = u.Hostname()
	}

	if u.Port() != "" {
		port = u.Port()
	}

	if u.Path != "" {
		dbname = strings.TrimPrefix(u.Path, "/")
	}

	if u.User != nil {
		user = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			password = pass
		}
	}

	return
}
