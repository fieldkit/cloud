package data

import (
	"errors"
	"io"

	"github.com/lib/pq"
)

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

type FieldkitInput struct {
	Input
}

type StringSlice []string

func (s StringSlice) Scan(data interface{}) error {
	return pq.Array(&s).Scan(data)
}

type FieldkitBinary struct {
	InputID   int32    `db:"input_id"`
	SchemaID  int32    `db:"schema_id"`
	ID        uint16   `db:"id"`
	Fields    []string `db:"fields"`
	Mapper    *Mapper  `db:"mapper"`
	Longitude *string  `db:"longitude"`
	Latitude  *string  `db:"longitude"`
}

// ByteReader implements the io.ByteReader interface. It maintains a single-
// byte buffer in order to avoid allocating a slice on every read.
type ByteReader struct {
	reader io.Reader
	buffer [1]byte
}

func NewByteReader(r io.Reader) *ByteReader {
	return &ByteReader{reader: r}
}

// ReadByte reads and returns the next byte from the reader.
func (r *ByteReader) ReadByte() (byte, error) {
	_, err := io.ReadFull(r, r.buffer[:])
	if err != nil {
		return 0x0, err
	}

	return r.buffer[0], nil
}

func (r *ByteReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

// // ReadMessage reads one message.
// func (m *FieldkitBinary) ReadMessage(r *ByteReader) (interface{}, error) {
// 	fieldValues := make([]interface{}, len(m.Fields))
// 	for i, field := range m.Fields {
// 		var fieldValue interface{}

// 		switch field {
// 		case VarintField:
// 			data, err := binary.ReadVarint(r)
// 			if err != nil {
// 				return nil, err
// 			}

// 			fieldValue = float64(data)

// 		case UvarintField:
// 			data, err := binary.ReadUvarint(r)
// 			if err != nil {
// 				return nil, err
// 			}

// 			fieldValue = float64(data)

// 		case Float32Field:
// 			data := make([]byte, 4)
// 			if _, err := io.ReadFull(r, data); err != nil {
// 				return nil, err
// 			}

// 			number := float64(math.Float32frombits(binary.LittleEndian.Uint32(data)))
// 			fieldValue = float64(number)

// 		case Float64Field:
// 			data := make([]byte, 8)
// 			if _, err := io.ReadFull(r, data); err != nil {
// 				return nil, err
// 			}

// 			number := math.Float64frombits(binary.LittleEndian.Uint64(data))
// 			fieldValue = float64(number)

// 		default:
// 			return nil, fmt.Errorf("unknown field type, %s", field)
// 		}

// 		fieldValues[i] = fieldValue
// 	}

// 	return m.Mapper.Map(fieldValues)
// }
