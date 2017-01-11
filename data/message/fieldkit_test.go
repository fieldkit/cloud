package message

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"testing"

	"github.com/O-C-R/fieldkit/data/jsondocument"
)

const (
	testFieldkitBinaryMessageTypeJSON = `
{
	"id": 1,
	"fields": [
		"varint",
		"uvarint",
		"float32",
		"float64"
	]
}
`
)

func writeMessage(w io.Writer, id uint64, fields ...interface{}) error {
	value := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(value, id)
	value = value[:n]
	if _, err := w.Write(value); err != nil {
		return err
	}

	for _, field := range fields {
		var data []byte
		switch field := field.(type) {
		case int64:
			data = make([]byte, binary.MaxVarintLen64)
			n := binary.PutVarint(data, field)
			data = data[:n]
		case uint64:
			data = make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(data, field)
			data = data[:n]
		case float32:
			data = make([]byte, 4)
			binary.LittleEndian.PutUint32(data, math.Float32bits(field))
		case float64:
			data = make([]byte, 8)
			binary.LittleEndian.PutUint64(data, math.Float64bits(field))
		}

		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}

func arrayIndexNumber(t *testing.T, d *jsondocument.Document, index int) float64 {
	data, err := d.ArrayIndex(index)
	if err != nil {
		t.Fatal(err)
	}

	number, err := data.Number()
	if err != nil {
		t.Fatal(err)
	}

	return number
}

func TestFieldkitBinaryReader(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	if err := writeMessage(buffer, uint64(1), int64(2), uint64(3), float32(math.E), float64(math.Pi)); err != nil {
		t.Fatal(err)
	}

	messageType := &MessageType{}
	if err := json.Unmarshal([]byte(testFieldkitBinaryMessageTypeJSON), messageType); err != nil {
		t.Fatal(err)
	}

	m := NewFieldkitBinaryReader(buffer)
	if err := m.AddMessageType(messageType); err != nil {
		t.Fatal(err)
	}

	t.Logf("%08b", buffer.Bytes())

	message, err := m.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(messageJSON))

	if uint64(arrayIndexNumber(t, message, 0)) != uint64(1) {
		t.Error("unexpected value at index 0")
	}

	if int64(arrayIndexNumber(t, message, 1)) != int64(2) {
		t.Error("unexpected value at index 1")
	}

	if uint64(arrayIndexNumber(t, message, 2)) != uint64(3) {
		t.Error("unexpected value at index 2")
	}

	if float32(arrayIndexNumber(t, message, 3)) != float32(math.E) {
		t.Error("unexpected value at index 3")
	}

	if float64(arrayIndexNumber(t, message, 4)) != float64(math.Pi) {
		t.Error("unexpected value at index 4")
	}
}
