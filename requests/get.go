package requests

import (
	"context"
	"database/sql"
	"errors"
	"wiki/database"
	"wiki/filesystem"

	"github.com/google/uuid"
)

func GetPage(ctx context.Context, db *sql.DB, dataDir string, id string) (Page, error) {
	var page Page
	var info *database.PageInfo
	var content string
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

	page = Page{info.UUID, info.Slug, info.Name, info.ArchiveDate, info.DeletedAt, content}

	if page.DeletedAt != nil {
		return Page{}, errors.New("not found")
	}

	return page, nil
}
