package webhook

import (
	"context"
	"encoding/hex"
	"encoding/json"
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

type ParsedMessage struct {
	original   *WebHookMessage
	deviceID   []byte
	deviceName string
	data       map[string]interface{}
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

	data := make(map[string]interface{})

	for _, module := range schema.Station.Modules {
		for _, sensor := range module.Sensors {
			value, err := m.evaluate(parser, sensor.Expression)
			if err != nil {
				return nil, fmt.Errorf("evaluating sensor expression '%s': %v", sensor.Name, err)
			}

			if sensor.Key != "" {
				data[sensor.Key] = value
			} else {
				data[sensor.Name] = value
			}
		}
	}

	return &ParsedMessage{
		deviceID:   deviceID,
		deviceName: deviceNameString,
		data:       data,
		receivedAt: receivedAt,
		ownerID:    schemaRegistration.OwnerID,
		schemaID:   schemaRegistration.ID,
		schema:     schema,
	}, nil
}

// Messages

type WebHookMessageReceived struct {
	SchemaID  int32 `json:"ttn_schema_id"`
	MessageID int64 `json:"ttn_message_id"`
}

type ProcessSchema struct {
	SchemaID int32 `json:"schema_id"`
}

// Schema

type WebHookSchemaSensor struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

type WebHookSchemaModule struct {
	Key     string                       `json:"key"`
	Name    string                       `json:"name"`
	Sensors []*WebHookSchemaSensor `json:"sensors"`
}

type WebHookSchemaStation struct {
	Key                  string                       `json:"key"`
	Model                string                       `json:"model"`
	IdentifierExpression string                       `json:"identifier"`
	NameExpression       string                       `json:"name"`
	ReceivedExpression   string                       `json:"received"`
	Modules              []*WebHookSchemaModule `json:"modules"`
}

type WebHookSchema struct {
	Station WebHookSchemaStation `json:"station"`
}

type WebHookSchemaRegistration struct {
	ID              int32      `db:"id"`
	OwnerID         int32      `db:"owner_id"`
	Name            string     `db:"name"`
	Token           []byte     `db:"token"`
	Body            []byte     `db:"body"`
	ReceivedAt      *time.Time `db:"received_at"`
	ProcessedAt     *time.Time `db:"processed_at"`
	ProcessInterval *int32     `db:"process_interval"`
}

func (r *WebHookSchemaRegistration) Parse() (*WebHookSchema, error) {
	s := &WebHookSchema{}
	if err := json.Unmarshal(r.Body, s); err != nil {
		return nil, err
	}
	return s, nil
}
