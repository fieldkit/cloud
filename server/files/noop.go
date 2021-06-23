package files

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/google/uuid"
)

type nopFilesArchive struct {
}

func NewNopFilesArchive() (a *nopFilesArchive) {
	return &nopFilesArchive{}
}

func (a *nopFilesArchive) String() string {
	return "noop"
}

func (a *nopFilesArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*ArchivedFile, error) {
	log := Logger(ctx).Sugar()

	cr := newCountingReader(reader)
	id := uuid.Must(uuid.NewRandom())

	fn := makeFileName(id.String())

	log.Infow("archiving", "content_type", contentType, "file_name", fn)

	io.Copy(ioutil.Discard, cr)

	ss := &ArchivedFile{
		Key:       id.String(),
		URL:       fn,
		BytesRead: cr.bytesRead,
	}

	return ss, nil
}

func (a *nopFilesArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	return nil, fmt.Errorf("unsupported")
}

func (a *nopFilesArchive) Opened(ctx context.Context, url string, opened *OpenedFile) (reopened *OpenedFile, err error) {
	return nil, nil
}

func (a *nopFilesArchive) DeleteByKey(ctx context.Context, key string) (err error) {
	return fmt.Errorf("unsupported")
}

func (a *nopFilesArchive) DeleteByURL(ctx context.Context, url string) (err error) {
	return fmt.Errorf("unsupported")
}

func (a *nopFilesArchive) Info(ctx context.Context, key string) (info *FileInfo, err error) {
	return nil, fmt.Errorf("unsupported")
}
