package requests

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"wiki/database"
	wikierrors "wiki/errors"
	"wiki/filesystem"
	"wiki/utils"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/google/uuid"
)

func GetPage(ctx context.Context, db *sql.DB, dataDir string, id string) (utils.Page, error) {
	var page utils.Page
	var info *database.PageInfo
	var content string
	var lastRev *database.RevInfo
	var pageId uuid.UUID
	var err error

	if err = uuid.Validate(id); err == nil {
		pageId, err = uuid.Parse(id)
		if err != nil {
			return utils.Page{}, err
		}
	} else {
		pageId, err = database.GetPageUUID(ctx, db, id)
		if err != nil {
			return utils.Page{}, err
		}
	}

	info, err = database.GetPageInfo(ctx, db, pageId)
	if err != nil {
		return utils.Page{}, err
	}

	content, err = filesystem.GetPageContent(ctx, db, dataDir, pageId)
	if err != nil {
		return utils.Page{}, err
	}

	if info.LastRevisionId != nil {
		lastRev, err = database.GetRevisionInfo(ctx, db, *info.LastRevisionId)
		if err != nil {
			return utils.Page{}, err
		}
	} else {
		lastRev = &database.RevInfo{}
	}

	page = utils.Page{UUID: info.UUID, Slug: info.Slug, Name: info.Name, ArchiveDate: info.ArchiveDate,
		DeletedAt: info.DeletedAt, LastEdit: lastRev.UUID, LastEditTime: lastRev.DateTime,
		Content: content}

	if page.DeletedAt != nil {
		return utils.Page{}, wikierrors.PageDeleted(pageId.String())
	}

	return page, nil
}

func GetPages(ctx context.Context, db *sql.DB, ind int, num int) ([]database.PageInfo, error) {
	var pages []database.PageInfo = make([]database.PageInfo, num)

	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pages").Scan(&count)
	if count != 0 && err != nil {
		return nil, err
	}
	count -= ind
	if count <= 0 {
		return pages, nil
	}

	uuids, err := db.QueryContext(ctx,
		"SELECT uuid FROM pages")
	if err != nil {
		return nil, err
	}
	for range ind {
		uuids.Next()
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

func GetPagesCategory(ctx context.Context, db *sql.DB, cat int, ind int, num int) ([]database.PageInfo, error) {
	var pages []database.PageInfo = make([]database.PageInfo, num)

	var count int
	err := db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM pages
		JOIN page_categories ON pages.uuid = page_categories.page_id
		WHERE page_categories.category=$1`,
		cat).Scan(&count)
	if err != nil {
		return nil, err
	}
	count -= ind
	if count <= 0 {
		return pages, nil
	}

	uuids, err := db.QueryContext(
		ctx,
		`SELECT uuid FROM pages
		JOIN page_categories ON pages.uuid = page_categories.page_id
		WHERE page_categories.category=$1`,
		cat)
	if err != nil {
		return nil, err
	}
	for range ind {
		uuids.Next()
	}

	for i := range pages {
		uuids.Next()
		if uuids == nil {
			break
		}
		var id uuid.UUID
		uuids.Scan(&id)
		pageInfo, err := database.GetPageInfo(ctx, db, id)
		if err != nil {
			i--
			continue
		}
		pages[i] = *pageInfo
	}

	return pages, nil
}

func GetRevision(ctx context.Context, db *sql.DB, dataDir string, revId string) (utils.Revision, error) {
	var err error
	var rev = utils.Revision{}

	if err := uuid.Validate(revId); err != nil {
		return utils.Revision{}, err
	}
	rev.UUID, err = uuid.Parse(revId)
	if err != nil {
		return utils.Revision{}, err
	}

	revInfo, err := database.GetRevisionInfo(ctx, db, rev.UUID)
	if err != nil {
		return utils.Revision{}, err
	}
	if revInfo == nil || revInfo.PageId == nil {
		return utils.Revision{}, err
	}
	pageInfo, err := database.GetPageInfo(ctx, db, *revInfo.PageId)
	if err != nil {
		return utils.Revision{}, err
	}
	if pageInfo.DeletedAt != nil {
		return utils.Revision{}, err
	}
	rev.PageId = *revInfo.PageId
	rev.Name = pageInfo.Name
	rev.RevDateTime = *revInfo.DateTime

	lastSnap, err := database.GetMostRecentSnapshot(ctx, db, rev.UUID)
	if err != nil {
		return utils.Revision{}, err
	}
	if lastSnap == nil {
		return utils.Revision{}, err
	}
	missingRevs, err := database.GetMissingRevisions(ctx, db, rev.UUID)
	if err != nil {
		return utils.Revision{}, err
	}
	rev.Content, err = filesystem.GetSnapshotContent(ctx, db, dataDir, lastSnap.UUID)
	if err != nil {
		return utils.Revision{}, err
	}

	// i hope and pray that this works
	// update: it worked. most errors were elsewhere :)
	for _, r := range missingRevs {
		if r.UUID == nil {
			continue
		}
		revContent, err := filesystem.GetRevisionContent(ctx, db, dataDir, *r.UUID)
		if err != nil {
			return utils.Revision{}, err
		}
		files, _, err := gitdiff.Parse(bytes.NewReader([]byte(revContent)))
		if err != nil {
			return utils.Revision{}, err
		}
		if len(files) == 0 {
			continue
		}
		src := bytes.NewReader([]byte(rev.Content))
		var dst bytes.Buffer

		err = gitdiff.Apply(&dst, src, files[0])
		if err != nil {
			if errors.Is(err, &gitdiff.Conflict{}) {
				return utils.Revision{}, err
			}
			return utils.Revision{}, err
		}
		rev.Content = dst.String()
	}

	return rev, nil
}

func GetRevisions(ctx context.Context, db *sql.DB, pageId string, ind int, num int) ([]database.RevInfo, error) {
	var revs []database.RevInfo = make([]database.RevInfo, num)

	var pageUUID uuid.UUID
	if err := uuid.Validate(pageId); err == nil {
		pageUUID, err = uuid.Parse(pageId)
		if err != nil {
			return nil, err
		}
	} else {
		pageUUID, err = database.GetPageUUID(ctx, db, pageId)
		if err != nil {
			return nil, err
		}
	}

	pageInfo, err := database.GetPageInfo(ctx, db, pageUUID)
	if err != nil {
		return nil, err
	}
	if pageInfo.DeletedAt != nil {
		return nil, err
	}

	var count int
	err = db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM revisions WHERE page_id=$1",
		pageUUID).Scan(&count)
	if count != 0 && err != nil {
		return nil, err
	}
	count -= ind
	if count <= 0 {
		return revs, nil
	}

	uuids, err := db.QueryContext(
		ctx,
		"SELECT uuid FROM revisions WHERE page_id=$1",
		pageUUID)
	if err != nil {
		return nil, err
	}
	for range ind {
		uuids.Next()
	}

	for i := range revs {
		uuids.Next()
		if uuids == nil {
			break
		}
		var id uuid.UUID
		uuids.Scan(&id)
		revInfo, err := database.GetRevisionInfo(ctx, db, id)
		if err != nil {
			return nil, err
		}
		revs[i] = *revInfo
	}

	return revs, nil
}
