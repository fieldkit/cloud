package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fieldkit/cloud/server/files"
)

type memoryFile struct {
	data        []byte
	meta        map[string]string
	contentType string
}

type InMemoryArchive struct {
	files map[string]*memoryFile
}

func NewInMemoryArchive(files map[string][]byte) (a *InMemoryArchive) {
	memoryFiles := make(map[string]*memoryFile)

	for key, data := range files {
		memoryFiles[key] = &memoryFile{
			data:        data,
			meta:        make(map[string]string),
			contentType: "",
		}
	}

	return &InMemoryArchive{
		files: memoryFiles,
	}
}

func (a *InMemoryArchive) String() string {
	return "memory"
}

func (a *InMemoryArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*files.ArchivedFile, error) {
	return nil, fmt.Errorf("unsupported")
}

func (a *InMemoryArchive) OpenByURL(ctx context.Context, url string) (of *files.OpenedFile, err error) {
	if file, ok := a.files[url]; ok {
		info, err := a.Info(ctx, url) // TODO Problem?
		if err != nil {
			return nil, err
		}

		body := ioutil.NopCloser(bytes.NewBuffer(file.data))

		of = &files.OpenedFile{
			FileInfo: *info,
			Body:     body,
		}

		return of, nil
	}
	return nil, fmt.Errorf("no such file: %s", url)
}

func (a *InMemoryArchive) DeleteByURL(ctx context.Context, url string) error {
	if _, ok := a.files[url]; ok {
		a.files[url] = nil
		return nil
	}
	return fmt.Errorf("no such file: %s", url)
}

func (a *InMemoryArchive) Opened(ctx context.Context, url string, opened *files.OpenedFile) (*files.OpenedFile, error) {
	return opened, nil
}

func (a *InMemoryArchive) Info(ctx context.Context, key string) (info *files.FileInfo, err error) {
	if file, ok := a.files[key]; ok {
		info = &files.FileInfo{
			Size:        int64(len(file.data)),
			ContentType: file.contentType,
			Meta:        file.meta,
		}
		return
	}
	return nil, fmt.Errorf("no such file: %s", key)
}
