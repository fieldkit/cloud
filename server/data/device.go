package data

import (
	"time"
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
