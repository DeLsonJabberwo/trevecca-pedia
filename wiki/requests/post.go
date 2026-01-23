package requests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
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
	_, err = utils.PushRevisionToDBFS(ctx, db, dataDir, revReq, rev.Content)
	if err != nil {
		fmt.Printf("Failure writing revision to db or fs: %s\n", err)
		return err
	}

	// also update the current page
	pageFilename := fmt.Sprintf("%s.md", rev.PageId)
	pageFilepath := filepath.Join(dataDir, "pages", pageFilename)
	os.WriteFile(pageFilepath, []byte(revReq.NewPage), 0644)
	// and the database stuff

	return nil

}
