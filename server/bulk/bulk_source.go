package bulk

import (
	"context"
	"time"
)

type BulkReading struct {
	Key        string  `json:"key"`
	Calibrated float64 `json:"calibrated"`
	Battery    bool    `json:"battery"`
}

type BulkMessage struct {
	DeviceID   []byte        `json:"device_id"`
	DeviceName string        `json:"device_name"`
	Time       time.Time     `json:"time"`
	Readings   []BulkReading `json:"readings"`
}

type BulkMessageBatch struct {
	Messages []*BulkMessage `json:"messages"`
}

type BulkSource interface {
	NextBatch(ctx context.Context, batch *BulkMessageBatch) error
}
