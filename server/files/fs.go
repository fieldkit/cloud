package files

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
)

const (
	Path = "./ingestion"
)

type LocalFilesArchive struct {
}

func NewLocalFilesArchive() (a *LocalFilesArchive) {
	return &LocalFilesArchive{}
}

func (a *LocalFilesArchive) Archive(ctx context.Context, meta *FileMeta, reader io.Reader) (*ArchivedFile, error) {
	log := Logger(ctx).Sugar()

	cr := newCountingReader(reader)

	id := uuid.Must(uuid.NewRandom())

	err := os.MkdirAll(Path, 0755)
	if err != nil {
		return nil, err
	}

	fn := makeFileName(id.String())

	log.Infow("archiving", "content_type", meta.ContentType, "file_name", fn)

	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	io.Copy(file, cr)

	ss := &ArchivedFile{
		ID:        id.String(),
		URL:       fn,
		BytesRead: cr.bytesRead,
	}

	return ss, nil

}

func (a *LocalFilesArchive) OpenByKey(ctx context.Context, key string) (io.ReadCloser, error) {
	return a.OpenByURL(ctx, makeFileName(key))
}

func (a *LocalFilesArchive) OpenByURL(ctx context.Context, url string) (io.ReadCloser, error) {
	log := Logger(ctx).Sugar()

	log.Infow("opening", "url", url)

	return os.Open(url)
}

func makeFileName(key string) string {
	return path.Join(Path, fmt.Sprintf("%v.fkpb", key))
}
