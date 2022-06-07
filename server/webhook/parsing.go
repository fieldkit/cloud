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

type ParsedAttribute struct {
	JSONValue  interface{} `json:"json_value"`
	Location   bool        `json:"location"`
	Associated bool        `json:"associated"`
}

type ParsedMessage struct {
	Original   *WebHookMessage
	Schema     *MessageSchemaStation
	SchemaID   int32
	OwnerID    int32
	ProjectID  *int32
	DeviceID   []byte
	ReceivedAt *time.Time
	Data       []*ParsedReading
	DeviceName *string
	Attributes map[string]*ParsedAttribute
}

func toFloatArray(x interface{}) ([]float64, bool) {
	if arrayValue, ok := x.([]interface{}); ok {
		values := make([]float64, 0)
		for _, opaque := range arrayValue {
			if v, ok := toFloat(opaque); ok {
				values = append(values, v)
			}
		}
		return values, true
	}
	return nil, false
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

func (m *WebHookMessage) tryParse(ctx context.Context, cache *JqCache, schemaRegistration *MessageSchemaRegistration, stationSchema *MessageSchemaStation, source interface{}) (p *ParsedMessage, err error) {
	log := Logger(ctx).Sugar()

	// Check condition expression if one is present. If this returns nothing we
	// skip this message w/o errors.
	if stationSchema.ConditionExpression != "" {
		if _, err := m.evaluate(ctx, cache, source, stationSchema.ConditionExpression); err != nil {
			if _, ok := err.(*EvaluationError); ok {
				// No luck, skipping and maybe another stationSchema will cover this message.
				return nil, nil
			}
			return nil, fmt.Errorf("evaluating condition-expression: %v", err)
		}
	}

	deviceIDRaw, err := m.evaluate(ctx, cache, source, stationSchema.IdentifierExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating identifier-expression: %v", err)
	}

	var deviceName *string

	if stationSchema.NameExpression != "" {
		if deviceNameRaw, err := m.evaluate(ctx, cache, source, stationSchema.NameExpression); err != nil {
			if _, ok := err.(*EvaluationError); !ok {
				return nil, fmt.Errorf("evaluating device-name-expression: %v", err)
			}
		} else {
			if deviceNameString, ok := deviceNameRaw.(string); !ok {
				return nil, fmt.Errorf("unexpected device-name value: %v", deviceNameString)
			} else {
				deviceName = &deviceNameString
			}
		}
	}

	var receivedAt *time.Time

	receivedAtRaw, err := m.evaluate(ctx, cache, source, stationSchema.ReceivedExpression)
	if err != nil {
		if _, ok := err.(*EvaluationError); !ok {
			return nil, fmt.Errorf("evaluating received-at-expression: %v", err)
		}
	} else {
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
		deviceID = []byte(deviceIDString)
	}

	sensors := make([]*ParsedReading, 0)

	for _, module := range stationSchema.Modules {
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
				log.Infow("evaluation-error", "error", err)
				/*
					if _, ok := err.(*EvaluationError); !ok {
						return nil, fmt.Errorf("evaluating sensor expression '%s': %v", sensor.Name, err)
					}
				*/
			} else {
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
	}

	attributes := make(map[string]*ParsedAttribute)

	if stationSchema.Attributes != nil {
		for _, attribute := range stationSchema.Attributes {
			jsonValue, err := m.evaluate(ctx, cache, source, attribute.Expression)
			if err != nil {
				return nil, fmt.Errorf("evaluating attribute expression '%s': %v", attribute.Name, err)
			}

			attributes[attribute.Name] = &ParsedAttribute{
				Location:   attribute.Location,
				Associated: attribute.Associated,
				JSONValue:  jsonValue,
			}
		}
	}

	return &ParsedMessage{
		DeviceID:   deviceID,
		DeviceName: deviceName,
		Data:       sensors,
		Attributes: attributes,
		ReceivedAt: receivedAt,
		OwnerID:    schemaRegistration.OwnerID,
		SchemaID:   schemaRegistration.ID,
		ProjectID:  schemaRegistration.ProjectID,
		Schema:     stationSchema,
	}, nil
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

	for _, stationSchema := range schema.Stations {
		if p, err = m.tryParse(ctx, cache, schemaRegistration, stationSchema, source); err != nil {
			return nil, err
		} else if p != nil {
			break
		}
	}

	return p, nil
}
