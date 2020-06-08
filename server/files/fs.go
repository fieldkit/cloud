package files

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"

	"github.com/h2non/filetype"
)

const (
	Path = "./.fs"
)

type localFilesArchive struct {
}

func NewLocalFilesArchive() (a *localFilesArchive) {
	return &localFilesArchive{}
}

func (a *localFilesArchive) String() string {
	return "fs"
}

func (a *localFilesArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*ArchivedFile, error) {
	log := Logger(ctx).Sugar()

	cr := newCountingReader(reader)

	id := uuid.Must(uuid.NewRandom())

	err := os.MkdirAll(Path, 0755)
	if err != nil {
		return nil, err
	}

	fn := makeFileName(id.String())

	log.Infow("archiving", "content_type", contentType, "file_name", fn)

	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	io.Copy(file, cr)

	ss := &ArchivedFile{
		Key:       id.String(),
		URL:       fn,
		BytesRead: cr.bytesRead,
	}

	return ss, nil

}

func (a *localFilesArchive) OpenByKey(ctx context.Context, key string) (of *OpenedFile, err error) {
	return a.OpenByURL(ctx, makeFileName(key))
}

func (a *localFilesArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("opening", "url", url)

	file, err := os.Open(url)
	if err != nil {
		return nil, err
	}

	info, err := a.Info(ctx, url)
	if err != nil {
		return nil, err
	}

	of = &OpenedFile{
		FileInfo: *info,
		Body:     file,
	}

	return
}

func (a *localFilesArchive) DeleteByKey(ctx context.Context, key string) error {
	return fmt.Errorf("unsupported")
}

func (a *localFilesArchive) DeleteByURL(ctx context.Context, url string) error {
	return fmt.Errorf("unsupported")
}

func (a *localFilesArchive) Info(ctx context.Context, keyOrUrl string) (info *FileInfo, err error) {
	names := []string{keyOrUrl, makeFileName(keyOrUrl)}
	for _, fn := range names {
		if si, err := os.Stat(fn); err == nil {
			header, _ := ioutil.ReadFile(fn) // TODO Read less data here.
			kind, _ := filetype.Match(header)
			contentType := kind.MIME.Value

			info = &FileInfo{
				Key:         keyOrUrl,
				Size:        si.Size(),
				Meta:        make(map[string]string),
				ContentType: contentType,
			}

			return info, nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", keyOrUrl)
}

func makeFileName(key string) string {
	return path.Join(Path, fmt.Sprintf("%v.fkpb", key))
}
