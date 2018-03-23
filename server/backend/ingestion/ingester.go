package ingestion

import (
	"context"
	"fmt"
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

type MessageIngester struct {
	Repository *Repository
}

const (
	FieldNameLatitude  = "latitude"
	FieldNameLongitude = "longitude"
	FieldNameAltitude  = "altitude"
	FieldNameTime      = "time"
)

func NewMessageIngester(r *Repository) *MessageIngester {
	return &MessageIngester{
		Repository: r,
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

func (i *MessageIngester) ApplySchema(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage, schema *MessageSchema, jsonSchema *JsonMessageSchema) (im *IngestedMessage, err error) {
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
		err := i.Repository.UpdateLocation(ctx, pm.SchemaId.Device, ds.Stream, location)
		if err != nil {
			return nil, NewErrorf(true, "%s: Unable to update location (%v)", pm.SchemaId, err)
		}
	}
	if ds.Stream.Location == nil {
		return nil, NewErrorf(false, "%s: Stream has no location.", pm.SchemaId)
	}

	im = &IngestedMessage{
		Schema:          schema,
		Time:            time,
		Location:        ds.Stream.Location,
		LocationUpdated: haveLocation,
		Fixed:           pm.Fixed,
		Fields:          mapped,
	}

	return
}

func (i *MessageIngester) ApplySchemas(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage) (im *IngestedMessage, err error) {
	errors := make([]error, 0)

	for _, schema := range ds.Schemas {
		jsonSchema, ok := schema.Schema.(*JsonMessageSchema)
		if ok {
			im, applyErr := i.ApplySchema(ctx, ds, pm, schema, jsonSchema)
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

func (i *MessageIngester) IngestProcessedMessage(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage) (im *IngestedMessage, err error) {
	if len(ds.Schemas) == 0 {
		return nil, NewErrorf(false, "(%s)(%s)[Error] (ApplySchemas) No device or schemas", pm.MessageId, pm.SchemaId)
	}

	im, err = i.ApplySchemas(ctx, ds, pm)
	if err != nil {
		return nil, NewErrorf(true, "(%s)(%s)[Error] (ApplySchemas) %v", pm.MessageId, pm.SchemaId, err)
	}

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
