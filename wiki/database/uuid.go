package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func GetUUID(ctx context.Context, db *sql.DB, id string) (uuid.UUID, error) {
	var pageUUID uuid.UUID
	if err := uuid.Validate(id); err == nil {
		pageUUID, err = uuid.Parse(id)
		if err != nil {
			return uuid.UUID{}, err
		}
	} else {
		pageUUID, err = getUUIDFromSlug(ctx, db, id)
		if err != nil {
			return uuid.UUID{}, err
		}
	}
	return pageUUID, nil
}

func getUUIDFromSlug(ctx context.Context, db *sql.DB, slug string) (uuid.UUID, error) {
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
