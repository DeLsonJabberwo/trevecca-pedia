package utils

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"wiki/database"

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
		pageUUID, err = database.GetPageUUID(ctx, db, id)
		if err != nil {
			return uuid.UUID{}, errors.New(strconv.Itoa(http.StatusNotFound))
		}
	}
	return pageUUID, nil
}
