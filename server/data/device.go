package data

import (
	"encoding/json"
	"time"

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

type Firmware struct {
	ID     int64          `db:"id"`
	Time   time.Time      `db:"time"`
	Module string         `db:"module"`
	URL    string         `db:"url"`
	ETag   string         `db:"etag"`
	Meta   types.JSONText `db:"meta"`
}

func (fw *Firmware) Profile() string {
	metaMap := make(map[string]string)
	err := json.Unmarshal(fw.Meta, &metaMap)
	if err == nil {
		profile := metaMap["Build-Profile"]
		if profile != "" {
			return profile
		}
	}
	return "<unknown>"
}

type DeviceFirmware struct {
	ID       int64     `db:"id"`
	DeviceID int64     `db:"device_id"`
	Time     time.Time `db:"time"`
	Module   string    `db:"module"`
	URL      string    `db:"url"`
	ETag     string    `db:"etag"`
}
