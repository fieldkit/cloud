package data

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type FieldkitBinary struct {
	InputID  int32    `db:"input_id"`
	SchemaID int32    `db:"schema_id"`
	ID       uint16   `db:"id"`
	Fields   []string `db:"fields"`
	Mapper   *Mapper  `db:"mapper"`
}

type FieldkitInput struct {
	Input
}

const (
	VarintField  = "varint"
	UvarintField = "uvarint"
	Float32Field = "float32"
	Float64Field = "float64"
)

var (
	fieldkitBinaryUnknownFieldkitBinaryError = errors.New("fieldkit binary: unknown message type")
	fieldkitBinaryKnownFieldTypes            = map[string]bool{
		VarintField:  true,
		UvarintField: true,
		Float32Field: true,
		Float64Field: true,
	}
)

// byteReader implements the io.ByteReader interface. It maintains a single-
// byte buffer in order to avoid allocating a slice on every read.
type byteReader struct {
	reader io.Reader
	buffer [1]byte
}

func newByteReader(r io.Reader) *byteReader {
	return &byteReader{reader: r}
}

// ReadByte reads and returns the next byte from the reader.
func (r *byteReader) ReadByte() (byte, error) {
	_, err := io.ReadFull(r, r.buffer[:])
	if err != nil {
		return 0x0, err
	}

	return r.buffer[0], nil
}

func (r *byteReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

// FieldkitBinaryReader implements the MessageReader interface.
type FieldkitBinaryReader struct {
	reader         *byteReader
	fieldkitBinary map[uint64]*FieldkitBinary
}

func NewFieldkitBinaryReader(r io.Reader) *FieldkitBinaryReader {
	reader := &FieldkitBinaryReader{
		reader:         newByteReader(r),
		fieldkitBinary: make(map[uint64]*FieldkitBinary),
	}

	return reader
}

// AddFieldkitBinary registers a FieldkitBinary
func (m *FieldkitBinaryReader) AddFieldkitBinary(fieldkitBinary *FieldkitBinary) error {
	m.fieldkitBinary[uint64(fieldkitBinary.ID)] = fieldkitBinary
	return nil
}

// ReadMessage reads one message.
func (m *FieldkitBinaryReader) ReadMessage() (interface{}, error) {
	id, err := binary.ReadUvarint(m.reader)
	if err != nil {
		return nil, err
	}

	fieldkitBinary, ok := m.fieldkitBinary[id]
	if !ok {
		return nil, fmt.Errorf("unknown message type, %d", id)
	}

	fieldValues := make([]interface{}, len(fieldkitBinary.Fields))
	for i, field := range fieldkitBinary.Fields {
		var fieldValue interface{}

		switch field {
		case VarintField:
			data, err := binary.ReadVarint(m.reader)
			if err != nil {
				return nil, err
			}

			fieldValue = float64(data)

		case UvarintField:
			data, err := binary.ReadUvarint(m.reader)
			if err != nil {
				return nil, err
			}

			fieldValue = float64(data)

		case Float32Field:
			data := make([]byte, 4)
			if _, err := io.ReadFull(m.reader, data); err != nil {
				return nil, err
			}

			number := float64(math.Float32frombits(binary.LittleEndian.Uint32(data)))
			fieldValue = float64(number)

		case Float64Field:
			data := make([]byte, 8)
			if _, err := io.ReadFull(m.reader, data); err != nil {
				return nil, err
			}

			number := math.Float64frombits(binary.LittleEndian.Uint64(data))
			fieldValue = float64(number)

		default:
			return nil, fmt.Errorf("unknown field type, %s", field)
		}

		fieldValues[i] = fieldValue
	}

	return fieldkitBinary.Mapper.Map(fieldValues)
}
