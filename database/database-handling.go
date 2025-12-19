package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func TestConnection(ctx context.Context, db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established.")
}

