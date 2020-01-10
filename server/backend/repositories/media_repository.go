package repositories

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/goadesign/goa"

	"github.com/h2non/filetype"
)

type SavedMedia struct {
	ID       string
	URL      string
	Size     int
	MimeType string
}

type LoadedMedia struct {
	ID       string
	Size     int
	MimeType string
	Reader   io.Reader
}

const (
	MaximumRequiredHeaderBytes = 262
)

type MediaRepository struct {
	bucketName string
	session    *session.Session
}

func NewMediaRepository(session *session.Session) (r *MediaRepository) {
	return &MediaRepository{
		bucketName: "fk-media",
		session:    session,
	}
}

func (r *MediaRepository) Save(ctx context.Context, rd *goa.RequestData) (sm *SavedMedia, err error) {
	log := Logger(ctx).Sugar()

	contentLength := rd.ContentLength

	headerLength := int64(MaximumRequiredHeaderBytes)
	if contentLength < headerLength {
		headerLength = contentLength
	}

	var buf bytes.Buffer
	headerReader := bufio.NewReader(io.TeeReader(io.LimitReader(rd.Body, headerLength), &buf))

	header := make([]byte, headerLength)
	_, err = io.ReadFull(headerReader, header)
	if err != nil {
		return nil, err
	}

	kind, _ := filetype.Match(header)
	if kind == filetype.Unknown {
		return nil, fmt.Errorf("unknown file type")
	}

	contentType := kind.MIME.Value

	metadata := make(map[string]*string)
	id := uuid.Must(uuid.NewRandom())
	uploader := s3manager.NewUploader(r.session)

	log.Infow("saving", "content_type", contentType, "id", id, "extension", kind.Extension, "mime_type", kind.MIME.Value, "bucket", r.bucketName)

	objReader := io.MultiReader(bytes.NewReader(buf.Bytes()), rd.Body)
	cr := NewCountingReader(objReader)

	o, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String(contentType),
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(id.String()),
		Body:        cr,
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("error uploading: %v", err)
	}

	log.Infow("saved", "content_type", contentType, "id", id, "url", o.Location, "bytes_read", cr.BytesRead, "bucket", r.bucketName)

	sm = &SavedMedia{
		ID:       id.String(),
		URL:      o.Location,
		Size:     cr.BytesRead,
		MimeType: kind.MIME.Value,
	}

	return
}

func (r *MediaRepository) LoadByURL(ctx context.Context, s3url string) (lm *LoadedMedia, err error) {
	u, err := url.Parse(s3url)
	if err != nil {
		return nil, err
	}

	return r.Load(ctx, u.Path[1:])
}

func (r *MediaRepository) Load(ctx context.Context, id string) (lm *LoadedMedia, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("loading", "id", id)

	goi := &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(id),
	}

	svc := s3.New(r.session)

	obj, err := svc.GetObject(goi)
	if err != nil {
		return nil, fmt.Errorf("error reading object %v: %v", id, err)
	}

	contentLength := 0

	if obj.ContentLength != nil {
		contentLength = int(*obj.ContentLength)
	}

	lm = &LoadedMedia{
		ID:       id,
		Size:     contentLength,
		MimeType: "",
		Reader:   obj.Body,
	}

	return
}

type CountingReader struct {
	target    io.Reader
	BytesRead int
}

func NewCountingReader(target io.Reader) *CountingReader {
	return &CountingReader{
		target: target,
	}
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.target.Read(p)
	r.BytesRead += n
	return n, err
}
