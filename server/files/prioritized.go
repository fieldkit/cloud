package files

import (
	"context"
	"io"

	"github.com/hashicorp/go-multierror"
)

type prioritizedFilesArchive struct {
	reading []FileArchive
	writing []FileArchive
}

func NewPrioritizedFilesArchive(reading []FileArchive, writing []FileArchive) (a FileArchive) {
	return &prioritizedFilesArchive{
		reading: reading,
		writing: writing,
	}
}

func (a *prioritizedFilesArchive) String() string {
	return "priotiized"
}

func (a *prioritizedFilesArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*ArchivedFile, error) {
	var errors error
	for _, a := range a.writing {
		file, err := a.Archive(ctx, contentType, meta, reader)
		if err == nil {
			return file, nil
		}
		errors = multierror.Append(errors)
	}
	return nil, errors
}

func (a *prioritizedFilesArchive) OpenByKey(ctx context.Context, key string) (io.ReadCloser, error) {
	var errors error
	for _, a := range a.reading {
		reader, err := a.OpenByKey(ctx, key)
		if err == nil {
			return reader, nil
		}
		errors = multierror.Append(errors)
	}
	return nil, errors
}

func (a *prioritizedFilesArchive) OpenByURL(ctx context.Context, url string) (io.ReadCloser, error) {
	var errors error
	for _, a := range a.reading {
		reader, err := a.OpenByURL(ctx, url)
		if err == nil {
			return reader, nil
		}
		errors = multierror.Append(errors)
	}
	return nil, errors
}

func (a *prioritizedFilesArchive) Info(ctx context.Context, key string) (info *FileInfo, err error) {
	var errors error
	for _, a := range a.reading {
		info, err := a.Info(ctx, key)
		if err == nil {
			return info, nil
		}
		errors = multierror.Append(errors)
	}
	return nil, errors
}
