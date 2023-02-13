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
	Key             string  `json:"key"`
	ModuleKeyPrefix string  `json:"module_key_prefix"`
	FullSensorKey   string  `json:"full_sensor_key"`
	Value           float64 `json:"value"`
	Battery         bool    `json:"battery"`
	Transient       bool    `json:"transient"`
}

type ParsedAttribute struct {
	JSONValue  interface{} `json:"json_value"`
	Location   bool        `json:"location"`
	Associated bool        `json:"associated"`
	Status     bool        `json:"status"`
	Hidden     bool        `json:"hidden"`
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
			return "", fmt.Errorf("error parsing query '%s': %w", query, err)
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
			return "", fmt.Errorf("error compiling query '%s': %w", query, err)
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
			return "", fmt.Errorf("query returned error '%s': %w", query, err)
		}

		if v != nil {
			return v, nil
		}
	}

	return "", &EvaluationError{NoReturn: true, Query: query}
}

func (m *WebHookMessage) evaluateCondition(ctx context.Context, cache *JqCache, expression string, allowEmpty bool, source interface{}) (bool, error) {
	log := Logger(ctx).Sugar()

	if value, err := m.evaluate(ctx, cache, source, expression); err != nil {
		if _, ok := err.(*EvaluationError); ok {
			// No luck, skipping and maybe another stationSchema will cover this message.
			return false, nil
		}
		return false, fmt.Errorf("evaluating condition-expression: %w", err)
	} else if value == nil {
		return false, nil
	} else {
		log.Infow("parsing:condition", "expression", expression, "value", value)

		switch v := value.(type) {
		case string: // TODO Consider a warning to encourage returning bool.
			if !allowEmpty && len(v) == 0 {
				return false, nil
			}
		case bool:
			return v, nil
		}

		return true, nil
	}
}

const (
	ExtractedKey = "extracted"
)

func (m *WebHookMessage) applyExtractors(ctx context.Context, cache *JqCache, stationSchema *MessageSchemaStation, source interface{}) (err error) {
	log := Logger(ctx).Sugar()

	// Because we'll end up inserting the extracted data into the message, we
	// need to be sure that's possible, which it only is on objects.
	sourceObject, ok := source.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected object source")
	}

	// Skip if an extractor has already succeeded for this message.
	if sourceObject[ExtractedKey] != nil {
		log.Infow("extractor:cached")
		return nil
	}

	for _, extractorSpec := range stationSchema.Extractors {
		extractor, err := FindExtractor(ctx, extractorSpec.Type)
		if err != nil {
			return err
		}

		extractingFrom, err := m.evaluate(ctx, cache, source, extractorSpec.Source)
		if err != nil {
			log.Infow("evaluation-error", "error", err)
		}

		extracted, err := extractor.Extract(ctx, source, extractingFrom)
		if err != nil {
			return fmt.Errorf("extraction error: %w", err)
		}

		log.Infow("extractor:done", "extracted", extracted)

		// This is annoying, wish there was a better way to "half serialize" an
		// object, until then we serialize to bytes and then deserialize into
		// opaque/unstructured types.
		serialized, err := json.Marshal(extracted)
		if err != nil {
			return fmt.Errorf("extraction json error: %w", err)
		}

		extractedOpaque := make(map[string]interface{})

		if err := json.Unmarshal(serialized, &extractedOpaque); err != nil {
			return fmt.Errorf("extraction json error: %w", err)
		}

		// Inserts the extracted data/payload into the source object so that the
		// jq expressions can pull information.
		sourceObject[ExtractedKey] = extractedOpaque
	}

	return nil
}

