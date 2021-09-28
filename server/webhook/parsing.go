package webhook

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/itchyny/gojq"
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

func toFloat(x interface{}) (float64, bool) {
	switch x := x.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case *big.Int:
		f, err := strconv.ParseFloat(x.String(), 64)
		return f, err == nil
	default:
		return 0.0, false
	}
}

func (m *WebHookMessage) evaluate(ctx context.Context, source interface{}, query string) (value interface{}, err error) {
	if query == "" {
		return "", fmt.Errorf("empty query")
	}

	parsed, err := gojq.Parse(query)
	if err != nil {
		return "", fmt.Errorf("error parsing query '%s': %v", query, err)
	}

	compiled, err := gojq.Compile(parsed,
		gojq.WithFunction("clamp", 0, 3, func(x interface{}, xs []interface{}) interface{} {
			if x, ok := toFloat(x); ok {
				if len(xs) == 0 {
					if x < 0 {
						return 0
					}
					if x > 100 {
						return 100
					}
					return x
				}
				if min, ok := toFloat(xs[0]); ok {
					if max, ok := toFloat(xs[1]); ok {
						if x < min {
							return min
						}
						if x > max {
							return max
						}
						return x
					}
				}
			}
			return fmt.Errorf("clamp cannot be applied to: %v, %v", x, xs)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("error compiling query '%s': %v", query, err)
	}

	iter := compiled.Run(source)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return "", fmt.Errorf("query returned error '%s': %v", query, err)
		}

		if v != nil {
			return v, nil
		}
	}

	return "", fmt.Errorf("query returned nothing '%s': %v", query, err)
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
		return nil, fmt.Errorf("error parsing schema: %v", err)
	}

	var source interface{}
	if err := json.Unmarshal(m.Body, &source); err != nil {
		return nil, fmt.Errorf("error parsing message: %v", err)
	}

	deviceIDRaw, err := m.evaluate(ctx, source, schema.Station.IdentifierExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating identifier-expression: %v", err)
	}

	deviceNameRaw, err := m.evaluate(ctx, source, schema.Station.NameExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating device-name-expression: %v", err)
	}

	deviceNameString, ok := deviceNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-name value: %v", deviceNameString)
	}

	receivedAtRaw, err := m.evaluate(ctx, source, schema.Station.ReceivedExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating received-at-expression: %v", err)
	}

	receivedAtString, ok := receivedAtRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected received-at value: %v", receivedAtRaw)
	}

	receivedAt, err := time.Parse("2006-01-02T15:04:05.999999999Z", receivedAtString)
	if err != nil {
		return nil, fmt.Errorf("malformed received-at value: %v", receivedAtRaw)
	}

	deviceIDString, ok := deviceIDRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-id value: %v", receivedAtRaw)
	}

	deviceID, err := hex.DecodeString(deviceIDString)
	if err != nil {
		return nil, fmt.Errorf("malformed device eui: %v", deviceIDRaw)
	}

	sensors := make([]*ParsedReading, 0)

	for _, module := range schema.Station.Modules {
		for _, sensor := range module.Sensors {
			maybeValue, err := m.evaluate(ctx, source, sensor.Expression)
			if err != nil {
				return nil, fmt.Errorf("evaluating sensor expression '%s': %v", sensor.Name, err)
			}

			if value, ok := toFloat(maybeValue); ok {
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
