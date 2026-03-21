package utils

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	UUID         uuid.UUID  `json:"uuid"`
	Slug         string     `json:"slug"`
	Name         string     `json:"name"`
	ArchiveDate  *time.Time `json:"archive_date"`
	LastEditUUID *uuid.UUID `json:"last_edit"`
	LastEditTime time.Time  `json:"last_edit_time"`
	Content      string     `json:"content"`
	Categories   []Category `json:"categories"`
}

type PageInfoPrev struct {
	UUID         uuid.UUID  `json:"uuid"`
	Slug         string     `json:"slug"`
	Name         string     `json:"name"`
	LastEditTime time.Time  `json:"last_edit_time"`
	ArchiveDate  *time.Time `json:"archive_date"`
	Preview      string     `json:"preview"`
}

type Category struct {
	ID       int        `json:"id"`
	Slug     string     `json:"slug"`
	Name     string     `json:"name"`
	FullSlug string     `json:"full_slug"`
	Children []Category `json:"children,omitempty"`
}

type CategoryFlat struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	FullSlug    string `json:"full_slug"`
	Depth       int    `json:"depth"`
	DisplayName string `json:"display_name"`
}

// Revision represents a page revision from the API
// Note: The API uses different field names for list vs detail endpoints
type Revision struct {
	UUID        uuid.UUID  `json:"uuid"`
	PageId      uuid.UUID  `json:"page_id"`
	RevDateTime time.Time  `json:"rev_date_time"`
	Author      string     `json:"author"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	ArchiveDate *time.Time `json:"archive_date"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Content     string     `json:"content"`
}
