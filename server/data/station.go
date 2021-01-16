package data

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"

	pb "github.com/fieldkit/app-protocol"
)

type StationArea struct {
	ID       int32        `db:"id,omitempty"`
	Name     string       `db:"name"`
	Geometry MultiPolygon `db:"geometry"`
}

type MultiPolygon struct {
	g *orb.MultiPolygon
}

func (l *MultiPolygon) Coordinates() [][][]float64 {
	coordinates := make([][][]float64, 0)
	for _, polygon := range *l.g {
		newPoly := make([][]float64, 0)
		for _, ring := range polygon {
			for _, c := range ring {
				newPoly = append(newPoly, []float64{c[0], c[1]})
			}
		}
		coordinates = append(coordinates, newPoly)
	}
	return coordinates
}

func (l *MultiPolygon) Scan(data interface{}) error {
	mp := &orb.MultiPolygon{}

	s := wkb.Scanner(mp)
	if err := s.Scan(data); err != nil {
		return err
	}

	l.g = mp

	return nil
}

type Station struct {
	ID                 int32      `db:"id,omitempty"`
	Name               string     `db:"name"`
	DeviceID           []byte     `db:"device_id"`
	OwnerID            int32      `db:"owner_id,omitempty"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at,omitempty"`
	Battery            *float32   `db:"battery"`
	Location           *Location  `db:"location"`
	LocationName       *string    `db:"location_name"`
	RecordingStartedAt *time.Time `db:"recording_started_at"`
	MemoryUsed         *int32     `db:"memory_used"`
	MemoryAvailable    *int32     `db:"memory_available"`
	FirmwareNumber     *int32     `db:"firmware_number"`
	FirmwareTime       *int64     `db:"firmware_time"`
	PlaceOther         *string    `db:"place_other"`
	PlaceNative        *string    `db:"place_native"`
	PhotoID            *int64     `db:"photo_id"`
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

func (s *Station) UpdateFromStatus(ctx context.Context, raw string) error {
	record, err := s.ParseHttpReply(raw)
	if err != nil {
		return err
	}

	if record.Status != nil {
		log := Logger(ctx).Sugar()

		status := record.Status

		if status.Identity != nil {
			log.Infow("identity", "identity",
				struct {
					Name         string `json:"name"`
					DeviceID     string `json:"device_id"`
					GenerationID string `json:"generation"`
					Firmware     string `json:"firmware"`
					Build        string `json:"build"`
				}{
					Name:         status.Identity.Name,
					DeviceID:     hex.EncodeToString(status.Identity.DeviceId),
					GenerationID: hex.EncodeToString(status.Identity.Generation),
					Firmware:     status.Identity.Firmware,
					Build:        status.Identity.Build,
				})

			if status.Identity.Name != "" {
				s.Name = status.Identity.Name
			}
			if status.Identity.Device != "" {
				s.Name = status.Identity.Device
			}
		}

		if status.Power != nil && status.Power.Battery != nil {
			log.Infow("battery", "device_id", s.DeviceID, "station_id", s.ID, "battery", status.Power)

			battery := float32(status.Power.Battery.Percentage)
			s.Battery = &battery
		}

		if status.Memory != nil {
			log.Infow("memory", "device_id", s.DeviceID, "station_id", s.ID, "memory", status.Memory)

			memoryUsed := int32(status.Memory.DataMemoryUsed)
			memoryInstalled := int32(status.Memory.DataMemoryInstalled)
			s.MemoryUsed = &memoryUsed
			s.MemoryAvailable = &memoryInstalled
		}

		if status.Gps != nil {
			gps := status.Gps

			log.Infow("gps", "device_id", s.DeviceID, "station_id", s.ID, "gps", gps)

			if gps.Fix > 0 {
				log.Infow("status: has gps fix", "station_id", s.ID, "gps", gps)
				s.Location = NewLocation([]float64{
					float64(gps.Longitude),
					float64(gps.Latitude),
				})
			} else if gps.Time > 0 && s.Location == nil {
				log.Infow("status: has old fix", "station_id", s.ID, "gps", gps)
				s.Location = NewLocation([]float64{
					float64(gps.Longitude),
					float64(gps.Latitude),
				})
			}
		} else {
			log.Infow("status: no gps in status", "station_id", s.ID)
		}

		if status.Recording != nil {
			log.Infow("recording", "device_id", s.DeviceID, "station_id", s.ID, "recording", status.Recording)

			if status.Recording.StartedTime > 0 {
				recordingStartedAt := time.Unix(int64(status.Recording.StartedTime), 0)
				s.RecordingStartedAt = &recordingStartedAt
			}
		}

		if status.Firmware != nil {
			log.Infow("firmware", "device_id", s.DeviceID, "station_id", s.ID, "firmware", status.Firmware)

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

type AggregatedDataSummary struct {
	StationID     int32      `db:"station_id"`
	Start         *time.Time `db:"start"`
	End           *time.Time `db:"end"`
	NumberSamples *int64     `db:"number_samples"`
}

type StationFull struct {
	Station        *Station
	Owner          *User
	Areas          []*StationArea
	Ingestions     []*Ingestion
	Configurations []*StationConfiguration
	Modules        []*StationModule
	Sensors        []*ModuleSensor
	DataSummary    *AggregatedDataSummary
	HasImages      bool
}

type EssentialStation struct {
	ID                 int64      `db:"id"`
	Name               string     `db:"name"`
	DeviceID           []byte     `db:"device_id"`
	OwnerID            int32      `db:"owner_id"`
	OwnerName          string     `db:"owner_name"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at,omitempty"`
	RecordingStartedAt *time.Time `db:"recording_started_at"`
	MemoryUsed         *int32     `db:"memory_used"`
	MemoryAvailable    *int32     `db:"memory_available"`
	FirmwareNumber     *int32     `db:"firmware_number"`
	FirmwareTime       *int64     `db:"firmware_time"`
	Location           *Location  `db:"location"`
	LastIngestionAt    *time.Time `db:"last_ingestion_at,omitempty"`
}
