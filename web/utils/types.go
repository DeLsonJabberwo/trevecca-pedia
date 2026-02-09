package utils

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	UUID			uuid.UUID	`json:"uuid"`
	Slug			string		`json:"slug"`
	Name			string		`json:"name"`
	ArchiveDate		*time.Time	`json:"archive_date"`
	LastEditUUID	uuid.UUID	`json:"last_edit"`
	LastEditTime	time.Time	`json:"last_edit_time"`
	Content			string		`json:"content"`
}

type PageInfoPrev struct {
	UUID			uuid.UUID	`json:"uuid"`
	Slug			string		`json:"slug"`
	Name			string		`json:"name"`
	ArchiveDate		*time.Time	`json:"archive_date"`
	LastEditUUID	*uuid.UUID	`json:"last_revision_id"`
	LastEditTime	time.Time	`json:"last_edit_time"`
	Preview			string		`json:"preview"`
}
