package filesystem

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/uuid"
)

func GetPageContent(ctx context.Context, db *sql.DB, dataDir string, pageId uuid.UUID) (string, error) {
	filename, err := GetPageFilename(ctx, db, dataDir, pageId.String())
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetRevisionContent(ctx context.Context, db *sql.DB, dataDir string, revId uuid.UUID) (string, error) {
	if revId == uuid.Nil {
		return "", nil
	}
	filename, err := GetRevisionFilename(ctx, db, dataDir, revId)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetSnapshotContent(ctx context.Context, db *sql.DB, dataDir string, snapId uuid.UUID) (string, error) {
	filename, err := GetSnapshotFilename(ctx, db, dataDir, snapId)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
