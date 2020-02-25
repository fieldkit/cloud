package files

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

type FileStreamArchiver struct {
}

func NewFileStreamArchiver() (a *FileStreamArchiver) {
	return &FileStreamArchiver{}
}

func (a *FileStreamArchiver) Archive(ctx context.Context, meta *FileMeta, reader io.Reader) (*SavedStream, error) {
	log := Logger(ctx).Sugar()

	countingReader := newCountingReader(reader)

	id := uuid.Must(uuid.NewRandom())

	path := "./ingestions"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s/%v.fkpb", path, id)

	log.Infow("archiving", "content_type", meta.ContentType, "file_name", fn)

	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	io.Copy(file, countingReader)

	ss := &SavedStream{
		ID:        id.String(),
		URL:       fn,
		BytesRead: countingReader.bytesRead,
	}

	return ss, nil

}
