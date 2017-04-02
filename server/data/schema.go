package data

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"
)

type JSONSchema struct {
	json   interface{}
	schema *gojsonschema.Schema
}

func (s *JSONSchema) MarshalJSON() ([]byte, error) {
	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(s.json); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (s *JSONSchema) UnmarshalJSON(data []byte) error {
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&s.json); err != nil {
		return err
	}

	var err error
	s.schema, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(s.json))
	if err != nil {
		return err
	}
	return nil
}

func (s *JSONSchema) Scan(data []byte) error {
	return s.UnmarshalJSON(data)
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

type Schema struct {
	ID         int32       `db:"id,omitempty"`
	ProjectID  *int32      `db:"project_id"`
	JSONSchema *JSONSchema `db:"json_schema"`
}
