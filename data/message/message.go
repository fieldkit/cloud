// Package message provides message-parser implementations. Messages are
// parsed (but not mapped) data points.
package message

import (
	"github.com/O-C-R/fieldkit/data/jsondocument"
)

// MessageReader is the interface that wraps the ReadMessage method.
type MessageReader interface {
	ReadMessage() (*jsondocument.Document, error)
}
