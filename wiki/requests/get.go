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
			return utils.Page{}, wikierrors.InvalidID(err)
		}
	} else {
		pageId, err = database.GetPageUUID(ctx, db, id)
		if err == sql.ErrNoRows {
			return utils.Page{}, wikierrors.PageNotFound()
		}
		if err != nil {
			return utils.Page{}, wikierrors.DatabaseError(err)
		}
	}

	info, err = database.GetPageInfo(ctx, db, pageId)
	if err != nil {
		return utils.Page{}, wikierrors.DatabaseError(err)
	}

	if info.DeletedAt != nil {
		return utils.Page{}, wikierrors.PageDeleted()
	}

	content, err = filesystem.GetPageContent(ctx, db, dataDir, pageId)
	if err != nil {
		return utils.Page{}, wikierrors.FilesystemError(err)
	}

	if info.LastRevisionId != nil {
		lastRev, err = database.GetRevisionInfo(ctx, db, *info.LastRevisionId)
		if err == sql.ErrNoRows {
			return utils.Page{}, wikierrors.RevisionNotFound()
		}
		if err != nil {
			return utils.Page{}, wikierrors.DatabaseError(err)
		}
	} else {
		lastRev = &database.RevInfo{}
	}

	page = utils.Page{UUID: info.UUID, Slug: info.Slug, Name: info.Name, ArchiveDate: info.ArchiveDate,
		DeletedAt: info.DeletedAt, LastEdit: lastRev.UUID, LastEditTime: lastRev.DateTime,
		Content: content}

	return page, nil
}

func GetPages(ctx context.Context, db *sql.DB, ind int, count int) ([]database.PageInfo, error) {
	var pagesCount int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pages WHERE deleted_at IS NULL").Scan(&pagesCount)
	if pagesCount != 0 && err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	pagesCount -= ind
	if pagesCount <= 0 {
		return []database.PageInfo{}, nil
	}

	uuids, err := db.QueryContext(ctx,
		"SELECT uuid FROM pages WHERE deleted_at IS NULL")
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	for range ind {
		uuids.Next()
	}

	var pages []database.PageInfo
	if count <= pagesCount {
		pages = make([]database.PageInfo, count)
	} else {
		pages = make([]database.PageInfo, pagesCount)
	}

	// i can almost guarantee this disgusting loop
	// could be done better and be much less tragic
	uuids.Next()
	for i := 0; i < len(pages); i++ {
		var id uuid.UUID
		uuids.Scan(&id)
		pageInfo, err := database.GetPageInfo(ctx, db, id)
		uuids.Next()
		if pageInfo == nil {
			continue
		}
		if err != nil {
			return nil, wikierrors.DatabaseError(err)
		}
		pages[i] = *pageInfo
	}

	return pages, nil
}

