package requests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"wiki/database"
	"wiki/utils"

	"github.com/aymanbagabas/go-udiff"
)

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

	// TODO: implement snapshots
	revId, err := utils.PushRevisionToDBFS(ctx, db, dataDir, revReq, rev.Content)
	if err != nil {
		fmt.Printf("Failure writing revision to db or fs: %s\n", err)
		return err
	}
	missingRevs, err := database.GetMissingRevisions(ctx, db, revId)
	if err != nil {
		return err
	}
	if len(missingRevs) >= 10 {
		_, err := utils.CreateSnapshot(ctx, db, dataDir, rev.PageId, revId)
		if err != nil {
			return err
		}
	}


	contentAtRev, err := utils.GetContentAtRevision(ctx, db, dataDir, rev.PageId, revId)
	if err != nil {
		return err
	}
	// also update the current page
	pageFilename := fmt.Sprintf("%s.md", rev.PageId)
	pageFilepath := filepath.Join(dataDir, "pages", pageFilename)
	os.WriteFile(pageFilepath, []byte(contentAtRev), 0644)
	// and the database stuff

	return nil

}
