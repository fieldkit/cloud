package files

import (
	"context"
	"errors"
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
	var errs *multierror.Error
	for _, a := range a.writing {
		file, err := a.Archive(ctx, contentType, meta, reader)
		if err == nil {
			return file, nil
		}
		errs = multierror.Append(errs, err)
	}

	err := errs.ErrorOrNil()

	if err == nil {
		return nil, errors.New("fatal file archive error")
	}

	return nil, err
}

func (a *prioritizedFilesArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	var errs *multierror.Error
	for _, a := range a.reading {
		reader, err := a.OpenByURL(ctx, url)
		if err == nil {
			return reader, nil
		}
		errs = multierror.Append(errs, err)
	}

	err = errs.ErrorOrNil()

	if err == nil {
		return nil, errors.New("fatal file archive error")
	}

	return nil, err
}

func (a *prioritizedFilesArchive) Opened(ctx context.Context, url string, opened *OpenedFile) error {
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
