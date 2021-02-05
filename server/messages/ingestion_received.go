package messages

type IngestionReceived struct {
	QueuedID int64 `json:"queued_id"`
	Verbose  bool  `json:"verbose"`
	UserID   int32 `json:"user_id"`
}
