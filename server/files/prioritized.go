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
	return "prioritized"
}

func (a *prioritizedFilesArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*ArchivedFile, error) {
	var errors *multierror.Error
	for _, a := range a.writing {
		file, err := a.Archive(ctx, contentType, meta, reader)
		if err == nil {
			return file, nil
		}
		errors = multierror.Append(err)
	}
	return nil, errors.ErrorOrNil()
}

func (a *prioritizedFilesArchive) OpenByKey(ctx context.Context, key string) (of *OpenedFile, err error) {
	var errors *multierror.Error
	for _, a := range a.reading {
		reader, err := a.OpenByKey(ctx, key)
		if err == nil {
			return reader, nil
		}
		errors = multierror.Append(err)
	}
	return nil, errors.ErrorOrNil()
}

func (a *prioritizedFilesArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	var errors *multierror.Error
	for _, a := range a.reading {
		reader, err := a.OpenByURL(ctx, url)
		if err == nil {
			return reader, nil
		}
		errors = multierror.Append(err)
	}
	return nil, errors.ErrorOrNil()
}

func (a *prioritizedFilesArchive) DeleteByKey(ctx context.Context, key string) error {
	return nil
}

func (a *prioritizedFilesArchive) DeleteByURL(ctx context.Context, url string) error {
	return nil
}

func (a *prioritizedFilesArchive) Info(ctx context.Context, key string) (info *FileInfo, err error) {
	var errors *multierror.Error
	for _, a := range a.reading {
		info, err := a.Info(ctx, key)
		if err == nil {
			return info, nil
		}

		log := Logger(ctx).Sugar()
		log.Warnw("error", "error", err)

		errors = multierror.Append(err)
	}
	return nil, errors.ErrorOrNil()
}
