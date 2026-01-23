package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func PushRevisionToDBFS(ctx context.Context, db *sql.DB, dataDir string, revReq RevisionRequest, diff string) (uuid.UUID, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback()

	var revUUID uuid.UUID
	err = tx.QueryRowContext(ctx, `
			INSERT INTO revisions (page_id, author)
			VALUES ($1, $2)
			RETURNING uuid;
			`, revReq.PageId, revReq.Author).Scan(&revUUID)
	if err != nil {
		fmt.Printf("Error writing to db: %s\n", err)
		return uuid.UUID{}, err
	}
	
	filename := fmt.Sprintf("%s.txt", revUUID)
	err = os.MkdirAll(filepath.Join(dataDir, "revisions"), 0755)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("creating revisions directory: %w", err)
	}
	err = os.WriteFile(filepath.Join(dataDir, "revisions", filename), []byte(diff), 0644)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = tx.Commit()
	if err != nil {
		os.Remove(filepath.Join(dataDir, "revisions", filename))
		return uuid.UUID{}, err
	}
	return revUUID, nil
}
