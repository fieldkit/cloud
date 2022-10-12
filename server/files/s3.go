package files

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/common/logging"
)

const (
	NoPrefix = ""
)

type S3FileArchive struct {
	session *session.Session
	metrics *logging.Metrics
	bucket  string
	prefix  string
}

func NewS3FileArchive(session *session.Session, metrics *logging.Metrics, bucket string, prefix string) (files *S3FileArchive, err error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	files = &S3FileArchive{
		session: session,
		metrics: metrics,
		bucket:  bucket,
		prefix:  prefix,
	}

	return
}

func (a *S3FileArchive) String() string {
	return "s3"
}

func (a *S3FileArchive) Archive(ctx context.Context, contentType string, meta map[string]string, reader io.Reader) (*ArchivedFile, error) {
	id := uuid.Must(uuid.NewRandom())

	key := id.String()
	if a.prefix != NoPrefix {
		key = a.prefix + id.String()
	}

	log := Logger(ctx).Sugar()

	uploader := s3manager.NewUploader(a.session)

	timer := a.metrics.FileUpload()

	defer timer.Send()

	metadata := make(map[string]*string)
	for key, value := range meta {
		metadata[key] = aws.String(value)
	}

	cr := newCountingReader(reader)

	r, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String(contentType),
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        cr,
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return nil, err
	}

	log.Infow("saved", "url", r.Location, "bytes_read", cr.bytesRead, "key", key)

	ss := &ArchivedFile{
		Key:       key,
		URL:       r.Location,
		BytesRead: cr.bytesRead,
	}

	return ss, err
}

func (a *S3FileArchive) OpenByURL(ctx context.Context, url string) (of *OpenedFile, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("opening", "url", url)

	object, err := common.GetBucketAndKey(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %w", err)
	}

	return a.open(ctx, object.Bucket, object.Key)
}

func (a *S3FileArchive) Opened(ctx context.Context, url string, opened *OpenedFile) (reopened *OpenedFile, err error) {
	return nil, nil
}

func (a *S3FileArchive) DeleteByURL(ctx context.Context, url string) (err error) {
	object, err := common.GetBucketAndKey(url)
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}

	svc := s3.New(a.session)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(object.Bucket), Key: aws.String(object.Key)})
	if err != nil {
		return fmt.Errorf("unable to delete object %q from bucket %q, %w", object.Key, object.Bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})

	return err
}

func (a *S3FileArchive) open(ctx context.Context, bucket, key string) (of *OpenedFile, err error) {
	svc := s3.New(a.session)

	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	log := Logger(ctx).Sugar()

	obj, err := svc.GetObject(goi)
	if err != nil {
		log.Errorw("reading:s3", "bucket", bucket, "key", key)
		return nil, fmt.Errorf("error reading object bucket=%v key=%v: %w", bucket, key, err)
	}

	contentLength := int64(0)
	if obj.ContentLength != nil {
		contentLength = *obj.ContentLength
	}

	contentType := ""
	if obj.ContentType != nil {
		contentType = *obj.ContentType
	}

	log.Infow("opened", "bucket", bucket, "key", key, "content_type", contentType, "content_length", contentLength)

	of = &OpenedFile{
		FileInfo: FileInfo{
			Key:         key,
			Size:        contentLength,
			ContentType: contentType,
			Meta:        make(map[string]string),
		},
		Body: obj.Body,
	}

	return
}

func (a *S3FileArchive) Info(ctx context.Context, key string) (info *FileInfo, err error) {
	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	}

	svc := s3.New(a.session)

	obj, err := svc.HeadObject(hoi)
	if err != nil {
		return nil, fmt.Errorf("error calling HeadObject(%v): %w", key, err)
	}

	maybeMap := common.SanitizeMeta(obj.Metadata)

	meta := make(map[string]string)

	for key, value := range maybeMap {
		meta[key] = *value
	}

	info = &FileInfo{
		Meta:        meta,
		Size:        0,
		ContentType: "",
	}

	return
}
