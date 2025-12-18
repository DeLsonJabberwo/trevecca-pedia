package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var ConnStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"

func TestConnection(ctx context.Context, db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established.")
}

