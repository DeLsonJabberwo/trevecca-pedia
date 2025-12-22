package requests

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"wiki/database"
	"wiki/filesystem"

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

	content, err = filesystem.GetPage(dataDir, pageId)
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
			log.Printf("error: %s\n", err)
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
