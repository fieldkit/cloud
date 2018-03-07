package backend

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	pb "github.com/fieldkit/data-protocol"
	"github.com/google/uuid"
	"log"
	"time"
)

type FkBinaryReader struct {
	DeviceId        string
	Location        *pb.DeviceLocation
	Time            int64
	NumberOfSensors uint32
	ReadingsSeen    uint32
	Modules         []string
	Sensors         map[uint32]*pb.SensorInfo
	Readings        map[uint32]float32
	Ingester        *ingestion.MessageIngester
	DocumentAdder   *DocumentAdder
}

func NewFkBinaryReader(b *Backend) *FkBinaryReader {
	sr := NewDatabaseSchemas(b.db)
	streams := NewDatabaseStreams(b.db)
	ingester := ingestion.NewMessageIngester(sr, streams)

	return &FkBinaryReader{
		Ingester:      ingester,
		Sensors:       make(map[uint32]*pb.SensorInfo),
		Readings:      make(map[uint32]float32),
		Modules:       make([]string, 0),
		DocumentAdder: NewDocumentAdder(b),
	}
}

type HashedData struct {
	DeviceId string
	Stream   string
	Time     int64
	Values   map[string]string
	Location []float64
}

func (br *FkBinaryReader) LocationArray() []float64 {
	if br.Location != nil {
		return []float64{float64(br.Location.Longitude), float64(br.Location.Latitude), float64(br.Location.Altitude)}
	} else {
		return []float64{}
	}
}

func (br *FkBinaryReader) ToHashingData(values map[string]string) (data []byte, err error) {
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
	values := make(map[string]string)
	for key, value := range br.Readings {
		values[br.Sensors[key].Name] = fmt.Sprintf("%f", value)
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

var (
	FkBinaryMessagesSpace = uuid.Must(uuid.Parse("0b8a5016-7410-4a1a-a2ed-2c48fec6903d"))
)

func (br *FkBinaryReader) Done() error {
	if br.ReadingsSeen > 0 {
		log.Printf("Ignored: Partial record (%v readings seen)", br.ReadingsSeen)
	}
	return nil
}

func (br *FkBinaryReader) Push(record *pb.DataRecord) error {
	if record.Metadata != nil {
		if br.DeviceId == "" {
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
					log.Printf("[Error] %v", err)
					return err
				}

				log.Printf("(%s)(%s)[Ingesting] %v, %v, %d values", pm.MessageId, pm.SchemaId, br.Modules, pm.Location, len(pm.MapValues))

				im, err := br.Ingester.IngestProcessedMessage(pm)
				if err != nil {
					log.Printf("(%s)(%s)[Error] %v", pm.MessageId, pm.SchemaId, err)
					return err
				}

				err = br.DocumentAdder.AddDocument(im)
				if err != nil {
					log.Printf("(%s)(%s)[Error] %v", pm.MessageId, pm.SchemaId, err)
					return err
				}

				log.Printf("(%s)(%s)[Success]", pm.MessageId, pm.SchemaId)
			}
		}
	}

	return nil
}
