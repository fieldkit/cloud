package data

import (
	"time"
)

type StationNote struct {
	ID        int32     `db:"id"`
	UserID    int32     `db:"user_id"`
	StationID *int32    `db:"station_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Body      string    `db:"body"`
}

type StationNotesWithUsers struct {
	Notes     []*StationNote
	UsersByID map[int32]*User
}
