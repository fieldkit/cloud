package messages

type IngestionReceived struct {
	QueuedID int64 `json:"queued_id"`
	Verbose  bool  `json:"verbose"`
	Refresh  bool  `json:"refresh"`
	UserID   int32 `json:"user_id"`
	SaveData bool  `json:"save_data"` // TODO Remove this
}

type IngestionStarted struct { // Event
	QueuedID    int64 `json:"queued_id"`
	IngestionID int64 `json:"ingestion_id"`
}

type IngestionFailed struct { // Event
	QueuedID    int64 `json:"queued_id"`
	IngestionID int64 `json:"ingestion_id"`
}

type IngestionCompleted struct { // Event
	QueuedID    int64 `json:"queued_id"`
	IngestionID int64 `json:"ingestion_id"`
}
