package requests

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	UUID		uuid.UUID	`json:"uuid"`
	Slug		string		`json:"slug"`
	Name		string		`json:"name"`
	ArchiveDate	*time.Time	`json:"archive_date"`
	DeletedAt	*time.Time	`json:"deleted_at"`
	Content		string		`json:"content"`
}
