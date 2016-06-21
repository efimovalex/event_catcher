package common

import (
	"github.com/efimovalex/EventKitAPI/adaptors/database"
)

type ListResponse struct {
	Events   []database.Event `json:"events"`
	NextPage string           `json:"next_page,omitempty"`
}
