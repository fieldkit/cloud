package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
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

	Ingester    *ingestion.MessageIngester
	Resolver    *ingestion.Resolver
	RecordAdder *RecordAdder
}

func NewFkBinaryReader(b *Backend) *FkBinaryReader {
	r := ingestion.NewRepository(b.db)
	ingester := ingestion.NewMessageIngester(r)

	return &FkBinaryReader{
		Modules:   make([]string, 0),
		Sensors:   make(map[uint32]*pb.SensorInfo),
		Readings:  make(map[uint32]float32),
		SourceIDs: make(map[int64]bool),

		Ingester:    ingester,
		Resolver:    ingestion.NewResolver(r),
		RecordAdder: NewRecordAdder(b),
	}
}

type HashedData struct {
	DeviceId string
	Stream   string
	Time     int64
	Values   map[string]interface{}
	Location []float64
}

func (br *FkBinaryReader) LocationArray() []float64 {
	if br.Location != nil {
		return []float64{float64(br.Location.Longitude), float64(br.Location.Latitude), float64(br.Location.Altitude)}
	} else {
		return []float64{}
	}
}

func (br *FkBinaryReader) ToHashingData(values map[string]interface{}) (data []byte, err error) {
	hashed := &HashedData{
		DeviceId: br.DeviceId,
		Stream:   "",
		Time:     br.Time,
		Values:   values,
		Location: br.LocationArray(),
	}

	data, err = json.Marshal(hashed)
	if err != nil {
		return nil, err
	}

	return
}

func ToUniqueHash(data []byte) (uuid.UUID, error) {
	return uuid.NewSHA1(FkBinaryMessagesSpace, data), nil
}

func (br *FkBinaryReader) CreateProcessedMessage() (pm *ingestion.ProcessedMessage, err error) {
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

	hd, err := br.ToHashingData(values)
	if err != nil {
		return nil, err
	}

	token, err := ToUniqueHash(hd)
	if err != nil {
		return nil, err
	}

	messageId := token.String()
	messageTime := time.Unix(br.Time, 0)

	if br.Location != nil {
		pm = &ingestion.ProcessedMessage{
			MessageId: ingestion.MessageId(messageId),
			SchemaId:  ingestion.NewSchemaId(ingestion.NewDeviceId(br.DeviceId), ""),
			Time:      &messageTime,
			Location:  br.LocationArray(),
			Fixed:     true,
			MapValues: values,
		}
	} else {
		pm = &ingestion.ProcessedMessage{
			MessageId: ingestion.MessageId(messageId),
			SchemaId:  ingestion.NewSchemaId(ingestion.NewDeviceId(br.DeviceId), ""),
			Time:      &messageTime,
			Location:  []float64{},
			Fixed:     false,
			MapValues: values,
		}
	}

	return
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

				pm, err := br.CreateProcessedMessage()
				if err != nil {
					return fmt.Errorf("Unable to create processed message (%v)", err)
				}

				ds, err := br.Resolver.ResolveDeviceAndSchemas(ctx, pm.SchemaId)
				if err != nil {
					return err
				}

				im, err := br.Ingester.IngestProcessedMessage(ctx, ds, pm)
				if err != nil {
					return err
				}

				err = br.RecordAdder.AddRecord(ctx, im)
				if err != nil {
					return err
				}

				br.SourceIDs[im.Schema.Ids.DeviceID] = true

				log.Printf("(%s)(%s)[Success] %v, %d values (location = %t), %v", pm.MessageId, pm.SchemaId, br.Modules, len(pm.MapValues), im.LocationUpdated, pm.Location)
			}
		}
	}

	return nil
}

func (br *FkBinaryReader) Done(ctx context.Context) error {
	if br.ReadingsSeen > 0 {
		log.Printf("Ignored: partial record (%v readings seen)", br.ReadingsSeen)
	}

	log.Printf("Processed %d records", br.RecordsProcessed)

	for id, _ := range br.SourceIDs {
		br.RecordAdder.EmitSourceChanged(id)
	}

	return nil
}
