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

type SavedStream struct {
	ID  string
	URL string
}

type StreamArchiver interface {
	Archive(ctx context.Context, headers *IncomingHeaders, read io.Reader) (*SavedStream, error)
}

type DevNullStreamArchiver struct {
}

func (a *DevNullStreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (*SavedStream, error) {
	Logger(ctx).Sugar().Infof("Streaming %s to /dev/null", headers.ContentType)

	return nil, nil
}

type FileStreamArchiver struct {
}

func (a *FileStreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (*SavedStream, error) {
	log := Logger(ctx).Sugar()

	if headers.FkUploadName == "" {
		log.Infof("No file name, streaming %s to /dev/null", headers.ContentType)
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		ss := &SavedStream{
			ID:  id.String(),
			URL: "https://foo/" + id.String(),
		}

		return ss, nil
	}

	fn := headers.FkUploadName

	log.Infof("Streaming %s to %s", headers.ContentType, fn)

	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	io.Copy(file, reader)

	ss := &SavedStream{
		ID:  fn,
		URL: fn,
	}

	return ss, nil

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

func (a *S3StreamArchiver) Archive(ctx context.Context, headers *IncomingHeaders, reader io.Reader) (*SavedStream, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(a.session)

	metadata := make(map[string]*string)

	metadata[FkProcessingHeaderName] = aws.String(headers.FkProcessing)
	metadata[FkVersionHeaderName] = aws.String(headers.FkVersion)
	metadata[FkBuildHeaderName] = aws.String(headers.FkBuild)
	metadata[FkDeviceIdHeaderName] = aws.String(headers.FkDeviceId)
	metadata[FkFileIdHeaderName] = aws.String(headers.FkFileId)
	metadata[FkFileNameHeaderName] = aws.String(headers.FkFileName)

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
		return nil, err
	}

	Logger(ctx).Sugar().Infof("Done: %s", r.Location)

	ss := &SavedStream{
		ID:  id.String(),
		URL: r.Location,
	}

	return ss, err
}
