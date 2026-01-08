package database

import (
	"context"
	"database/sql"
)

func ValidateCategory(ctx context.Context, db *sql.DB, cat string) int {
	id := 0
	err := db.QueryRowContext(
		ctx,
		`SELECT id FROM categories
		WHERE id=$1`,
		cat).Scan(&id)
	if err != nil {
		id = 0
	}
	if id != 0 {
		return id
	}
	err = db.QueryRowContext(
		ctx,
		`SELECT id FROM categories
		WHERE slug=$1`, cat).Scan(&id)
	if err != nil {
		id = 0
	}
	return id
}
