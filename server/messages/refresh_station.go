package messages

import (
	"time"
)

type RefreshStation struct {
	StationID   int32
	Completely  bool
	Verbose     bool
	UserID      int32
	HowRecently time.Duration
}
