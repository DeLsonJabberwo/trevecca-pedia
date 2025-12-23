package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func GetPageUUID(ctx context.Context, db *sql.DB, slug string) (uuid.UUID, error) {
	var pageId uuid.UUID
	err := db.QueryRowContext(
			ctx,
			"SELECT uuid FROM pages WHERE slug=$1", slug).
			Scan(&pageId)
	if err != nil {
		return uuid.Nil, err
	}
	return pageId, nil
}

func GetPageInfo(ctx context.Context, db *sql.DB, uuid uuid.UUID) (*PageInfo, error) {
	var p PageInfo
	err := db.QueryRowContext(
			ctx,
			"SELECT * FROM pages WHERE uuid=$1", uuid.String()).
			Scan(&p.UUID, &p.Slug, &p.Name, &p.LastRevisionId, &p.ArchiveDate, &p.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func GetPageNameUUIDs(ctx context.Context, db *sql.DB) ([]NameUUID, error) {
	var r []NameUUID
	rows, err := db.QueryContext(
				ctx,
				"SELECT name, uuid FROM pages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row NameUUID
		err := rows.Scan(&row.Name, &row.UUID)
		if err != nil {
			return nil, err
		}
		r = append(r, row)
	}
	return r, nil
}

func GetPageRevisionsInfo(ctx context.Context, db *sql.DB, pageId uuid.UUID) ([]RevInfo, error) {
	var revs []RevInfo
	rows, err := db.QueryContext(
				ctx,
				"SELECT uuid, date_time, author FROM revisions WHERE page_id=$1 ORDER BY date_time",
				pageId.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row RevInfo
		err := rows.Scan(&row.UUID, &row.DateTime, &row.Author)
		if err != nil {
			return nil, err
		}
		revs = append(revs, row)
	}

	return revs, nil
}

func GetRevisionInfo(ctx context.Context, db *sql.DB, revId uuid.UUID) *RevInfo {
	var rev RevInfo
	db.QueryRowContext(
				ctx,
				"SELECT uuid, page_id, date_time, author FROM revisions WHERE uuid=$1",
				revId).Scan(&rev.UUID, &rev.PageId, &rev.DateTime, &rev.Author)
	return &rev
}

func GetMostRecentSnapshot(ctx context.Context, db *sql.DB, revId uuid.UUID) SnapInfo {
	revInfo := GetRevisionInfo(ctx, db, revId)
	pageId := revInfo.PageId
	var snapId uuid.UUID
	var snap SnapInfo
	var snapCount int
	db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM snapshots
		WHERE page=$1`,
		pageId).
		Scan(&snapCount)
	if snapCount == 1 {
		db.QueryRowContext(
			ctx,
			`SELECT uuid FROM snapshots
			WHERE page=$1`,
			pageId).
			Scan(&snapId)
	} else {
		db.QueryRowContext(
			ctx,
			`SELECT uuid FROM snapshots
			JOIN revisions ON snapshots.revision = revisions.uuid
			WHERE snapshots.page=$1
			ORDER BY revisions.date_time`,
			pageId).
			Scan(&snapId)
	}
	db.QueryRowContext(
		ctx,
		"SELECT * FROM snapshots WHERE uuid=$1",
		snapId).
		Scan(&snap.UUID, &snap.Page, &snap.Revision)
	return snap
}

func GetMissingRevisions(ctx context.Context, db *sql.DB, revId uuid.UUID) ([]RevInfo, error) {
	var revs []RevInfo
	var count int
	var snapRevTime time.Time

	snap := GetMostRecentSnapshot(ctx, db, revId)
	if snap.Revision == nil {
		snapRevTime = time.Time{}
	} else {
		snap_rev := GetRevisionInfo(ctx, db, *snap.Revision)
		snapRevTime = *snap_rev.DateTime
	}
	db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM revisions WHERE date_time > $1",
		snapRevTime).
		Scan(&count)
	revs = make([]RevInfo, count)

	revIds, err := db.QueryContext(
				ctx,
				"SELECT uuid FROM revisions WHERE date_time >  $1 ORDER BY date_time ASC",
				snapRevTime)
	if err != nil {
		return nil, err
	}

	for i := 0; revIds.Next(); i++ {
		var id uuid.UUID
		revIds.Scan(&id)
		revs[i] = *GetRevisionInfo(ctx, db, id)
	}

	return revs, nil
}

