package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/O-C-R/fieldkit/data/jsondocument"
)

const (
	VarintField  = "varint"
	UvarintField = "uvarint"
	Float32Field = "float32"
	Float64Field = "float64"
)

var (
	fieldkitBinaryUnknownMessageTypeError = errors.New("unknown message type")
	knownFieldTypes                       = map[string]bool{
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

type MessageType struct {
	ID     uint64   `json:id`
	Fields []string `json:fields`
}

// FieldkitBinaryReader implements the MessageReader interface.
type FieldkitBinaryReader struct {
	reader       *byteReader
	messageTypes map[uint64][]string
}

func NewFieldkitBinaryReader(r io.Reader) *FieldkitBinaryReader {
	return &FieldkitBinaryReader{
		reader:       newByteReader(r),
		messageTypes: make(map[uint64][]string),
	}
}

// AddMessageType registers a MessageType
func (m *FieldkitBinaryReader) AddMessageType(messageType *MessageType) error {
	fields := make([]string, len(messageType.Fields))
	for i, field := range messageType.Fields {
		if !knownFieldTypes[field] {
			return fmt.Errorf("unknown field type, %s", field)
		}

		fields[i] = field
	}

	m.messageTypes[messageType.ID] = fields
	return nil
}

// ReadMessage reads one message.
func (m *FieldkitBinaryReader) ReadMessage() (*jsondocument.Document, error) {
	id, err := binary.ReadUvarint(m.reader)
	if err != nil {
		return nil, err
	}

	fields, ok := m.messageTypes[id]
	if !ok {
		return nil, fieldkitBinaryUnknownMessageTypeError
	}

	message := jsondocument.Array()
	if err := message.SetDocument("/-", jsondocument.Number(float64(id))); err != nil {
		return nil, err
	}

	for _, field := range fields {
		switch field {
		case VarintField:
			data, err := binary.ReadVarint(m.reader)
			if err != nil {
				return nil, err
			}

			if err := message.SetDocument("/-", jsondocument.Number(float64(data))); err != nil {
				return nil, err
			}

		case UvarintField:
			data, err := binary.ReadUvarint(m.reader)
			if err != nil {
				return nil, err
			}

			if err := message.SetDocument("/-", jsondocument.Number(float64(data))); err != nil {
				return nil, err
			}

		case Float32Field:
			data := make([]byte, 4)
			if _, err := io.ReadFull(m.reader, data); err != nil {
				return nil, err
			}

			number := float64(math.Float32frombits(binary.LittleEndian.Uint32(data)))
			if err := message.SetDocument("/-", jsondocument.Number(number)); err != nil {
				return nil, err
			}

		case Float64Field:
			data := make([]byte, 8)
			if _, err := io.ReadFull(m.reader, data); err != nil {
				return nil, err
			}

			number := math.Float64frombits(binary.LittleEndian.Uint64(data))
			if err := message.SetDocument("/-", jsondocument.Number(number)); err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unknown field type, %s", field)
		}
	}

	return message, nil
}
