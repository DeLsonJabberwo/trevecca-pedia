package utils

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"wiki/database"
	"wiki/filesystem"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
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

func CreateSnapshot(ctx context.Context, db *sql.DB, dataDir string, pageId uuid.UUID, revId uuid.UUID) (uuid.UUID, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	var snapUUID uuid.UUID
	err = tx.QueryRowContext(ctx, `
			INSERT INTO snapshots (page, revision)
			VALUES ($1, $2)
			RETURNING uuid;
			`, pageId, revId).Scan(snapUUID)
	if err != nil {
		return uuid.UUID{}, err
	}

	snapContent, err := GetContentAtRevision(ctx, db, dataDir, pageId, revId)
	if err != nil {
		return uuid.UUID{}, err
	}

	filename := fmt.Sprintf("%s.md", snapUUID)
	err = os.MkdirAll(filepath.Join(dataDir, "snapshots"), 0755)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("creating snapshots directory: %w", err)
	}
	err = os.WriteFile(filepath.Join(dataDir, "snapshots", filename), []byte(snapContent), 0644)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = tx.Commit()
	if err != nil {
		os.Remove(filepath.Join(dataDir, "snapshots", filename))
		return uuid.UUID{}, err
	}
	return uuid.UUID{}, nil
}

func GetContentAtRevision(ctx context.Context, db *sql.DB, dataDir string, pageId uuid.UUID, revId uuid.UUID) (string, error) {
	lastSnap := database.GetMostRecentSnapshot(ctx, db, revId)
	missingRevs, err := database.GetMissingRevisions(ctx, db, revId)
	if err != nil {
		return "", err
	}
	revContent, err := filesystem.GetSnapshotContent(ctx, db, dataDir, lastSnap.UUID)
	if err != nil {
		return "", err
	}

	// i hope and pray that this works
	// update: it worked. most errors were elsewhere :)
	for _, r := range missingRevs {
		revContent, err := filesystem.GetRevisionContent(ctx, db, dataDir, *r.UUID)
		if err != nil {
			return "", err
		}
		files, _, err := gitdiff.Parse(bytes.NewReader([]byte(revContent)))
		if err != nil {
			return "", fmt.Errorf("couldn't parse revision: %w", err)
		}
		if len(files) == 0 {
			continue
		}
		src := bytes.NewReader([]byte(revContent))
		var dst bytes.Buffer

		err = gitdiff.Apply(&dst, src, files[0])
		if err != nil {
			if errors.Is(err, &gitdiff.Conflict{}) {
				return "", fmt.Errorf("conflict while applying revision: %w", err)
			}
			return "", fmt.Errorf("applying revision: %w", err)
		}
		revContent = dst.String()
	}
	return revContent, nil
}
