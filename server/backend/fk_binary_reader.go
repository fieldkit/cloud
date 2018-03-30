package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type FormattedMessageReceiver interface {
	HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) error
}

type FkBinaryReader struct {
	RecordsProcessed uint32

	DeviceId        string
	Location        *pb.DeviceLocation
	Time            int64
	NumberOfSensors uint32
	ReadingsSeen    uint32

	Modules  []string
	Sensors  map[uint32]*pb.SensorInfo
	Readings map[uint32]float32

	Receiver FormattedMessageReceiver
}

func NewFkBinaryReader(receiver FormattedMessageReceiver) *FkBinaryReader {
	return &FkBinaryReader{
		Modules:  make([]string, 0),
		Sensors:  make(map[uint32]*pb.SensorInfo),
		Readings: make(map[uint32]float32),
		Receiver: receiver,
	}
}

func (br *FkBinaryReader) Read(ctx context.Context, body io.Reader) error {
	unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var record pb.DataRecord
		err := proto.Unmarshal(b, &record)
		if err != nil {
			// We keep reading, this may just be a protocol version issue.
			log.Printf("Error unmarshalling record: %v", err)
			return nil, nil
		}

		err = br.Push(ctx, &record)
		if err, ok := err.(*ingestion.IngestError); ok {
			if err.Critical {
				return nil, err
			} else {
				log.Printf("Error: %v", err)
			}
		} else if err != nil {
			return nil, err
		}

		return &record, nil
	})

	_, err := stream.ReadLengthPrefixedCollection(body, unmarshalFunc)
	if err != nil {
		return err
	}

	if br.ReadingsSeen > 0 {
		log.Printf("Ignored: partial record (%v readings seen)", br.ReadingsSeen)
	}

	log.Printf("Processed %d records", br.RecordsProcessed)

	return nil
}

func (br *FkBinaryReader) Push(ctx context.Context, record *pb.DataRecord) error {
	br.RecordsProcessed += 1

	if record.Metadata != nil {
		if record.Metadata.DeviceId != nil && len(record.Metadata.DeviceId) > 0 {
			br.DeviceId = hex.EncodeToString(record.Metadata.DeviceId)
		}
		if record.Metadata.Sensors != nil {
			if br.NumberOfSensors == 0 {
				for _, sensor := range record.Metadata.Sensors {
					br.Sensors[sensor.Sensor] = sensor
					br.NumberOfSensors += 1
				}
			}
		}
		if record.Metadata.Modules != nil {
			br.Modules = make([]string, 0)
			for _, m := range record.Metadata.Modules {
				br.Modules = append(br.Modules, m.Name)
			}
		}
	}
	if record.LoggedReading != nil {
		reading := record.LoggedReading.Reading
		location := record.LoggedReading.Location

		if location != nil {
			br.Location = location
		}

		if reading != nil {
			if br.NumberOfSensors == 0 {
				log.Printf("Ignored: Unknown sensor. (%+v)", record)
				return nil
			}

			if reading.Sensor == 0 {
				br.ReadingsSeen = 0
			}

			br.Readings[reading.Sensor] = reading.Value
			br.ReadingsSeen += 1

			if br.ReadingsSeen == br.NumberOfSensors {
				br.Time = int64(record.LoggedReading.Reading.Time)
				br.ReadingsSeen = 0

				fm, err := br.getFormattedMessage()
				if err != nil {
					return fmt.Errorf("Unable to create formatted message (%v)", err)
				}

				if err := br.Receiver.HandleFormattedMessage(ctx, fm); err != nil {
					return err
				}
			}
		}
	}

	if record.Log != nil {
		if false {
			fmt.Printf("%-10d %-30s %s\n", record.Log.Uptime, record.Log.Facility, record.Log.Message)
		}
	}

	return nil
}

func (br *FkBinaryReader) getLocationArray() []float64 {
	if br.Location != nil {
		return []float64{float64(br.Location.Longitude), float64(br.Location.Latitude), float64(br.Location.Altitude)}
	} else {
		return []float64{}
	}
}

func (br *FkBinaryReader) getFormattedMessage() (fm *ingestion.FormattedMessage, err error) {
	values := make(map[string]interface{})
	for key, value := range br.Readings {
		if math.IsNaN(float64(value)) {
			values[br.Sensors[key].Name] = "NaN"
		} else if math.IsInf(float64(value), 0) {
			values[br.Sensors[key].Name] = "NaN"
		} else {
			values[br.Sensors[key].Name] = value
		}
	}

	messageTime := time.Unix(br.Time, 0)
	location := br.getLocationArray()

	messageId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	fm = &ingestion.FormattedMessage{
		MessageId: ingestion.MessageId(messageId.String()),
		SchemaId:  ingestion.NewSchemaId(ingestion.NewDeviceId(br.DeviceId), ""),
		Time:      &messageTime,
		Location:  location,
		Fixed:     len(location) > 0,
		MapValues: values,
		Modules:   br.Modules,
	}

	return
}
