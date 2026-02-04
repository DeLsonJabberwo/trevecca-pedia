package utils

import (
	"database/sql"
	"path/filepath"
	wikierrors "wiki/errors"
)

func GetDatabase() (*sql.DB, error) {
	var connStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	return db, nil
}

func GetDataDir() string {
	return filepath.Join("..", "wiki-fs")

}
