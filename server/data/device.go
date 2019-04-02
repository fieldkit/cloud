package data

import (
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx/types"
)

type Device struct {
	SourceID int64  `db:"source_id"`
	Key      string `db:"key"`
	Token    string `db:"token"`
}

type DeviceSchema struct {
	ID       int64  `db:"id"`
	DeviceID int64  `db:"device_id"`
	SchemaID int64  `db:"schema_id"`
	Key      string `db:"key"`
}

type DeviceSource struct {
	Source
	Device
}

type DeviceJSONSchema struct {
	DeviceSchema
	RawSchema
}

type DeviceLocation struct {
	ID        int64      `db:"id"`
	DeviceID  int64      `db:"device_id"`
	Timestamp *time.Time `db:"timestamp"`
	Location  *Location  `db:"location"`
}

type DeviceStreamLocation struct {
	ID        int64     `db:"id"`
	DeviceID  string    `db:"device_id"`
	Timestamp time.Time `db:"timestamp"`
	Location  *Location `db:"location"`
}

type DeviceStreamLocationAndPlace struct {
	DeviceStreamLocation
	Places pq.StringArray `db:"places"`
}

type Firmware struct {
	ID      int64          `db:"id"`
	Time    time.Time      `db:"time"`
	Module  string         `db:"module"`
	Profile string         `db:"profile"`
	URL     string         `db:"url"`
	ETag    string         `db:"etag"`
	Meta    types.JSONText `db:"meta"`
}

type DeviceFirmware struct {
	ID       int64     `db:"id"`
	DeviceID int64     `db:"device_id"`
	Time     time.Time `db:"time"`
	Module   string    `db:"module"`
	Profile  string    `db:"profile"`
	URL      string    `db:"url"`
	ETag     string    `db:"etag"`
}

type DeviceFile struct {
	ID       int64          `db:"id"`
	Time     time.Time      `db:"time"`
	StreamID string         `db:"stream_id"`
	Firmware string         `db:"firmware"`
	DeviceID string         `db:"device_id"`
	Size     int64          `db:"size"`
	FileID   string         `db:"file_id"`
	URL      string         `db:"url"`
	Meta     types.JSONText `db:"meta"`
	Children *[]int64       `db:"children"`
}

type DeviceStream struct {
	ID       int64          `db:"id"`
	Time     time.Time      `db:"time"`
	StreamID string         `db:"stream_id"`
	Firmware string         `db:"firmware"`
	DeviceID string         `db:"device_id"`
	Size     int64          `db:"size"`
	FileID   string         `db:"file_id"`
	URL      string         `db:"url"`
	Meta     types.JSONText `db:"meta"`
	Children *[]int64       `db:"children"`
}

type DeviceSummary struct {
	DeviceID      string    `db:"device_id"`
	LastFileID    string    `db:"last_stream_id"`
	LastFileTime  time.Time `db:"last_stream_time"`
	NumberOfFiles int       `db:"number_of_files"`
	LogsSize      int       `db:"logs_size"`
	DataSize      int       `db:"data_size"`
}
