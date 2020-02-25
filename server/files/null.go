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
