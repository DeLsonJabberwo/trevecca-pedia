package requests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"wiki/database"
	"wiki/utils"
	"wiki/filesystem"
	wikierrors "wiki/errors"

	"github.com/aymanbagabas/go-udiff"
)

func DeletePage(ctx context.Context, db *sql.DB, dataDir string, delReq utils.DeletePageRequest) error {
	pageUUID, err := database.GetUUID(ctx, db, delReq.Slug)
	if err != nil {
		return wikierrors.DatabaseError(err)
	}
	pageInfo, err := database.GetPageInfo(ctx, db, pageUUID)
	if err != nil {
		return wikierrors.DatabaseError(err)
	}

	pageDeleted, err := database.GetPageDeleted(ctx, db, pageInfo.UUID) 
	if err != nil {
		return wikierrors.DatabaseError(err)
	}
	if pageDeleted {
		return wikierrors.PageDeleted()
	}

	// remove from database
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return wikierrors.InternalError(err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE pages
		SET deleted_at=NOW()
		WHERE uuid=$1;
	`, pageInfo.UUID)
	if err != nil {
		tx.Rollback()
		return wikierrors.DatabaseError(err)
	}

	tx.Commit()
	return nil
}

func PostRevision(ctx context.Context, db *sql.DB, dataDir string, revReq utils.RevisionRequest) error {
	var rev utils.Revision
	var err error

	fmt.Printf("PostRevision: starting for page %s\n", revReq.PageId)

	currPage, err := GetPage(ctx, db, dataDir, revReq.PageId)
	if err != nil {
		fmt.Printf("PostRevision: GetPage failed: %s\n", err)
		return err
	}
	fmt.Printf("PostRevision: got page %s\n", currPage.Slug)

	rev.PageId = currPage.UUID
	rev.Author = revReq.Author

	filename := filepath.Join(dataDir, "pages", fmt.Sprintf("%s.md", rev.PageId))
	diff := udiff.Unified(filename, filename, currPage.Content, revReq.NewPage)
	rev.Content = diff
	fmt.Printf("PostRevision: diff created\n")

	revId, err := utils.PushRevisionToDBFS(ctx, db, dataDir, revReq, rev.Content)
	if err != nil {
		fmt.Printf("PostRevision: PushRevisionToDBFS failed: %s\n", err)
		return wikierrors.DatabaseFilesystemError(err)
	}
	fmt.Printf("PostRevision: revision pushed, id=%s\n", revId)

	missingRevs, err := database.GetMissingRevisions(ctx, db, revId)
	if err != nil {
		fmt.Printf("PostRevision: GetMissingRevisions failed: %s\n", err)
		return wikierrors.DatabaseError(err)
	}
	fmt.Printf("PostRevision: missing revs count=%d\n", len(missingRevs))

	if len(missingRevs) >= 10 {
		_, err := utils.CreateSnapshot(ctx, db, dataDir, rev.PageId, revId)
		if err != nil {
			
			fmt.Printf("PostRevision: CreateSnapshot failed (non-fatal): %s\n", err)
		}
	}

	slugPath, err := filesystem.GetPageFilename(ctx, db, dataDir, rev.PageId.String())
	if err != nil {
		fmt.Printf("PostRevision: GetPageFilename failed: %s\n", err)
		return wikierrors.FilesystemError(err)
	}
	fmt.Printf("PostRevision: writing to %s\n", slugPath)

	err = os.WriteFile(slugPath, []byte(revReq.NewPage), 0644)
	if err != nil {
		fmt.Printf("PostRevision: WriteFile failed: %s\n", err)
		return wikierrors.FilesystemError(err)
	}
	fmt.Printf("PostRevision: file written\n")

	_, err = db.ExecContext(ctx, `
		UPDATE pages
		SET last_revision_id = $1
			
		WHERE uuid = $2;
	`, revId, rev.PageId)
	if err != nil {
		fmt.Printf("PostRevision: db update failed: %s\n", err)
		return wikierrors.DatabaseError(fmt.Errorf("updating page last_revision_id: %w", err))
	}
	fmt.Printf("PostRevision: complete\n")

	return nil
}