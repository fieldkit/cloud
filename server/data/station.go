package data

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/app-protocol"
)

type Station struct {
	ID                 int32          `db:"id,omitempty"`
	Name               string         `db:"name"`
	DeviceID           []byte         `db:"device_id"`
	OwnerID            int32          `db:"owner_id,omitempty"`
	CreatedAt          time.Time      `db:"created_at,omitempty"`
	UpdatedAt          time.Time      `db:"updated_at,omitempty"`
	StatusJSON         types.JSONText `db:"status_json"`
	Private            bool           `db:"private"`
	Battery            *float32       `db:"battery"`
	Location           *Location      `db:"location"`
	LocationName       *string        `db:"location_name"`
	RecordingStartedAt *int64         `db:"recording_started_at"`
	MemoryUsed         *int32         `db:"memory_used"`
	MemoryAvailable    *int32         `db:"memory_available"`
	FirmwareNumber     *int32         `db:"firmware_number"`
	FirmwareTime       *int64         `db:"firmware_time"`
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

func (s *Station) ParseHttpReply(raw string) (*pb.HttpReply, error) {
	bytes, err := DecodeBinaryString(raw)
	if err != nil {
		return nil, err
	}

	buffer := proto.NewBuffer(bytes)
	_, err = buffer.DecodeVarint()
	if err != nil {
		return nil, err
	}

	record := &pb.HttpReply{}
	if err := buffer.Unmarshal(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Station) UpdateFromStatus(raw string) error {
	record, err := s.ParseHttpReply(raw)
	if err != nil {
		return err
	}

	if record.Status != nil {
		status := record.Status

		if status.Power != nil && status.Power.Battery != nil {
			battery := float32(status.Power.Battery.Percentage)
			s.Battery = &battery
		}

		if status.Memory != nil {
			memoryUsed := int32(status.Memory.DataMemoryUsed)
			memoryInstalled := int32(status.Memory.DataMemoryInstalled)
			s.MemoryUsed = &memoryUsed
			s.MemoryAvailable = &memoryInstalled
		}

		if status.Gps != nil {
			gps := status.Gps
			if gps.Fix > 0 {
				s.Location = NewLocation([]float64{
					float64(gps.Longitude),
					float64(gps.Latitude),
				})
			}
		}

		if status.Recording != nil {
			recordingStartedAt := int64(status.Recording.StartedTime)
			s.RecordingStartedAt = &recordingStartedAt
		}

		if status.Firmware != nil {
			firmwareTime := int64(status.Firmware.Timestamp)
			s.FirmwareTime = &firmwareTime

			if number, err := strconv.Atoi(status.Firmware.Number); err == nil {
				firmwareNumber := int32(number)
				s.FirmwareNumber = &firmwareNumber
			}
		}
	}

	return nil
}

func (s *Station) DeviceIDHex() string {
	return hex.EncodeToString(s.DeviceID)
}

type StationFull struct {
	Station    *Station
	Owner      *User
	Ingestions []*Ingestion
	Media      []*MediaForStation
	Modules    []*StationModule
	Sensors    []*ModuleSensor
}
