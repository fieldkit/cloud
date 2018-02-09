package ingestion

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type IngestedMessage struct {
	Schema   *JsonMessageSchema
	Time     *time.Time
	Location *Location
	Fields   map[string]interface{}
}

const (
	FieldNameLatitude  = "latitude"
	FieldNameLongitude = "longitude"
	FieldNameAltitude  = "altitude"
	FieldNameTime      = "time"
)

func determineTime(pm *ProcessedMessage, ms *JsonMessageSchema, m map[string]interface{}) (t *time.Time, err error) {
	if ms.UseProviderTime {
		t = pm.Time
	} else {
		if !ms.HasTime {
			return nil, fmt.Errorf("%s: No time information.", pm.SchemaId)
		}

		raw := m[FieldNameTime].(string)
		unix, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, err
		}
		parsed := time.Unix(unix, 0)
		t = &parsed
	}

	return
}

func determineLocation(pm *ProcessedMessage, ms *JsonMessageSchema, m map[string]interface{}, t *time.Time) (l *Location, err error) {
	if ms.UseProviderLocation {
		return &Location{UpdatedAt: t, Coordinates: pm.Location}, nil
	}
	if ms.HasLocation {
		coordinates := make([]float64, 0)
		for _, key := range []string{FieldNameLongitude, FieldNameLatitude, FieldNameAltitude} {
			f, err := strconv.ParseFloat(m[key].(string), 64)
			if err != nil {
				return nil, err
			}
			coordinates = append(coordinates, f)
		}
		if len(coordinates) < 2 {
			return nil, fmt.Errorf("Not enough coordinates.")
		}
		return &Location{
			UpdatedAt:   t,
			Coordinates: coordinates,
		}, nil
	}
	return nil, nil
}

func (i *MessageIngester) ApplySchema(pm *ProcessedMessage, ms *JsonMessageSchema) (im *IngestedMessage, err error) {
	// This works with MapValues, too as they'll be zero for now.
	if pm.ArrayValues != nil && len(pm.ArrayValues) != len(ms.Fields) {
		return nil, fmt.Errorf("%s: fields S=%v != M=%v", pm.SchemaId, len(ms.Fields), len(pm.ArrayValues))
	}

	mapped := make(map[string]interface{})

	for i, field := range ms.Fields {
		mapped[ToSnake(field.Name)] = pm.ArrayValues[i]
	}

	if pm.MapValues != nil {
		for k, v := range pm.MapValues {
			mapped[k] = v
		}
	}

	// TODO: Eventually we'll be able to link a secondary 'Location' stream
	// to any other stream.
	stream, err := i.Streams.LookupStream(pm.SchemaId.Device)
	if err != nil {
		return nil, err
	}

	if stream == nil {
		return nil, fmt.Errorf("%s: Unable to find stream", pm.SchemaId.Device)
	}

	time, err := determineTime(pm, ms, mapped)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get time (%s)", pm.SchemaId, err)
	}

	location, err := determineLocation(pm, ms, mapped, time)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get location (%s)", pm.SchemaId, err)
	}
	if location != nil {
		if location.Valid() {
			log.Printf("(%s)(%s)[Updating Location] %v", pm.MessageId, pm.SchemaId, location)
			i.Streams.UpdateLocation(pm.SchemaId.Device, location)
			stream.Location = location
		}
	}
	if stream.Location == nil {
		return nil, fmt.Errorf("%s: Stream has no location, yet.", pm.SchemaId.Device)
	}

	im = &IngestedMessage{
		Schema:   ms,
		Time:     time,
		Location: stream.Location,
		Fields:   mapped,
	}

	return
}

func (i *MessageIngester) ApplySchemas(pm *ProcessedMessage, schemas []interface{}) (im *IngestedMessage, err error) {
	errors := make([]error, 0)

	if len(schemas) == 0 {
		errors = append(errors, fmt.Errorf("No matching schemas"))
	}

	for _, schema := range schemas {
		jsonSchema, ok := schema.(*JsonMessageSchema)
		if ok {
			im, applyErr := i.ApplySchema(pm, jsonSchema)
			if applyErr != nil {
				errors = append(errors, applyErr)
			} else {
				return im, nil
			}
		} else {
			log.Printf("Unexpected MessageSchema type: %v", schema)
		}
	}
	return nil, fmt.Errorf("%v", makePretty(errors))
}

func makePretty(errors []error) (m string) {
	strs := []string{}
	for _, err := range errors {
		strs = append(strs, fmt.Sprintf("(%v)", err))
	}

	return strings.Join(strs, ", ")
}

func (i *MessageIngester) Ingest(raw *RawMessage) (im *IngestedMessage, pm *ProcessedMessage, err error) {
	i.Statistics.Processed += 1

	mp, err := IdentifyMessageProvider(raw)
	if err != nil {
		return nil, nil, err
	}
	if mp == nil {
		return nil, nil, fmt.Errorf("No message provider: (ContentType: %s)", raw.ContentType)
	}

	pm, err = mp.ProcessMessage(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("(ProcessMessage) %v", err)
	}
	if pm != nil {
		im, err := i.IngestProcessedMessage(pm)
		if err != nil {
			return nil, pm, err
		}
		return im, pm, nil
	}

	return
}

func (i *MessageIngester) IngestProcessedMessage(pm *ProcessedMessage) (im *IngestedMessage, err error) {
	schemas, err := i.Schemas.LookupSchema(pm.SchemaId)
	if err != nil {
		return nil, fmt.Errorf("(LookupSchema) %v", err)
	}
	im, err = i.ApplySchemas(pm, schemas)
	if err != nil {
		return nil, fmt.Errorf("(ApplySchemas) %v", err)
	}
	i.Statistics.Successes += 1
	return im, nil
}

func (i *MessageIngester) HandleMessage(raw *RawMessage) error {
	_, pm, err := i.Ingest(raw)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("(%s)(%s)[Success]", pm.MessageId, pm.SchemaId)
	}
	return nil
}

type IngestionStatistics struct {
	Processed uint64
	Successes uint64
}

type MessageIngester struct {
	RawMessageHandler
	Schemas    SchemaRepository
	Streams    StreamsRepository
	Statistics IngestionStatistics
}

func NewMessageIngester(sr SchemaRepository, streams StreamsRepository) *MessageIngester {
	return &MessageIngester{
		Schemas: sr,
		Streams: streams,
	}
}
