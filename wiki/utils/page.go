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
	DeletedAt		*time.Time	`json:"deleted_at"`
	LastEdit		*uuid.UUID	`json:"last_edit"`
	LastEditTime	*time.Time	`json:"last_edit_time"`
	Content			string		`json:"content"`
}

type PageInfoPrev struct {
	UUID			uuid.UUID	`json:"uuid"`
	Slug			string		`json:"slug"`
	Name			string		`json:"name"`
	LastEditTime	time.Time	`json:"last_edit_time"`
	ArchiveDate		time.Time	`json:"archive_date"`
	Preview			string		`json:"preview"`
}

type IndexInfo struct {
	Slug			string		`json:"slug"`
	Name			string		`json:"name"`
	LastModified	time.Time	`json:"last_modified"`
	ArchiveDate		time.Time	`json:"archive_date"`
	Content			string		`json:"content"`
}

type NewPageRequest struct {
	Slug			string		`json:"slug"`
	Name			string		`json:"name"`
	Author			string		`json:"author"`
	ArchiveDate		*time.Time	`json:"archive_date"`
	Content			string		`json:"content"`
}

type DeletePageRequest struct {
	Slug			string		`json:"slug"`
	User			string		`json:"user"`
}

