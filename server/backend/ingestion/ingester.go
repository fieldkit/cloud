package ingestion

import (
	"context"
	"fmt"
	_ "log"
	"strconv"
	"strings"
	"time"
)

type IngestedMessage struct {
	Schema          *MessageSchema
	Time            *time.Time
	Location        *Location
	LocationUpdated bool
	Fixed           bool
	Fields          map[string]interface{}
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

const (
	FieldNameLatitude  = "latitude"
	FieldNameLongitude = "longitude"
	FieldNameAltitude  = "altitude"
	FieldNameTime      = "time"
)

func NewMessageIngester(sr SchemaRepository, streams StreamsRepository) *MessageIngester {
	return &MessageIngester{
		Schemas: sr,
		Streams: streams,
	}
}

func determineTime(pm *ProcessedMessage, ms *JsonMessageSchema, m map[string]interface{}) (t *time.Time, err error) {
	if ms.UseProviderTime {
		t = pm.Time
	} else {
		if !ms.HasTime {
			return nil, NewErrorf(true, "No time information")
		}

		raw := m[FieldNameTime].(string)
		unix, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, NewErrorf(true, "Unable to parse time '%s' (%v)", raw, err)
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

	if !ms.HasLocation {
		return nil, nil
	}

	coordinates := make([]float64, 0)
	for _, key := range []string{FieldNameLongitude, FieldNameLatitude, FieldNameAltitude} {
		f, err := strconv.ParseFloat(m[key].(string), 64)
		if err != nil {
			return nil, NewErrorf(true, "Unable to parse coordinate '%s' (%v)", m[key], err)
		}
		coordinates = append(coordinates, f)
	}

	if len(coordinates) < 2 {
		return nil, NewErrorf(true, "Not enough coordinates (%v)", coordinates)
	}

	return &Location{
		UpdatedAt:   t,
		Coordinates: coordinates,
	}, nil
}

func (i *MessageIngester) ApplySchema(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage, schema *MessageSchema, jsonSchema *JsonMessageSchema) (im *IngestedMessage, err error) {
	// This works with MapValues, too as they'll be zero for now.
	if pm.ArrayValues != nil && len(pm.ArrayValues) != len(jsonSchema.Fields) {
		return nil, NewErrorf(true, "%s: fields S=%v != M=%v", pm.SchemaId, len(jsonSchema.Fields), len(pm.ArrayValues))
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
		return nil, NewErrorf(true, "%s: Unable to find stream (%v)", pm.SchemaId, err)
	}

	time, err := determineTime(pm, jsonSchema, mapped)
	if err != nil {
		return nil, NewErrorf(true, "%s: Unable to get time (%v)", pm.SchemaId, err)
	}

	location, err := determineLocation(pm, jsonSchema, mapped, time)
	if err != nil {
		return nil, NewErrorf(true, "%s: Unable to get location (%v)", pm.SchemaId, err)
	}

	haveLocation := location != nil && location.Valid()
	if haveLocation {
		err := i.Streams.UpdateLocation(ctx, pm.SchemaId.Device, stream, location)
		if err != nil {
			return nil, NewErrorf(true, "%s: Unable to update location (%v)", pm.SchemaId, err)
		}
	}
	if stream.Location == nil {
		return nil, NewErrorf(false, "%s: Stream has no location.", pm.SchemaId)
	}

	im = &IngestedMessage{
		Schema:          schema,
		Time:            time,
		Location:        stream.Location,
		LocationUpdated: haveLocation,
		Fixed:           pm.Fixed,
		Fields:          mapped,
	}

	return
}

func (i *MessageIngester) ApplySchemas(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage, schemas []*MessageSchema) (im *IngestedMessage, err error) {
	errors := make([]error, 0)

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
			return nil, NewErrorf(true, "Unexpected MessageSchema type: %v", schema)
		}
	}

	return nil, NewErrorf(true, "%v", makePretty(errors))
}

func (i *MessageIngester) IngestProcessedMessage(ctx context.Context, cache *IngestionCache, pm *ProcessedMessage) (im *IngestedMessage, err error) {
	schemas, err := cache.LookupSchemas(pm.SchemaId, func(id SchemaId) ([]*MessageSchema, error) {
		return i.Schemas.LookupSchemas(ctx, pm.SchemaId)
	})
	if err != nil {
		return nil, NewErrorf(true, "(%s)(%s)[Error] (LookupSchema) %v", pm.MessageId, pm.SchemaId, err)
	}

	if len(schemas) == 0 {
		return nil, NewErrorf(false, "(%s)(%s)[Error] (ApplySchemas) No device or schemas", pm.MessageId, pm.SchemaId)
	}

	im, err = i.ApplySchemas(ctx, cache, pm, schemas)
	if err != nil {
		return nil, NewErrorf(true, "(%s)(%s)[Error] (ApplySchemas) %v", pm.MessageId, pm.SchemaId, err)
	}

	i.Statistics.Successes += 1
	return im, nil
}

func (i *MessageIngester) CreateProcessedMessagFromRawMessage(ctx context.Context, raw *RawMessage) (pm *ProcessedMessage, err error) {
	mp, err := IdentifyMessageProvider(raw)
	if err != nil {
		return nil, err
	}
	if mp == nil {
		return nil, NewErrorf(true, "No message provider: (ContentType: %s)", raw.ContentType)
	}

	pm, err = mp.ProcessMessage(raw)
	if err != nil {
		return nil, NewErrorf(true, "(ProcessMessage) %v", err)
	}

	return pm, nil
}

func makePretty(errors []error) (m string) {
	strs := []string{}
	for _, err := range errors {
		strs = append(strs, fmt.Sprintf("(%v)", err))
	}

	return strings.Join(strs, ", ")
}

type IngestError struct {
	Cause    error
	Critical bool
}

func (e *IngestError) Error() string {
	return e.Cause.Error()
}

func NewErrorf(critical bool, f string, a ...interface{}) *IngestError {
	return &IngestError{
		Cause:    fmt.Errorf(f, a...),
		Critical: critical,
	}
}

func NewError(cause error) *IngestError {
	if err, ok := cause.(*IngestError); ok {
		return &IngestError{
			Cause:    err,
			Critical: err.Critical,
		}
	}
	return &IngestError{
		Cause:    cause,
		Critical: true,
	}

}
