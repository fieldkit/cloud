package ingestion

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	FkBinaryMessagesSpace = uuid.Must(uuid.Parse("0b8a5016-7410-4a1a-a2ed-2c48fec6903d"))
)

type DatabaseIds struct {
	SchemaID int64
	DeviceID int64
}

type MessageSchema struct {
	Ids    DatabaseIds
	Schema interface{}
}

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	Ids                 interface{}
	UseProviderTime     bool
	UseProviderLocation bool
	HasTime             bool
	HasLocation         bool
	Fields              []JsonSchemaField
}

type MessageId string

type DeviceId string

func NewDeviceId(id string) DeviceId {
	return DeviceId(id)
}

func NewProviderDeviceId(provider, id string) DeviceId {
	return DeviceId(fmt.Sprintf("%s-%s", provider, id))
}

func (s DeviceId) String() string {
	return fmt.Sprintf("%s", string(s))
}

type SchemaId struct {
	Device DeviceId
	Stream string
}

func (s SchemaId) String() string {
	if s.Stream == "" {
		return s.Device.String()
	}
	return fmt.Sprintf("%s-%s", s.Device.String(), s.Stream)
}

func NewSchemaId(deviceId DeviceId, stream string) SchemaId {
	return SchemaId{
		Device: deviceId,
		Stream: stream,
	}
}

type FormattedMessage struct {
	MessageId   MessageId
	SchemaId    SchemaId
	Time        *time.Time
	Location    []float64
	Fixed       bool
	ArrayValues []string
	MapValues   map[string]interface{}
	Modules     []string
}

type ProcessedMessage struct {
	Schema          *MessageSchema
	Time            *time.Time
	Location        *Location
	LocationUpdated bool
	Fixed           bool
	Fields          map[string]interface{}
}

type SchemaApplier struct {
}

const (
	FieldNameLatitude  = "latitude"
	FieldNameLongitude = "longitude"
	FieldNameAltitude  = "altitude"
	FieldNameTime      = "time"
)

func NewSchemaApplier() *SchemaApplier {
	return &SchemaApplier{}
}

func (i *SchemaApplier) ApplySchema(ds *DeviceStream, fm *FormattedMessage, schema *MessageSchema, jsonSchema *JsonMessageSchema) (im *ProcessedMessage, err error) {
	// This works with MapValues, too as they'll be zero for now.
	if fm.ArrayValues != nil && len(fm.ArrayValues) != len(jsonSchema.Fields) {
		return nil, fmt.Errorf("%s: fields S=%v != M=%v", fm.SchemaId, len(jsonSchema.Fields), len(fm.ArrayValues))
	}

	mapped := make(map[string]interface{})

	for i, field := range jsonSchema.Fields {
		mapped[ToSnake(field.Name)] = fm.ArrayValues[i]
	}

	if fm.MapValues != nil {
		for k, v := range fm.MapValues {
			mapped[k] = v
		}
	}

	time, err := getTime(fm, jsonSchema, mapped)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get time (%v)", fm.SchemaId, err)
	}

	location, err := getLocation(fm, jsonSchema, mapped, time)
	if err != nil {
		return nil, fmt.Errorf("%s: Unable to get location (%v)", fm.SchemaId, err)
	}

	haveLocation := location != nil && location.Valid()
	if haveLocation {
		ds.Stream.Location = location
	}
	if ds.Stream.Location == nil {
		return nil, fmt.Errorf("%s: Stream has no location.", fm.SchemaId)
	}

	im = &ProcessedMessage{
		Schema:          schema,
		Time:            time,
		Location:        ds.Stream.Location,
		LocationUpdated: haveLocation,
		Fixed:           fm.Fixed,
		Fields:          mapped,
	}

	return
}

func (i *SchemaApplier) ApplySchemas(ds *DeviceStream, fm *FormattedMessage) (im *ProcessedMessage, err error) {
	if len(ds.Schemas) == 0 {
		return nil, NewErrorf(false, "(%s)(%s)[Error] (ApplySchemas) No device or schemas", fm.MessageId, fm.SchemaId)
	}

	errors := make([]error, 0)

	for _, schema := range ds.Schemas {
		jsonSchema, ok := schema.Schema.(*JsonMessageSchema)
		if ok {
			fm, applyErr := i.ApplySchema(ds, fm, schema, jsonSchema)
			if applyErr != nil {
				errors = append(errors, applyErr)
			} else {
				return fm, nil
			}
		} else {
			return nil, fmt.Errorf("Unexpected MessageSchema type: %v", schema)
		}
	}

	return nil, fmt.Errorf("%v", joinErrors(errors))
}

func getTime(fm *FormattedMessage, ms *JsonMessageSchema, m map[string]interface{}) (t *time.Time, err error) {
	if ms.UseProviderTime {
		t = fm.Time
	} else {
		if !ms.HasTime {
			return nil, fmt.Errorf("No time information")
		}

		raw := m[FieldNameTime].(string)
		unix, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse time '%s' (%v)", raw, err)
		}

		parsed := time.Unix(unix, 0)
		t = &parsed
	}

	return
}

func getLocation(fm *FormattedMessage, ms *JsonMessageSchema, m map[string]interface{}, t *time.Time) (l *Location, err error) {
	if ms.UseProviderLocation {
		return &Location{UpdatedAt: t, Coordinates: fm.Location}, nil
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

func joinErrors(errors []error) (m string) {
	strs := []string{}
	for _, err := range errors {
		strs = append(strs, fmt.Sprintf("(%v)", err))
	}

	return strings.Join(strs, ", ")
}
