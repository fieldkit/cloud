package messages

import (
	"time"
)

type RefreshStation struct {
	StationID   int32         `json:"station_id"`
	Completely  bool          `json:"completely"`
	SkipManual  bool          `json:"skip_manual,omitempty"`
	Verbose     bool          `json:"verbose"`
	UserID      int32         `json:"user_id"`
	HowRecently time.Duration `json:"how_recently"`
}
