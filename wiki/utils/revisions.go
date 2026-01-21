package utils

import (
	"time"

	"github.com/google/uuid"
)

type Revision struct {
	UUID			uuid.UUID	`json:"uuid"`
	PageId			uuid.UUID	`json:"page_id"`
	Name			string		`json:"name"`
	RevDateTime		time.Time	`json:"rev_date_time"`
	Author			string		`json:"author"`
	Content			string		`json:"content"`
}

type RevisionRequest struct {
	PageId			string		`json:"page_id"`
	Author			string		`json:"author"`
	NewPage			string		`json:"new_page"`
}

