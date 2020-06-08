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
	ContentType  string
	DeviceID     []byte
	GenerationID []byte
	Blocks       []int64
	Flags        []int64
}

type FileInfo struct {
	Size        int64
	ContentType string
	Meta        map[string]string
}

type FileArchive interface {
	Archive(ctx context.Context, contentType string, meta map[string]string, read io.Reader) (*ArchivedFile, error)
	OpenByKey(ctx context.Context, key string) (io.ReadCloser, error)
	OpenByURL(ctx context.Context, url string) (io.ReadCloser, error)
	Info(ctx context.Context, key string) (info *FileInfo, err error)
	String() string
}
