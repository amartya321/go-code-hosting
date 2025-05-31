package model

import (
	"database/sql"
	"time"
)

type Repository struct {
	ID          int            `json:"id"`
	OwnerID     int            `json:"owner_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	IsPrivate   bool           `json:"is_private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
