package files

import (
	"context"
	"io"
)

type DevNullFileArchive struct {
}

func (a *DevNullFileArchive) Archive(ctx context.Context, meta *FileMeta, reader io.Reader) (*ArchivedFile, error) {
	Logger(ctx).Sugar().Infof("Fileing %s to /dev/null", meta.ContentType)

	return nil, nil
}

func (a *DevNullFileArchive) OpenByKey(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, nil
}

func (a *DevNullFileArchive) OpenByURL(ctx context.Context, url string) (io.ReadCloser, error) {
	return nil, nil
}
