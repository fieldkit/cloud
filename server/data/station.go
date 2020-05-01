package data

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Station struct {
	ID                 int32          `db:"id,omitempty"`
	Name               string         `db:"name"`
	DeviceID           []byte         `db:"device_id"`
	OwnerID            int32          `db:"owner_id,omitempty"`
	CreatedAt          time.Time      `db:"created_at,omitempty"`
	StatusJSON         types.JSONText `db:"status_json"`
	Private            bool           `db:"private"`
	Battery            *float32       `db:"battery"`
	RecordingStartedAt *int64         `db:"recording_started_at"`
	MemoryUsed         *int64         `db:"memory_used"`
	MemoryAvailable    *int64         `db:"memory_available"`
	FirmwareNumber     *int64         `db:"firmware_number"`
	FirmwareTime       *int64         `db:"firmware_time"`
}

type StationLog struct {
	ID        int32  `db:"id,omitempty"`
	StationID int32  `db:"station_id,omitempty"`
	Timestamp string `db:"timestamp"`
	Body      string `db:"body"`
}

func (s *Station) SetStatus(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.StatusJSON = jsonData
	return nil
}

func (s *Station) GetStatus() (map[string]interface{}, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(s.StatusJSON, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (s *Station) DeviceIDHex() string {
	return hex.EncodeToString(s.DeviceID)
}

type StationFull struct {
	Station    *Station
	Owner      *User
	Ingestions []*Ingestion
	Media      []*FieldNoteMediaForStation
	Modules    []*Module
}

type Module struct {
	ID         int64  `db:"id"`
	StationID  int64  `db:"station_id"`
	HardwareID string `db:"hardware_id"`
	Name       string `db:"name"`
}

type Sensor struct {
	ID       int64  `db:"id"`
	ModuleID int64  `db:"module_id"`
	Name     string `db:"name"`
	Units    string `db:"units"`
}
