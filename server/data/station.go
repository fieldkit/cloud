package data

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Station struct {
	ID         int32          `db:"id,omitempty"`
	Name       string         `db:"name"`
	DeviceID   []byte         `db:"device_id"`
	OwnerID    int32          `db:"owner_id,omitempty"`
	CreatedAt  time.Time      `db:"created_at,omitempty"`
	StatusJSON types.JSONText `db:"status_json"`
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
