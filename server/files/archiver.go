package files

import (
	"context"
	"io"
)

type SavedStream struct {
	ID        string
	URL       string
	BytesRead int
}

type FileMeta struct {
	ContentType string
	DeviceID    []byte
	Generation  []byte
	Blocks      []int64
	Flags       []int64
}

type StreamArchiver interface {
	Archive(ctx context.Context, meta *FileMeta, read io.Reader) (*SavedStream, error)
}
