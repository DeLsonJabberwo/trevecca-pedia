package utils

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"wiki/database"
	wikierrors "wiki/errors"
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
	lastSnap, err := database.GetMostRecentSnapshot(ctx, db, revId)
	if err == sql.ErrNoRows {
		return "", wikierrors.RevisionNotFound()
	}
	if err != nil {
		return "", wikierrors.DatabaseError(err)
	}
	missingRevs, err := database.GetMissingRevisions(ctx, db, revId)
	if err != nil {
		return "", wikierrors.DatabaseError(err)
	}
	revContent, err := filesystem.GetSnapshotContent(ctx, db, dataDir, lastSnap.UUID)
	if err != nil {
		return "", wikierrors.FilesystemError(err)
	}

	// i hope and pray that this works
	// update: it worked. most errors were elsewhere :)
	for _, r := range missingRevs {
		revContent, err := filesystem.GetRevisionContent(ctx, db, dataDir, *r.UUID)
		if err != nil {
			return "", wikierrors.FilesystemError(err)
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

func GetPageInfoPreview(ctx context.Context, db *sql.DB, dataDir string, pageId uuid.UUID) (*PageInfoPrev, error) {
	pageInfo, err := database.GetPageInfo(ctx, db, pageId)
	if err != nil {
		return nil, err
	}
	preview, err := filesystem.GetPagePreview(ctx, db, dataDir, pageId, 250)
	if err != nil {
		return nil, err
	}
	var lastEditTime time.Time
	if pageInfo.LastRevisionId == nil {
		pageInfo.LastRevisionId = &uuid.Nil
		lastEditTime = time.Time{}
	} else {
		err := db.QueryRowContext(ctx, `
			SELECT date_time FROM revisions WHERE uuid=$1;
		`, pageInfo.LastRevisionId).Scan(&lastEditTime)
		if err != nil {
			return nil, err
		}
	}
	if pageInfo.ArchiveDate == nil {
		pageInfo.ArchiveDate = &time.Time{}
	}
	return &PageInfoPrev{
		UUID: pageInfo.UUID,
		Slug: pageInfo.Slug,
		Name: pageInfo.Name,
		LastEditTime: lastEditTime,
		ArchiveDate: *pageInfo.ArchiveDate,
		Preview: preview,
	}, nil
}

func GetIndexInfo(ctx context.Context, db *sql.DB, dataDir string, pageId string) (*IndexInfo, error) {
	pageUUID, err := database.GetUUID(ctx, db, pageId)
	if err != nil {
		return nil, err
	}
	var indexInfo IndexInfo
	var lastRev uuid.UUID
	var archiveDate *time.Time
	err = db.QueryRowContext(ctx, `
		SELECT slug, name, last_revision_id, archive_date
		FROM pages WHERE uuid=$1;
	`, pageUUID).Scan(&indexInfo.Slug, &indexInfo.Name, &lastRev, &archiveDate)
	if err != nil {
		return nil, err
	}
	fmt.Printf("lastRev: %s\n", lastRev)
	fmt.Printf("archiveDate: %s\n", archiveDate)
	if lastRev != uuid.Nil {
		err = db.QueryRowContext(ctx, `
		SELECT date_time FROM revisions WHERE uuid=$1; 
		`, lastRev).Scan(&indexInfo.LastModified)
		if err != nil {
			return nil, err
		}
	} else {
		indexInfo.LastModified = time.Time{}
	}
	if archiveDate != nil {
		indexInfo.ArchiveDate = *archiveDate
	} else {
		indexInfo.ArchiveDate = time.Time{}
	}
	indexInfo.Content, err = filesystem.GetPageContent(ctx, db, dataDir, pageUUID)
	if err != nil {
		return nil, err
	}
	return &indexInfo, nil
}

