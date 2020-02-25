package files

import (
	"context"
	"io"
)

type ArchivedFile struct {
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

type FileArchive interface {
	Archive(ctx context.Context, meta *FileMeta, read io.Reader) (*ArchivedFile, error)
	OpenByKey(ctx context.Context, key string) (io.ReadCloser, error)
	OpenByURL(ctx context.Context, url string) (io.ReadCloser, error)
}
