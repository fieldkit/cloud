package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

var (
	FkBinaryMessagesSpace = uuid.Must(uuid.Parse("0b8a5016-7410-4a1a-a2ed-2c48fec6903d"))
)

type FkBinaryReader struct {
	RecordsProcessed uint32

	DeviceId        string
	Location        *pb.DeviceLocation
	Time            int64
	NumberOfSensors uint32
	ReadingsSeen    uint32

	Modules   []string
	Sensors   map[uint32]*pb.SensorInfo
	Readings  map[uint32]float32
	SourceIDs map[int64]bool

	SchemaApplier *ingestion.SchemaApplier
	Repository    *ingestion.Repository
	Resolver      *ingestion.Resolver
	RecordAdder   *ingestion.RecordAdder
}

func NewFkBinaryReader(b *Backend) *FkBinaryReader {
	r := ingestion.NewRepository(b.db)

	return &FkBinaryReader{
		Modules:   make([]string, 0),
		Sensors:   make(map[uint32]*pb.SensorInfo),
		Readings:  make(map[uint32]float32),
		SourceIDs: make(map[int64]bool),

		Repository:    r,
		SchemaApplier: ingestion.NewSchemaApplier(),
		Resolver:      ingestion.NewResolver(r),
		RecordAdder:   ingestion.NewRecordAdder(r)}
}

func (br *FkBinaryReader) LocationArray() []float64 {
	if br.Location != nil {
		return []float64{float64(br.Location.Longitude), float64(br.Location.Latitude), float64(br.Location.Altitude)}
	} else {
		return []float64{}
	}
}

func (br *FkBinaryReader) CreateFormattedMessage() (fm *ingestion.FormattedMessage, err error) {
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
	location := br.LocationArray()

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
	}

	return
}

func (br *FkBinaryReader) Ingest(ctx context.Context) error {
	fm, err := br.CreateFormattedMessage()
	if err != nil {
		return fmt.Errorf("Unable to create processed message (%v)", err)
	}

	ds, err := br.Resolver.ResolveDeviceAndSchemas(ctx, fm.SchemaId)
	if err != nil {
		return err
	}

	pm, err := br.SchemaApplier.ApplySchemas(ds, fm)
	if err != nil {
		return err
	}

	err = br.RecordAdder.AddRecord(ctx, ds, pm)
	if err != nil {
		return err
	}

	br.SourceIDs[pm.Schema.Ids.DeviceID] = true

	log.Printf("(%s)(%s)[Success] %v, %d values (location = %t), %v", fm.MessageId, fm.SchemaId, br.Modules, len(fm.MapValues), pm.LocationUpdated, fm.Location)

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

				if err := br.Ingest(ctx); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (br *FkBinaryReader) Done(ctx context.Context, sourceChanges ingestion.SourceChangesPublisher) error {
	if br.ReadingsSeen > 0 {
		log.Printf("Ignored: partial record (%v readings seen)", br.ReadingsSeen)
	}

	log.Printf("Processed %d records", br.RecordsProcessed)

	for id, _ := range br.SourceIDs {
		sourceChanges.SourceChanged(ingestion.NewSourceChange(id))
	}

	return nil
}