func (m *WebHookMessage) tryParse(ctx context.Context, cache *JqCache, schemaRegistration *MessageSchemaRegistration, stationSchema *MessageSchemaStation, source interface{}) (p *ParsedMessage, err error) {
	log := Logger(ctx).Sugar()

	if err := m.applyExtractors(ctx, cache, stationSchema, source); err != nil {
		return nil, fmt.Errorf("apply-extractors: %w", err)
	}

	// Check condition expression if one is present. If this returns nothing we
	// skip this message w/o errors.
	if stationSchema.ConditionExpression != "" {
		if pass, err := m.evaluateCondition(ctx, cache, stationSchema.ConditionExpression, true, source); err != nil {
			return nil, err
		} else if !pass {
			return nil, nil
		}
	}

	deviceIDRaw, err := m.evaluate(ctx, cache, source, stationSchema.IdentifierExpression)
	if err != nil {
		return nil, fmt.Errorf("evaluating identifier-expression: %w", err)
	}

	var deviceName *string
	if stationSchema.NameExpression != "" {
		if deviceNameRaw, err := m.evaluate(ctx, cache, source, stationSchema.NameExpression); err != nil {
			if _, ok := err.(*EvaluationError); !ok {
				return nil, fmt.Errorf("evaluating device-name-expression: %w", err)
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
	if stationSchema.ReceivedExpression != "" {
		receivedAtRaw, err := m.evaluate(ctx, cache, source, stationSchema.ReceivedExpression)
		if err != nil {
			if _, ok := err.(*EvaluationError); !ok {
				return nil, fmt.Errorf("evaluating received-at-expression: %w", err)
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
							// NOTE: 2022-08-17 16:56:11.835000-04:00
							parsed, err = time.Parse("2006-01-02 15:04:05.000000Z07:00", receivedAtString)
							if err != nil {
								// NOTE: 2022-08-17 16:56:11-04:00
								parsed, err = time.Parse("2006-01-02 15:04:05Z07:00", receivedAtString)
								if err != nil {
									return nil, fmt.Errorf("malformed received-at value: %v", receivedAtRaw)
								}
							}
						}
					}
				}

				receivedAt = &parsed
			} else if receivedAtNumber, ok := receivedAtRaw.(float64); ok {
				parsed := time.Unix(0, int64(receivedAtNumber)*int64(time.Millisecond)).UTC()

				receivedAt = &parsed
			} else {
				return nil, fmt.Errorf("unexpected received-at value: %v", receivedAtRaw)
			}
		}
	}

	deviceIDString, ok := deviceIDRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected device-id value: %v", deviceIDRaw)
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

			log.Infow("parsing:sensor", "sensor_key", sensor.Key)

			expectedKey := strcase.ToLowerCamel(sensor.Key)
			if expectedKey != sensor.Key {
				return nil, fmt.Errorf("unexpected sensor-key formatting '%s' (expected '%s')", sensor.Key, expectedKey)
			}

			if sensor.ConditionExpression != "" {
				if pass, err := m.evaluateCondition(ctx, cache, sensor.ConditionExpression, false, source); err != nil {
					return nil, err
				} else if !pass {
					log.Infow("sensor:skipped:condition")
					continue
				}
			}

			moduleKeyPrefix := module.KeyPrefix()
			fullSensorKey := fmt.Sprintf("%s.%s", moduleKeyPrefix, sensor.Key)

			if sensor.Expression != "" {
				maybeValue, err := m.evaluate(ctx, cache, source, sensor.Expression)
				if err != nil {
					log.Infow("evaluation-error", "error", err)
				} else {
					if value, ok := toFloat(maybeValue); ok {
						filtered := false
						if sensor.Filter != nil && len(*sensor.Filter) == 2 {
							if value < (*sensor.Filter)[0] || value > (*sensor.Filter)[1] {
								filtered = true
							}
						}
						if !filtered {
							reading := &ParsedReading{
								Key:             sensor.Key,
								ModuleKeyPrefix: moduleKeyPrefix,
								FullSensorKey:   fullSensorKey,
								Battery:         sensor.Battery,
								Transient:       sensor.Transient,
								Value:           value,
							}

							sensors = append(sensors, reading)
						} else {
							log.Infow("sensor:skipped:filtered")
						}
					} else {
						return nil, fmt.Errorf("non-numeric sensor value '%s'/'%s': %v", sensor.Name, sensor.Expression, maybeValue)
					}
				}
			} else {
				log.Infow("sensor:skipped:no-expression")
			}
		}
	}

	attributes := make(map[string]*ParsedAttribute)

	if stationSchema.Attributes != nil {
		for _, attribute := range stationSchema.Attributes {
			jsonValue, err := m.evaluate(ctx, cache, source, attribute.Expression)
			if err != nil {
				if _, ok := err.(*EvaluationError); !ok {
					return nil, fmt.Errorf("evaluating attribute expression '%s': %w", attribute.Name, err)
				}
			}

			attributes[attribute.Name] = &ParsedAttribute{
				Location:   attribute.Location,
				Associated: attribute.Associated,
				Hidden:     attribute.Hidden,
				Status:     attribute.Status,
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

func (m *WebHookMessage) unrollArrays(ctx context.Context, source interface{}) ([]interface{}, error) {
	if array, ok := source.([]interface{}); ok {
		return array, nil
	}

	return []interface{}{source}, nil
}

func (m *WebHookMessage) Parse(ctx context.Context, cache *JqCache, schemas map[int32]*MessageSchemaRegistration) ([]*ParsedMessage, error) {
	if m.SchemaID == nil {
		return nil, fmt.Errorf("missing schema id")
	}

	schemaRegistration, ok := schemas[*m.SchemaID]
	if !ok {
		return nil, fmt.Errorf("missing schema")
	}

	schema, err := schemaRegistration.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %w", err)
	}

	var source interface{}
	if err := json.Unmarshal(m.Body, &source); err != nil {
		return nil, fmt.Errorf("error parsing message: %w", err)
	}

	unrolled, err := m.unrollArrays(ctx, source)
	if err != nil {
		return nil, err
	}

	parsed := make([]*ParsedMessage, 0)

	for _, sourceObject := range unrolled {
		for _, stationSchema := range schema.Stations {
			if stationSchema.Flatten {
				sourceObject = flattenObjects(sourceObject)
			}
			if p, err := m.tryParse(ctx, cache, schemaRegistration, stationSchema, sourceObject); err != nil {
				return nil, err
			} else if p != nil {
				parsed = append(parsed, p)
				break
			}
		}
	}

	return parsed, nil
}

func flattenObjects(source interface{}) interface{} {
	if sourceMap, ok := source.(map[string]interface{}); ok {
		flattened := make(map[string]interface{})
		for key, value := range sourceMap {
			if childMap, ok := value.(map[string]interface{}); ok {
				for childKey, childValue := range childMap {
					flattened[childKey] = childValue
				}
			} else {
				flattened[key] = value
			}
		}
		return flattened
	}
	return source
}
