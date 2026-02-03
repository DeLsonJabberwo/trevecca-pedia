package requests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"wiki/database"
	"wiki/utils"
	wikierrors "wiki/errors"

	"github.com/aymanbagabas/go-udiff"
)

func DeletePage(ctx context.Context, db *sql.DB, dataDir string, delReq utils.DeletePageRequest) error {
	pageUUID, err := database.GetPageUUID(ctx, db, delReq.Slug)
	if err != nil {
		return wikierrors.DatabaseError(err)
	}
	pageInfo, err := database.GetPageInfo(ctx, db, pageUUID)
	if err != nil {
		return wikierrors.DatabaseError(err)
	}

	if pageInfo.DeletedAt != nil {
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

	currPage, err := GetPage(ctx, db, dataDir, revReq.PageId)
	if err != nil {
		return err
	}
	
	rev.PageId = currPage.UUID

	// probably validate user and page tokens and such
	rev.Author = revReq.Author

	// create the diff and make the revision
	filename := filepath.Join(dataDir, "pages", fmt.Sprintf("%s.md", rev.PageId))
	diff := udiff.Unified(filename, filename, currPage.Content, revReq.NewPage)

	rev.Content = diff

	revId, err := utils.PushRevisionToDBFS(ctx, db, dataDir, revReq, rev.Content)
	if err != nil {
		return wikierrors.DatabaseFilesystemError(err)
	}
	missingRevs, err := database.GetMissingRevisions(ctx, db, revId)
	if err != nil {
		return wikierrors.DatabaseError(err)
	}
	if len(missingRevs) >= 10 {
		_, err := utils.CreateSnapshot(ctx, db, dataDir, rev.PageId, revId)
		if err != nil {
			return wikierrors.DatabaseFilesystemError(err)
		}
	}


	contentAtRev, err := utils.GetContentAtRevision(ctx, db, dataDir, rev.PageId, revId)
	if err != nil {
		return wikierrors.DatabaseFilesystemError(err)
	}
	// also update the current page
	pageFilename := fmt.Sprintf("%s.md", rev.PageId)
	pageFilepath := filepath.Join(dataDir, "pages", pageFilename)
	err = os.WriteFile(pageFilepath, []byte(contentAtRev), 0644)
	if err != nil {
		return wikierrors.FilesystemError(err)
	}
	// and the database stuff

	return nil

}
