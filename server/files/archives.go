package files

import (
	"context"
	"io"
)

type ArchivedFile struct {
	Key       string
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
	Key         string
	Size        int64
	ContentType string
	Meta        map[string]string
}

type OpenedFile struct {
	FileInfo
	Body io.ReadCloser
}

type FileArchive interface {
	Archive(ctx context.Context, contentType string, meta map[string]string, read io.Reader) (*ArchivedFile, error)
	OpenByURL(ctx context.Context, url string) (f *OpenedFile, err error)
	DeleteByURL(ctx context.Context, url string) (err error)
	Opened(ctx context.Context, url string, opened *OpenedFile) (reopened *OpenedFile, err error)
	Info(ctx context.Context, key string) (info *FileInfo, err error)
	String() string
}
