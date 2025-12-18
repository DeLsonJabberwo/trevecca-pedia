package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func GetPageInfo(ctx context.Context, db *sql.DB, uuid uuid.UUID) (*PageInfo, error) {
	var p PageInfo
	err := db.QueryRowContext(
			ctx,
			"SELECT * FROM pages WHERE uuid=$1", uuid.String()).
			Scan(&p.UUID, &p.Name, &p.LastRevisionId, &p.ArchiveDate, &p.DeletedAt)
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

