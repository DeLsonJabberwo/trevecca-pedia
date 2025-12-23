package requests

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"wiki/database"
	"wiki/filesystem"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/google/uuid"
)

func GetPage(ctx context.Context, db *sql.DB, dataDir string, id string) (Page, error) {
	var page Page
	var info *database.PageInfo
	var content string
	var lastRev *database.RevInfo
	var pageId uuid.UUID
	var err error

	if err = uuid.Validate(id); err == nil {
		pageId, err = uuid.Parse(id)
		if err != nil {
			return Page{}, err
		}
	} else {
		pageId, err = database.GetPageUUID(ctx, db, id)
		if err != nil {
			return Page{}, err
		}
	}

	info, err = database.GetPageInfo(ctx, db, pageId)
	if err != nil {
		return Page{}, err
	}

	content, err = filesystem.GetPageContent(dataDir, pageId)
	if err != nil {
		return Page{}, err
	}

	if info.LastRevisionId != nil {
		lastRev = database.GetRevisionInfo(ctx, db, *info.LastRevisionId)
	} else {
		lastRev = &database.RevInfo{}
	}

	page = Page{info.UUID, info.Slug, info.Name, info.ArchiveDate, info.DeletedAt, lastRev.UUID, lastRev.DateTime, content}

	if page.DeletedAt != nil {
		return Page{}, errors.New("not found")
	}

	return page, nil
}

func GetPages(ctx context.Context, db *sql.DB, ind int, num int) ([]database.PageInfo, error) {
	var pages []database.PageInfo = make([]database.PageInfo, num)

	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pages OFFSET $1;", ind).Scan(&count)
	if count != 0 && err != nil {
		return nil, err
	}
	if count == 0 {
		return pages, nil
	}

	uuids, err := db.QueryContext(ctx,
				"SELECT uuid FROM pages OFFSET $1;", ind)
	if err != nil {
		return nil, err
	}

	// i can almost guarantee this disgusting loop
	// could be done better and be much less tragic
	uuids.Next()
	for i := 0; i < num && i < count; i++ {
		var id uuid.UUID
		uuids.Scan(&id)
		pageInfo, err := database.GetPageInfo(ctx, db, id)
		uuids.Next()
		if pageInfo == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		if pageInfo.DeletedAt != nil {
			i--
			continue
		}
		pages[i] = *pageInfo
	}

	return pages, nil
}

func GetRevision(ctx context.Context, db *sql.DB, dataDir string, revId string) (Revision, error) {
	var err error
	var rev = Revision{}

	if err := uuid.Validate(revId); err != nil {
		return Revision{}, err
	}
	rev.UUID, err = uuid.Parse(revId)
	if err != nil {
		return Revision{}, err
	}

	revInfo := database.GetRevisionInfo(ctx, db, rev.UUID)
	pageInfo, err := database.GetPageInfo(ctx, db, *revInfo.PageId)
	if err != nil {
		return Revision{}, err
	}
	if pageInfo.DeletedAt != nil {
		return Revision{}, errors.New("404")
	}
	rev.PageId = *revInfo.PageId
	rev.Name = pageInfo.Name
	rev.RevDateTime = *revInfo.DateTime

	lastSnap := database.GetMostRecentSnapshot(ctx, db, rev.UUID)
	missingRevs, err := database.GetMissingRevisions(ctx, db, rev.UUID)
	if err != nil {
		return Revision{}, err
	}
	rev.Content, err = filesystem.GetSnapshotContent(dataDir, lastSnap.UUID)
	if err != nil {
		return Revision{}, err
	}

	// i hope and pray that this works
	// update: it worked. most errors were elsewhere :)
	for _, r := range missingRevs {
		revContent, err := filesystem.GetRevisionContent(dataDir, *r.UUID)
		if err != nil {
			return Revision{}, err
		}
		files, _, err := gitdiff.Parse(bytes.NewReader([]byte(revContent)))
		if err != nil {
			return Revision{}, fmt.Errorf("couldn't parse revision: %w", err)
		}
		if len(files) == 0 {
			continue
		}
		src := bytes.NewReader([]byte(rev.Content))
		var dst bytes.Buffer

		err = gitdiff.Apply(&dst, src, files[0])
		if err != nil {
			if errors.Is(err, &gitdiff.Conflict{}) {
				return Revision{}, fmt.Errorf("conflict while applying revision: %w", err)
			}
			return Revision{}, fmt.Errorf("applying revision: %w", err)
		}
		rev.Content = dst.String()
	}

	return rev, nil
}
