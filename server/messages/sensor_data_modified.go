package messages

import (
	"time"
)

type SensorDataModified struct {
	ModifiedAt  time.Time `json:"modified_at"`
	PublishedAt time.Time `json:"published_at"`
	StationID   *int32    `json:"station_id"`
	UserID      int32     `json:"user_id"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}
