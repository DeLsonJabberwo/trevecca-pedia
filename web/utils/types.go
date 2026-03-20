package utils

import (
	"encoding/json"
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

// UnmarshalJSON implements custom JSON unmarshaling to handle both API formats
func (r *Revision) UnmarshalJSON(data []byte) error {
	// Try format 1: lowercase with underscores (individual revision endpoint)
	type RevisionFormat1 struct {
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

	var f1 RevisionFormat1
	if err := json.Unmarshal(data, &f1); err == nil && !f1.RevDateTime.IsZero() {
		r.UUID = f1.UUID
		r.PageId = f1.PageId
		r.RevDateTime = f1.RevDateTime
		r.Author = f1.Author
		r.Slug = f1.Slug
		r.Name = f1.Name
		r.ArchiveDate = f1.ArchiveDate
		r.DeletedAt = f1.DeletedAt
		r.Content = f1.Content
		return nil
	}

	// Try format 2: capitalized field names (list endpoint)
	type RevisionFormat2 struct {
		UUID        *uuid.UUID `json:"UUID"`
		PageId      *uuid.UUID `json:"PageId"`
		DateTime    *time.Time `json:"DateTime"`
		Author      *string    `json:"Author"`
		Slug        string     `json:"Slug"`
		Name        string     `json:"Name"`
		ArchiveDate *time.Time `json:"ArchiveDate"`
		DeletedAt   *time.Time `json:"DeletedAt"`
	}

	var f2 RevisionFormat2
	if err := json.Unmarshal(data, &f2); err != nil {
		return err
	}

	if f2.UUID != nil {
		r.UUID = *f2.UUID
	}
	if f2.PageId != nil {
		r.PageId = *f2.PageId
	}
	if f2.DateTime != nil {
		r.RevDateTime = *f2.DateTime
	}
	if f2.Author != nil {
		r.Author = *f2.Author
	}
	r.Slug = f2.Slug
	r.Name = f2.Name
	r.ArchiveDate = f2.ArchiveDate
	r.DeletedAt = f2.DeletedAt

	return nil
}
