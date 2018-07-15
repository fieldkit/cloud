package backend

import (
	"context"
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StreamArchiver interface {
	Archive(ctx context.Context, headers *IncomingHeaders, read io.Reader) (string, error)
}

type DevNullStreamArchiver struct {
}

func (a *DevNullStreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (string, error) {
	Logger(ctx).Sugar().Infof("Streaming %s to /dev/null", headers.ContentType)

	return "", nil
}

type FileStreamArchiver struct {
}

func (a *FileStreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (string, error) {
	log := Logger(ctx).Sugar()

	if headers.FkUploadName != "" {
		fn := headers.FkUploadName

		log.Infof("Streaming %s to %s", headers.ContentType, fn)

		file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return "", err
		}

		defer file.Close()

		io.Copy(file, reader)
	} else {
		log.Infof("Streaming %s to /dev/null (No UploadName)", headers.ContentType)

	}

	return "", nil
}

type S3StreamArchiver struct {
	session    *session.Session
	bucketName string
}

func NewS3StreamArchiver(session *session.Session, bucketName string) *S3StreamArchiver {
	return &S3StreamArchiver{
		session:    session,
		bucketName: bucketName,
	}
}

func (a *S3StreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(a.session)

	metadata := make(map[string]*string)

	metadata[FkProcessingHeaderName] = aws.String(headers.FkProcessing)
	metadata[FkVersionHeaderName] = aws.String(headers.FkVersion)
	metadata[FkBuildHeaderName] = aws.String(headers.FkBuild)
	metadata[FkDeviceIdHeaderName] = aws.String(headers.FkDeviceId)
	metadata[FkFileIdHeaderName] = aws.String(headers.FkFileId)

	r, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String(headers.ContentType),
		Bucket:      aws.String(a.bucketName),
		Key:         aws.String(id.String()),
		Body:        reader,
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return "", err
	}

	Logger(ctx).Sugar().Infof("Done: %s", r.Location)

	return id.String(), err
}
