package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fieldkit/cloud/server/files"
)

type InMemoryArchive struct {
	files map[string][]byte
}

func NewInMemoryArchive(files map[string][]byte) (a *InMemoryArchive) {
	return &InMemoryArchive{
		files: files,
	}
}

func (a *InMemoryArchive) String() string {
	return "memory"
}

func (a *InMemoryArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*files.ArchivedFile, error) {
	return nil, fmt.Errorf("unsupported")
}

func (a *InMemoryArchive) OpenByKey(ctx context.Context, key string) (io.ReadCloser, error) {
	if d, ok := a.files[key]; ok {
		return ioutil.NopCloser(bytes.NewBuffer(d)), nil
	}
	return nil, fmt.Errorf("no such file: %s", key)
}

func (a *InMemoryArchive) OpenByURL(ctx context.Context, url string) (io.ReadCloser, error) {
	if d, ok := a.files[url]; ok {
		return ioutil.NopCloser(bytes.NewBuffer(d)), nil
	}
	return nil, fmt.Errorf("no such file: %s", url)
}

func (a *InMemoryArchive) Info(ctx context.Context, key string) (meta map[string]string, err error) {
	return nil, fmt.Errorf("unsupported")
}
