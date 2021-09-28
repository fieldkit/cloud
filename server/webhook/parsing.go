package webhook

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/elgs/gojq"
)

type WebHookMessage struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	SchemaID  *int32    `db:"schema_id"`
	Headers   *string   `db:"headers"`
	Body      []byte    `db:"body"`
}

type ParsedReading struct {
	Name     string  `json:"name"`
	Key      string  `json:"key"`
	Value    float64 `json:"value"`
	Battery  bool    `json:"battery"`
	Location bool    `json:"location"`
}

type ParsedMessage struct {
	original   *WebHookMessage
	deviceID   []byte
	deviceName string
	data       []*ParsedReading
	receivedAt time.Time
	schema     *WebHookSchema
	schemaID   int32
	ownerID    int32
}

func (m *WebHookMessage) evaluate(parser *gojq.JQ, query string) (value interface{}, err error) {
	if query == "" {
		return "", fmt.Errorf("empty query")
	}
	raw, err := parser.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying '%s': %v", query, err)
	}

	return raw, nil
}

func (m *WebHookMessage) Parse(ctx context.Context, schemas map[int32]*WebHookSchemaRegistration) (p *ParsedMessage, err error) {
	if m.SchemaID == nil {
		return nil, fmt.Errorf("missing schema")
	}

	schemaRegistration, ok := schemas[*m.SchemaID]
	if !ok {
		return nil, fmt.Errorf("missing schema")
	}

	schema, err := schemaRegistration.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %s", err)
	}

	parser, err := gojq.NewStringQuery(string(m.Body))
	if err != nil {
		return nil, fmt.Errorf("error creating body parser: %v", err)
	}

	deviceIDRaw, err := m.evaluate(parser, schema.Station.IdentifierExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating identifier-expression: %v", err)
	}

	deviceNameRaw, err := m.evaluate(parser, schema.Station.NameExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating device-name-expression: %v", err)
	}

	deviceNameString, ok := deviceNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-name value: %v", deviceNameString)
	}

	receivedAtRaw, err := m.evaluate(parser, schema.Station.ReceivedExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating received-at-expression: %v", err)
	}

	receivedAtString, ok := receivedAtRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected received-at value: %v", receivedAtRaw)
	}

	receivedAt, err := time.Parse("2006-01-02T15:04:05.999999999Z", receivedAtString)
	if err != nil {
		return nil, fmt.Errorf("malformed received at (%s)", receivedAtRaw)
	}

	deviceIDString, ok := deviceIDRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-id value: %v", receivedAtRaw)
	}

	deviceID, err := hex.DecodeString(deviceIDString)
	if err != nil {
		return nil, fmt.Errorf("malformed device eui: %s", deviceIDRaw)
	}

	sensors := make([]*ParsedReading, 0)

	for _, module := range schema.Station.Modules {
		for _, sensor := range module.Sensors {
			maybeValue, err := m.evaluate(parser, sensor.Expression)
			if err != nil {
				return nil, fmt.Errorf("evaluating sensor expression '%s': %v", sensor.Name, err)
			}

			if value, ok := maybeValue.(float64); ok {
				reading := &ParsedReading{
					Name:     sensor.Name,
					Key:      sensor.Key,
					Battery:  sensor.Battery,
					Location: sensor.Location,
					Value:    value,
				}

				sensors = append(sensors, reading)
			} else {
				return nil, fmt.Errorf("non-numeric sensor value '%s'/'%s': %v", sensor.Name, sensor.Expression, maybeValue)
			}
		}
	}

	return &ParsedMessage{
		deviceID:   deviceID,
		deviceName: deviceNameString,
		data:       sensors,
		receivedAt: receivedAt,
		ownerID:    schemaRegistration.OwnerID,
		schemaID:   schemaRegistration.ID,
		schema:     schema,
	}, nil
}
