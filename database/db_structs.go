package database

import (
	"time"

	"github.com/google/uuid"
)

type PageInfo struct {
	UUID			uuid.UUID	`db:"uuid"`
	Slug			string		`db:"slug"`
	Name			string		`db:"name"`
	LastRevisionId	*uuid.UUID	`db:"last_revision_id"`
	ArchiveDate		*time.Time	`db:"archive_date"`
	DeletedAt		*time.Time	`db:"deleted_at"`
}

type NameUUID struct {
	Name	string		`db:"name"`
	UUID	uuid.UUID	`db:"uuid"`
}

type RevInfo struct {
	UUID		uuid.UUID	`db:"uuid"`
	DateTime	time.Time	`db:"date_time"`
	Author		string		`db:"author"`
}
