package backend

import (
	"bytes"
	"context"
	"fmt"
	_ "io"
	"mime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/fieldkit/cloud/server/errors"
)

type StreamProcessor interface {
	Process(ctx context.Context, id string) error
}

type S3StreamProcessor struct {
	session    *session.Session
	bucketName string
	ingester   *StreamIngester
}

func NewS3StreamProcessor(session *session.Session, ingester *StreamIngester, bucketName string) *S3StreamProcessor {
	return &S3StreamProcessor{
		session:    session,
		bucketName: bucketName,
		ingester:   ingester,
	}
}

func (a *S3StreamProcessor) Process(ctx context.Context, id string) error {
	log := Logger(ctx).Sugar()

	svc := s3.New(a.session)

	hoo, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return errors.Structured(fmt.Errorf("S3 Error: %v", err), "stream_id", id)
	}

	if hoo.ContentLength == nil || hoo.ContentType == nil {
		return errors.Structured(fmt.Errorf("Missing ContentLength/ContentType"), "stream_id", id)
	}

	url := fmt.Sprintf("https://%s.amazonaws.com/%s/%s", *a.session.Config.Region, a.bucketName, id)

	metadata := make(map[string]string)
	for key, value := range hoo.Metadata {
		metadata[key] = *value
	}

	mediaType, _, err := mime.ParseMediaType(*hoo.ContentType)
	if err != nil {
		return err
	}

	headers := &IncomingHeaders{
		ContentType:   *hoo.ContentType,
		ContentLength: int32(*hoo.ContentLength),
		MediaType:     mediaType,
		FkProcessing:  metadata[FkProcessingHeaderName],
		FkVersion:     metadata[FkVersionHeaderName],
		FkBuild:       metadata[FkBuildHeaderName],
		FkDeviceId:    metadata[FkDeviceIdHeaderName],
		FkFileId:      metadata[FkFileIdHeaderName],
	}

	log.Infow("downloading", "size", hoo.ContentLength, "stream_id", id, "url", url, "modified", hoo.LastModified, "metadata", metadata)

	downloader := s3manager.NewDownloader(a.session)

	buff := &aws.WriteAtBuffer{}
	_, err = downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(a.bucketName),
			Key:    aws.String(id),
		})
	if err != nil {
		return errors.Structured(fmt.Errorf("S3 Error: %v", err), "stream_id", id)
	}

	reader := bytes.NewReader(buff.Bytes())

	return a.ingester.download(ctx, headers, reader)
}
