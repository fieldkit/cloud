package messages

import "time"

type IngestionReceived struct {
	QueuedID int64 `json:"queued_id"`
	Verbose  bool  `json:"verbose"`
	Refresh  bool  `json:"refresh"`
	UserID   int32 `json:"user_id"`
	SaveData bool  `json:"save_data"` // TODO Remove this
}

type IngestionCompleted struct { // Event
	QueuedID    int64     `json:"queued_id"`
	CompletedAt time.Time `json:"completed_at"`
	StationID   *int32    `json:"station_id"`
	UserID      int32     `json:"user_id"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}
