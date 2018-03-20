package ingestion

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type IngestedMessage struct {
	Schema   *MessageSchema
	Time     *time.Time
	Location *Location
	Fields   map[string]interface{}
	Fixed    bool
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

func (i *MessageIngester) ApplySchema(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage, schema *MessageSchema, jsonSchema *JsonMessageSchema) (im *IngestedMessage, err error) {
	// This works with MapValues, too as they'll be zero for now.
	if pm.ArrayValues != nil && len(pm.ArrayValues) != len(jsonSchema.Fields) {
		return nil, fmt.Errorf("%s: fields S=%v != M=%v", pm.SchemaId, len(jsonSchema.Fields), len(pm.ArrayValues))
	}

	mapped := make(map[string]interface{})

	for i, field := range jsonSchema.Fields {
		mapped[ToSnake(field.Name)] = pm.ArrayValues[i]
	}

	if pm.MapValues != nil {
		for k, v := range pm.MapValues {
			mapped[k] = v
		}
	}

	// TODO: Eventually we'll be able to link a secondary 'Location' stream to any other stream.
	stream, err := cache.LookupStream(pm.SchemaId.Device, func(id DeviceId) (*Stream, error) {
		return i.Streams.LookupStream(ctx, id)
	})
	if err != nil {
		return nil, err
	}

	if stream == nil {
		return nil, fmt.Errorf("%s: Unable to find stream", pm.SchemaId.Device)
	}

	time, err := determineTime(pm, jsonSchema, mapped)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get time (%s)", pm.SchemaId, err)
	}

	location, err := determineLocation(pm, jsonSchema, mapped, time)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get location (%s)", pm.SchemaId, err)
	}
	if location != nil {
		if location.Valid() {
			log.Printf("(%s)(%s)[Updating Location] %v", pm.MessageId, pm.SchemaId, location)
			i.Streams.UpdateLocation(ctx, pm.SchemaId.Device, stream, location)
			stream.Location = location
		}
	}
	if stream.Location == nil {
		return nil, fmt.Errorf("%s: Stream has no location, yet.", pm.SchemaId.Device)
	}

	im = &IngestedMessage{
		Schema:   schema,
		Time:     time,
		Location: stream.Location,
		Fixed:    pm.Fixed,
		Fields:   mapped,
	}

	return
}

func (i *MessageIngester) ApplySchemas(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage, schemas []*MessageSchema) (im *IngestedMessage, err error) {
	errors := make([]error, 0)

	if len(schemas) == 0 {
		errors = append(errors, fmt.Errorf("No matching schemas"))
	}

	for _, schema := range schemas {
		jsonSchema, ok := schema.Schema.(*JsonMessageSchema)
		if ok {
			im, applyErr := i.ApplySchema(ctx, cache, pm, schema, jsonSchema)
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

func (i *MessageIngester) IngestProcessedMessage(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage) (im *IngestedMessage, err error) {
	schemas, err := cache.LookupSchemas(pm.SchemaId, func(id SchemaId) ([]*MessageSchema, error) {
		return i.Schemas.LookupSchemas(ctx, pm.SchemaId)
	})
	if err != nil {
		return nil, fmt.Errorf("(LookupSchema) %v", err)
	}

	im, err = i.ApplySchemas(ctx, cache, pm, schemas)
	if err != nil {
		return nil, fmt.Errorf("(ApplySchemas) %v", err)
	}

	i.Statistics.Successes += 1
	return im, nil
}

func (i *MessageIngester) IngestRawMessage(ctx context.Context, cache *IngestionCache, raw *RawMessage) (im *IngestedMessage, pm *ProcessedMessage, err error) {
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
		im, err := i.IngestProcessedMessage(ctx, cache, pm)
		if err != nil {
			return nil, pm, err
		}
		return im, pm, nil
	}

	return
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
