package messages

type IngestionReceived struct {
	QueuedID int64
	Verbose  bool
	UserID   int32
}
