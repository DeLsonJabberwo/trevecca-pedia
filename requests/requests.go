package requests

import "github.com/google/uuid"

type PageRequest struct {
	UUID	*uuid.UUID	`json:"uuid"`
	Slug	*string		`json:"slug"`
}
