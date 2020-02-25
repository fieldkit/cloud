package files

import (
	"context"
	"io"
)

type DevNullStreamArchiver struct {
}

func (a *DevNullStreamArchiver) Archive(ctx context.Context, meta *FileMeta, reader io.Reader) (*SavedStream, error) {
	Logger(ctx).Sugar().Infof("Streaming %s to /dev/null", meta.ContentType)

	return nil, nil
}