func GetPagesCategory(ctx context.Context, db *sql.DB, cat int, ind int, count int) ([]database.PageInfo, error) {
	var pagesCount int
	err := db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM pages
		JOIN page_categories ON pages.uuid = page_categories.page_id
		WHERE page_categories.category=$1 AND pages.deleted_at IS NULL`,
		cat).Scan(&pagesCount)
	if pagesCount != 0 && err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	pagesCount -= ind
	if pagesCount <= 0 {
		return []database.PageInfo{}, nil
	}

	uuids, err := db.QueryContext(
		ctx,
		`SELECT uuid FROM pages
		JOIN page_categories ON pages.uuid = page_categories.page_id
		WHERE page_categories.category=$1 AND pages.deleted_at IS NULL`,
		cat)
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	for range ind {
		uuids.Next()
	}

	var pages []database.PageInfo
	if count <= pagesCount {
		pages = make([]database.PageInfo, count)
	} else {
		pages = make([]database.PageInfo, pagesCount)
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
		return utils.Revision{}, wikierrors.InvalidID(err)
	}
	rev.UUID, err = uuid.Parse(revId)
	if err != nil {
		return utils.Revision{}, wikierrors.InvalidID(err)
	}

	revInfo, err := database.GetRevisionInfo(ctx, db, rev.UUID)
	if err == sql.ErrNoRows {
		return utils.Revision{}, wikierrors.RevisionNotFound()
	}
	if err != nil {
		return utils.Revision{}, wikierrors.DatabaseError(err)
	}
	if revInfo == nil || revInfo.PageId == nil {
		return utils.Revision{}, wikierrors.RevisionNotFound()
	}
	pageInfo, err := database.GetPageInfo(ctx, db, *revInfo.PageId)
	if err == sql.ErrNoRows {
		return utils.Revision{}, wikierrors.PageNotFound()
	}
	if err != nil {
		return utils.Revision{}, wikierrors.DatabaseError(err)
	}
	if pageInfo.DeletedAt != nil {
		return utils.Revision{}, wikierrors.PageDeleted()
	}
	rev.PageId = *revInfo.PageId
	rev.Name = pageInfo.Name
	rev.RevDateTime = *revInfo.DateTime

	lastSnap, err := database.GetMostRecentSnapshot(ctx, db, rev.UUID)
	if err == sql.ErrNoRows{
		return utils.Revision{}, wikierrors.SnapshotNotFound()
	}
	if err != nil {
		return utils.Revision{}, wikierrors.DatabaseError(err)
	}
	if lastSnap == nil {
		return utils.Revision{}, wikierrors.SnapshotNotFound()
	}
	missingRevs, err := database.GetMissingRevisions(ctx, db, rev.UUID)
	if err != nil {
		return utils.Revision{}, wikierrors.DatabaseError(err)
	}
	rev.Content, err = filesystem.GetSnapshotContent(ctx, db, dataDir, lastSnap.UUID)
	if err != nil {
		return utils.Revision{}, wikierrors.FilesystemError(err)
	}

	// i hope and pray that this works
	// update: it worked. most errors were elsewhere :)
	for _, r := range missingRevs {
		if r.UUID == nil {
			continue
		}
		revContent, err := filesystem.GetRevisionContent(ctx, db, dataDir, *r.UUID)
		if err != nil {
			return utils.Revision{}, wikierrors.FilesystemError(err)
		}
		files, _, err := gitdiff.Parse(bytes.NewReader([]byte(revContent)))
		if err != nil {
			return utils.Revision{}, wikierrors.InternalError(err)
		}
		if len(files) == 0 {
			continue
		}
		src := bytes.NewReader([]byte(rev.Content))
		var dst bytes.Buffer

		err = gitdiff.Apply(&dst, src, files[0])
		if err != nil {
			if errors.Is(err, &gitdiff.Conflict{}) {
				return utils.Revision{}, wikierrors.RevisionConflict(err)
			}
			return utils.Revision{}, wikierrors.InternalError(err)
		}
		rev.Content = dst.String()
	}

	return rev, nil
}

func GetRevisions(ctx context.Context, db *sql.DB, pageId string, ind int, count int) ([]database.RevInfo, error) {
	var revs []database.RevInfo = make([]database.RevInfo, count)

	var pageUUID uuid.UUID
	if err := uuid.Validate(pageId); err == nil {
		pageUUID, err = uuid.Parse(pageId)
		if err != nil {
			return nil, wikierrors.InvalidID(err)
		}
	} else {
		pageUUID, err = database.GetPageUUID(ctx, db, pageId)
		if err == sql.ErrNoRows {
			return nil, wikierrors.PageNotFound()
		}
		if err != nil {
			return nil, wikierrors.DatabaseError(err)
		}
	}

	pageInfo, err := database.GetPageInfo(ctx, db, pageUUID)
	if err == sql.ErrNoRows {
		return nil, wikierrors.PageNotFound()
	}
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	if pageInfo.DeletedAt != nil {
		return nil, wikierrors.PageDeleted()
	}

	var rowCount int
	err = db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM revisions WHERE page_id=$1",
		pageUUID).Scan(&rowCount)
	if rowCount != 0 && err != nil {
		return nil, wikierrors.DatabaseError(err)
	}
	rowCount -= ind
	if rowCount <= 0 {
		return revs, nil
	}

	uuids, err := db.QueryContext(
		ctx,
		"SELECT uuid FROM revisions WHERE page_id=$1",
		pageUUID)
	if err != nil {
		return nil, wikierrors.DatabaseError(err)
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
			return nil, wikierrors.DatabaseError(err)
		}
		revs[i] = *revInfo
	}

	return revs, nil
}
