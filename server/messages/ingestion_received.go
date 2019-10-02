package messages

import "time"

type IngestionReceived struct {
	Time time.Time
	ID   int64
	URL  string
}
