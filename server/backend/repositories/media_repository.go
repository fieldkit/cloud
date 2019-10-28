package repositories

import (
	"context"
	"fmt"
	"io"

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
	MimeType string
	Size     int
}

type LoadedMedia struct {
	ID     string
	URL    string
	Reader io.Reader
}

const (
	MinimumRequiredBytes = 262
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

	contentType := rd.Header["Content-Type"]
	if len(contentType) != 1 || contentType[0] == "" {
		return nil, fmt.Errorf("invalid content type (empty)")
	}

	cr := NewCountingReader(rd.Body)
	hpr := NewHeaderPreservingReader(cr, 262)

	// NOTE Ideally we'd buffer the first 262 bytes and then classify the
	// image and then decide if we should go to S3 or not.

	id := uuid.Must(uuid.NewRandom())

	log.Infow("saving", "content_type", contentType, "id", id)

	metadata := make(map[string]*string)

	uploader := s3manager.NewUploader(r.session)

	o, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String(contentType[0]),
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(id.String()),
		Body:        hpr,
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return nil, err
	}

	log.Infow("saved", "content_type", contentType, "id", id, "url", o.Location, "bytes_read", cr.BytesRead, "header_size", len(hpr.Header))

	kind, _ := filetype.Match(hpr.Header)
	if kind == filetype.Unknown {
		svc := s3.New(r.session)

		_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(r.bucketName), Key: aws.String(id.String())})
		if err != nil {
			return nil, fmt.Errorf("unable to delete object %q from bucket %q, %v", id, r.bucketName, err)
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(r.bucketName),
			Key:    aws.String(id.String()),
		})
		if err != nil {
			return nil, fmt.Errorf("error deleting invalid file typed upload", err)
		}

		return nil, fmt.Errorf("unknown file type")
	}

	log.Infow("file type", "extension", kind.Extension, "mime_type", kind.MIME.Value)

	sm = &SavedMedia{
		ID:       id.String(),
		URL:      o.Location,
		Size:     cr.BytesRead,
		MimeType: kind.MIME.Value,
	}

	return
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

	lm = &LoadedMedia{
		ID:     id,
		URL:    "",
		Reader: obj.Body,
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

type HeaderPreservingReader struct {
	target io.Reader
	Header []byte
}

func NewHeaderPreservingReader(target io.Reader, sizeToPreserve int) *HeaderPreservingReader {
	return &HeaderPreservingReader{
		target: target,
		Header: make([]byte, 0, sizeToPreserve),
	}
}

func (r *HeaderPreservingReader) Read(p []byte) (n int, err error) {
	n, err = r.target.Read(p)
	avail := cap(r.Header) - len(r.Header)
	if avail > 0 {
		copying := avail
		if copying > n {
			copying = n
		}
		r.Header = append(r.Header, p[:copying]...)
	}
	return n, err
}
