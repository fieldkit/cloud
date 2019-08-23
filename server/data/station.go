package data

import (
	"time"
)

type Station struct {
	ID        int32     `db:"id,omitempty"`
	Name      string    `db:"name"`
	DeviceID  []byte    `db:"device_id"`
	OwnerID   int32     `db:"owner_id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}

type StationLog struct {
	ID        int32  `db:"id,omitempty"`
	StationID int32  `db:"station_id,omitempty"`
	Body      string `db:"body"`
	Timestamp string `db:"timestamp"`
}
