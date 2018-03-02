package backend

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	pb "github.com/fieldkit/data-protocol"
	"github.com/google/uuid"
	"log"
)

type FkBinaryReader struct {
	DeviceId        string
	Location        *pb.DeviceLocation
	Time            int64
	NumberOfSensors uint32
	ReadingsSeen    uint32
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
		DocumentAdder: NewDocumentAdder(b),
	}
}

func (br *FkBinaryReader) CreateHttpJsonMessage() *ingestion.HttpJsonMessage {
	values := make(map[string]string)
	for key, value := range br.Readings {
		values[br.Sensors[key].Name] = fmt.Sprintf("%f", value)
	}

	if br.Location != nil {
		return &ingestion.HttpJsonMessage{
			Location: []float64{float64(br.Location.Longitude), float64(br.Location.Latitude), float64(br.Location.Altitude)},
			Fixed:    br.Location.Fix == 1,
			Time:     br.Time,
			Device:   br.DeviceId,
			Stream:   "",
			Values:   values,
		}
	}

	return &ingestion.HttpJsonMessage{
		Fixed:  false,
		Time:   br.Time,
		Device: br.DeviceId,
		Stream: "",
		Values: values,
	}
}

var (
	FkBinaryMessagesSpace = uuid.Must(uuid.Parse("0b8a5016-7410-4a1a-a2ed-2c48fec6903d"))
)

func ToHashingData(im *ingestion.HttpJsonMessage) (data []byte, err error) {
	data, err = json.Marshal(im)
	if err != nil {
		return nil, err
	}
	return
}

func ToUniqueHash(im *ingestion.HttpJsonMessage) (uuid.UUID, error) {
	data, err := ToHashingData(im)
	if err != nil {
		return uuid.New(), err
	}
	return uuid.NewSHA1(FkBinaryMessagesSpace, data), err
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

			br.Readings[reading.Sensor] = reading.Value
			br.ReadingsSeen += 1

			if br.ReadingsSeen == br.NumberOfSensors {
				br.Time = int64(record.LoggedReading.Reading.Time)
				br.ReadingsSeen = 0

				message := br.CreateHttpJsonMessage()
				token, err := ToUniqueHash(message)
				if err != nil {
					return err
				}

				messageId := ingestion.MessageId(token.String())

				log.Printf("(%s)[Ingesting] %+v", messageId, message)

				pm, err := message.ToProcessedMessage(messageId)
				if err != nil {
					log.Printf("(%s)[Error] %v", messageId, err)
					return err
				}

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
