package webhook

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/gojq"

	"github.com/iancoleman/strcase"
)

type ParsedReading struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Battery   bool    `json:"battery"`
	Transient bool    `json:"transient"`
}

type ParsedMessage struct {
	Original   *WebHookMessage
	DeviceID   []byte
	DeviceName string
	Data       []*ParsedReading
	ReceivedAt time.Time
	Schema     *MessageSchema
	SchemaID   int32
	OwnerID    int32
	Location   []float64
}

func toFloat(x interface{}) (float64, bool) {
	switch x := x.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	case *big.Int:
		f, err := strconv.ParseFloat(x.String(), 64)
		return f, err == nil
	default:
		return 0.0, false
	}
}

type JqCache struct {
	compiled map[string]*gojq.Code
}

type EvaluationError struct {
	Query    string
	NoReturn bool
}

func (m *EvaluationError) Error() string {
	return "EvaluationError"
}

func (m *WebHookMessage) evaluate(ctx context.Context, cache *JqCache, source interface{}, query string) (value interface{}, err error) {
	if query == "" {
		return "", fmt.Errorf("empty query")
	}

	compiled, ok := cache.compiled[query]
	if !ok {
		parsed, err := gojq.Parse(query)
		if err != nil {
			return "", fmt.Errorf("error parsing query '%s': %v", query, err)
		}

		compiled, err = gojq.Compile(parsed,
			gojq.WithFunction("coerce_to_number", 0, 0, func(x interface{}, xs []interface{}) interface{} {
				if x, ok := x.(string); ok {
					f, err := strconv.ParseFloat(strings.Replace(x, "i", "", 1), 32)
					if err != nil {
						return err
					}
					return f
				}
				return x
			}),
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

		if cache.compiled == nil {
			cache.compiled = make(map[string]*gojq.Code)
		}

		cache.compiled[query] = compiled
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

	return "", &EvaluationError{NoReturn: true, Query: query}
}

func (m *WebHookMessage) Parse(ctx context.Context, cache *JqCache, schemas map[int32]*MessageSchemaRegistration) (p *ParsedMessage, err error) {
	if m.SchemaID == nil {
		return nil, fmt.Errorf("missing schema id")
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

	// Check condition expression if one is present. If this returns nothing we
	// skip this message w/o errors.
	if schema.Station.ConditionExpression != "" {
		if _, err := m.evaluate(ctx, cache, source, schema.Station.ConditionExpression); err != nil {
			if _, ok := err.(*EvaluationError); ok {
				return nil, nil
			}
			return nil, fmt.Errorf("evaluating condition-expression: %v", err)
		}
	}

	deviceIDRaw, err := m.evaluate(ctx, cache, source, schema.Station.IdentifierExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating identifier-expression: %v", err)
	}

	deviceNameString := ""

	if schema.Station.NameExpression != "" {
		deviceNameRaw, err := m.evaluate(ctx, cache, source, schema.Station.NameExpression)
		if err != nil {
			return nil, fmt.Errorf("evaluating device-name-expression: %v", err)
		}

		deviceNameString, ok := deviceNameRaw.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected device-name value: %v", deviceNameString)
		}
	}

	receivedAtRaw, err := m.evaluate(ctx, cache, source, schema.Station.ReceivedExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating received-at-expression: %v", err)
	}

	var receivedAt *time.Time
	if receivedAtString, ok := receivedAtRaw.(string); ok {
		parsed, err := time.Parse("2006-01-02T15:04:05.999999999Z", receivedAtString)
		if err != nil {
			parsed, err = time.Parse("2006-01-02 15:04:05.999999999+00:00", receivedAtString)
			if err != nil {
				// NOTE: NOAA Tidal data was missing seconds.
				parsed, err = time.Parse("2006-01-02 15:04+00:00", receivedAtString)
				if err != nil {
					return nil, fmt.Errorf("malformed received-at value: %v", receivedAtRaw)
				}
			}
		}

		receivedAt = &parsed
	} else if receivedAtNumber, ok := receivedAtRaw.(float64); ok {
		parsed := time.Unix(0, int64(receivedAtNumber)*int64(time.Millisecond))

		receivedAt = &parsed
	} else {
		return nil, fmt.Errorf("unexpected received-at value: %v", receivedAtRaw)
	}

	if receivedAt == nil {
		return nil, fmt.Errorf("missing received-at")
	}

	deviceIDString, ok := deviceIDRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-id value: %v", receivedAtRaw)
	}

	if len(deviceIDString) == 0 {
		return nil, fmt.Errorf("malformed device eui: %v", deviceIDRaw)
	}

	deviceID, err := hex.DecodeString(deviceIDString)
	if err != nil {
		if false {
			return nil, fmt.Errorf("malformed device eui: %v", deviceIDRaw)
		}
		deviceID = []byte(deviceIDString)
	}

	sensors := make([]*ParsedReading, 0)

	for _, module := range schema.Station.Modules {
		for _, sensor := range module.Sensors {
			if sensor.Key == "" {
				return nil, fmt.Errorf("empty sensor-key")
			}

			expectedKey := strcase.ToLowerCamel(sensor.Key)
			if expectedKey != sensor.Key {
				return nil, fmt.Errorf("unexpected sensor-key formatting '%s' (expected '%s')", sensor.Key, expectedKey)
			}

			maybeValue, err := m.evaluate(ctx, cache, source, sensor.Expression)
			if err != nil {
				return nil, fmt.Errorf("evaluating sensor expression '%s': %v", sensor.Name, err)
			}

			if value, ok := toFloat(maybeValue); ok {
				reading := &ParsedReading{
					Key:       sensor.Key,
					Battery:   sensor.Battery,
					Transient: sensor.Transient,
					Value:     value,
				}

				sensors = append(sensors, reading)
			} else {
				return nil, fmt.Errorf("non-numeric sensor value '%s'/'%s': %v", sensor.Name, sensor.Expression, maybeValue)
			}
		}
	}

	var location []float64

	if schema.Station.LongitudeExpression != "" && schema.Station.LatitudeExpression != "" {
		maybeLongitude, err := m.evaluate(ctx, cache, source, schema.Station.LongitudeExpression)
		if err != nil {
			return nil, fmt.Errorf("evaluating longitude expression: %v", err)
		}

		maybeLatitude, err := m.evaluate(ctx, cache, source, schema.Station.LatitudeExpression)
		if err != nil {
			return nil, fmt.Errorf("evaluating latitude expression: %v", err)
		}

		if maybeLongitude != nil && maybeLatitude != nil {
			if longitude, ok := maybeLongitude.(float64); ok {
				if latitude, ok := maybeLatitude.(float64); ok {
					location = []float64{longitude, latitude}
				}
			}
		}
	}

	return &ParsedMessage{
		DeviceID:   deviceID,
		DeviceName: deviceNameString,
		Data:       sensors,
		ReceivedAt: *receivedAt,
		OwnerID:    schemaRegistration.OwnerID,
		SchemaID:   schemaRegistration.ID,
		Schema:     schema,
		Location:   location,
	}, nil
}
