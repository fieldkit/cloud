package data

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

var (
	invalidJSONSchemaError = errors.New("invalid JSON schema")
)

type JSONSchema struct {
	json   interface{}
	schema *gojsonschema.Schema
}

func NewJSONSchema(json string) (js *JSONSchema, err error) {
	js = &JSONSchema{
		json: json,
	}
	return
}

func (s *JSONSchema) MarshalJSON() ([]byte, error) {
	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(s.json); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (s *JSONSchema) UnmarshalJSON(data []byte) error {
	var jsonData interface{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&jsonData); err != nil {
		return err
	}

	return s.UnmarshalGo(jsonData)
}

func (s *JSONSchema) UnmarshalGo(jsonData interface{}) error {
	var err error
	s.schema, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(jsonData))
	if err != nil {
		return err
	}

	s.json = jsonData
	return nil
}

func (s *JSONSchema) Scan(data interface{}) error {
	if data, ok := data.([]byte); ok {
		return s.UnmarshalJSON(data)
	}

	return invalidJSONSchemaError
}

func (s *JSONSchema) Value() (driver.Value, error) {
	return s.MarshalJSON()
}

func (s *JSONSchema) Validate(data interface{}) (bool, map[string]string, error) {
	result, err := s.schema.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return false, nil, err
	}

	if result.Valid() {
		return true, nil, nil
	}

	errors := map[string]string{}
	for _, jsonErr := range result.Errors() {
		errors[jsonErr.Field()] = jsonErr.Description()
	}

	return false, errors, nil
}

type RawSchema struct {
	ID         int32   `db:"id,omitempty"`
	ProjectID  *int32  `db:"project_id"`
	JSONSchema *string `db:"json_schema"`
}

type Schema struct {
	ID         int32       `db:"id,omitempty"`
	ProjectID  *int32      `db:"project_id"`
	JSONSchema *JSONSchema `db:"json_schema"`
}
