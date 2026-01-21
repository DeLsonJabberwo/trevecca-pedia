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
	_, err := db.Exec(`
			INSERT INTO revisions (page_id, author)
			VALUES ($1, $2);
			`, revReq.PageId, revReq.Author)
	if err != nil {
		fmt.Printf("Error writing to db: %s\n", err)
		return uuid.UUID{}, err
	}
	var revUUID uuid.UUID
	err = db.QueryRowContext(ctx, `
			SELECT uuid FROM revisions
			WHERE page_id=$1
			ORDER BY date_time DESC
			LIMIT 1;
			`, revReq.PageId).Scan(&revUUID)
	if err != nil {
		return uuid.UUID{}, err
	}
	
	filename := fmt.Sprintf("%s.txt", revUUID)
	os.MkdirAll(filepath.Join(dataDir, "revisions"), 0755)
	err = os.WriteFile(filepath.Join(dataDir, "revisions", filename), []byte(diff), 0644)
	if err != nil {
		log.Printf("Couldn't write revision to file: %s\n", err)
		return uuid.UUID{}, err
	}

	return revUUID, nil
}
