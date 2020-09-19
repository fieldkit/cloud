package data

import (
	"time"
)

type Collection struct {
	ID          int32     `db:"id,omitempty"`
	OwnerID     int32     `db:"owner_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Tags        string    `db:"tags"`
	Private     bool      `db:"private"`
	Origin      string    `db:"origin"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CollectionUser struct {
	UserID    int32 `db:"user_id"`
	ProjectID int32 `db:"project_id"`
}
