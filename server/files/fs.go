package files

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"

	"github.com/h2non/filetype"
)

const (
	Path = ".fs"
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

	copied, err := io.Copy(file, cr)
	if err != nil {
		return nil, err
	}

	log.Infow("saved", "bytes_read", copied, "key", id.String())

	ss := &ArchivedFile{
		Key:       id.String(),
		URL:       fn,
		BytesRead: int(copied),
	}

	return ss, nil

}

func (a *localFilesArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	log := Logger(ctx).Sugar()

	fn := makeFileName(url)

	log.Infow("opening", "url", url, "filename", fn)

	file, err := os.Open(fn)
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

func (a *localFilesArchive) Opened(ctx context.Context, url string, opened *OpenedFile) (reopened *OpenedFile, err error) {
	log := Logger(ctx).Sugar()

	fn := makeFileName(url)

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		log.Infow("opened, resaving", "url", url, "fn", fn)

		file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		_, err = io.Copy(file, opened.Body)
		if err != nil {
			return nil, err
		}

		rof, err := a.OpenByURL(ctx, url)
		if err != nil {
			return nil, err
		}
		if rof == nil {
			panic("wtf")
		}
		return rof, nil
	} else {
		log.Infow("opened, skipped", "url", url, "fn", fn)
	}

	return nil, nil
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

func makeFileName(keyOrUrl string) string {
	u, err := url.Parse(keyOrUrl)
	if err != nil {
		return path.Join(Path, fmt.Sprintf("%v", keyOrUrl))
	}
	if strings.HasPrefix(u.Path, Path) {
		return u.Path
	}
	return path.Join(Path, fmt.Sprintf("%v", u.Path))
}
