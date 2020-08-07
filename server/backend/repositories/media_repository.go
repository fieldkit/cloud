package repositories

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/h2non/filetype"

	"github.com/fieldkit/cloud/server/files"
)

type SavedMedia struct {
	Key      string
	URL      string
	Size     int64
	MimeType string
}

type LoadedMedia struct {
	Key      string
	Size     int64
	MimeType string
	Reader   io.Reader
}

const (
	MaximumRequiredHeaderBytes = 262
)

type MediaRepository struct {
	files files.FileArchive
}

func NewMediaRepository(files files.FileArchive) (r *MediaRepository) {
	return &MediaRepository{
		files: files,
	}
}

func (r *MediaRepository) Save(ctx context.Context, reader io.ReadCloser, contentLength int64, givenContentType string) (sm *SavedMedia, err error) {
	log := Logger(ctx).Sugar()

	if contentLength == 0 {
		return nil, fmt.Errorf("empty")
	}

	headerLength := int64(MaximumRequiredHeaderBytes)
	if contentLength < headerLength {
		headerLength = contentLength
	}

	var buf bytes.Buffer
	headerReader := bufio.NewReader(io.TeeReader(io.LimitReader(reader, headerLength), &buf))

	header := make([]byte, headerLength)
	_, err = io.ReadFull(headerReader, header)
	if err != nil {
		return nil, err
	}

	contentType := givenContentType
	kind, _ := filetype.Match(header)
	if kind == filetype.Unknown {
		log.Warnw("detect file type failed", "given_content_type", givenContentType)
	} else {
		contentType = kind.MIME.Value
	}
	if contentType == "application/octet-stream" {
		return nil, fmt.Errorf("invalid content type")
	}

	objReader := io.MultiReader(bytes.NewReader(buf.Bytes()), reader)
	cr := NewCountingReader(objReader)
	metadata := make(map[string]string)
	af, err := r.files.Archive(ctx, contentType, metadata, cr)
	if err != nil {
		return nil, err
	}

	log.Infow("saved", "content_type", contentType, "id", af.Key, "bytes_read", cr.BytesRead)

	sm = &SavedMedia{
		Key:      af.Key,
		URL:      af.URL,
		Size:     cr.BytesRead,
		MimeType: kind.MIME.Value,
	}

	return
}

func (r *MediaRepository) DeleteByURL(ctx context.Context, url string) (err error) {
	return r.files.DeleteByURL(ctx, url)
}

func (r *MediaRepository) DeleteByKey(ctx context.Context, key string) (err error) {
	return r.files.DeleteByKey(ctx, key)
}

func (r *MediaRepository) LoadByURL(ctx context.Context, url string) (lm *LoadedMedia, err error) {
	opened, err := r.files.OpenByURL(ctx, url)
	if err != nil {
		return nil, err
	}

	if opened == nil {
		return nil, fmt.Errorf("file archive bug, nil opened file: %v", r.files)
	}

	lm = &LoadedMedia{
		Key:      opened.Key,
		Size:     opened.Size,
		MimeType: opened.ContentType,
		Reader:   opened.Body,
	}

	return
}

func (r *MediaRepository) LoadByKey(ctx context.Context, key string) (lm *LoadedMedia, err error) {
	opened, err := r.files.OpenByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if opened == nil {
		return nil, fmt.Errorf("file archive bug, nil opened file: %v", r.files)
	}

	lm = &LoadedMedia{
		Key:      opened.Key,
		Size:     opened.Size,
		MimeType: opened.ContentType,
		Reader:   opened.Body,
	}

	return
}

type CountingReader struct {
	target    io.Reader
	BytesRead int64
}

func NewCountingReader(target io.Reader) *CountingReader {
	return &CountingReader{
		target: target,
	}
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.target.Read(p)
	r.BytesRead += int64(n)
	return n, err
}
