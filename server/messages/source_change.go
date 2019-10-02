package messages

import (
	"time"
)

type SourceChange struct {
	SourceID    int64
	DeviceID    string
	FileTypeIDs []string
	QueuedAt    time.Time
}

func NewSourceChange(sourceId int64, deviceID string, fileTypeIDs []string) SourceChange {
	return SourceChange{
		SourceID:    sourceId,
		DeviceID:    deviceID,
		QueuedAt:    time.Now(),
		FileTypeIDs: fileTypeIDs,
	}
}
