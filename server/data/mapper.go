package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/fieldkit/cloud/server/data/jsondocument"
)

var (
	invalidMapperError = errors.New("invalid mapper")
)

type Mapper struct {
	pointers map[string]string
}

func NewMapper(pointers map[string]string) *Mapper {
	m := &Mapper{
		pointers: map[string]string{},
	}

	for dstPointer, srcPointer := range pointers {
		m.pointers[dstPointer] = srcPointer
	}

	return m
}

func (m *Mapper) Map(src *jsondocument.Document) (*jsondocument.Document, error) {
	dst := jsondocument.Object()
	for dstPointer, srcPointer := range m.pointers {
		value, err := src.Document(srcPointer)
		if err != nil {
			return nil, err
		}

		if err := dst.SetDocument(dstPointer, value); err != nil {
			return nil, err
		}
	}

	return dst, nil
}

func (m *Mapper) Pointers() map[string]string {
	pointers := map[string]string{}
	for dstPointer, srcPointer := range m.pointers {
		pointers[dstPointer] = srcPointer
	}

	return pointers
}

func (m *Mapper) Scan(data interface{}) error {
	dataBytes, ok := data.([]byte)
	if !ok {
		return invalidMapperError
	}

	m.pointers = map[string]string{}
	return json.Unmarshal(dataBytes, &m.pointers)
}

func (m *Mapper) Value() (driver.Value, error) {
	return json.Marshal(m.pointers)
}
