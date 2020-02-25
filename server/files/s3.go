package files

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/fieldkit/cloud/server/common"
)

type S3FileArchive struct {
	session    *session.Session
	bucketName string
}

func NewS3FileArchive(session *session.Session, bucketName string) *S3FileArchive {
	return &S3FileArchive{
		session:    session,
		bucketName: bucketName,
	}
}

func (a *S3FileArchive) Archive(ctx context.Context, meta *FileMeta, reader io.Reader) (*ArchivedFile, error) {
	id := uuid.Must(uuid.NewRandom())

	log := Logger(ctx).Sugar()

	uploader := s3manager.NewUploader(a.session)

	metadata := make(map[string]*string)
	metadata[common.FkDeviceIdHeaderName] = aws.String(hex.EncodeToString(meta.DeviceID))
	metadata[common.FkGenerationHeaderName] = aws.String(hex.EncodeToString(meta.Generation))
	metadata[common.FkBlocksIdHeaderName] = aws.String(fmt.Sprintf("%v", meta.Blocks))
	metadata[common.FkFlagsIdHeaderName] = aws.String(fmt.Sprintf("%v", meta.Flags))

	countingReader := newCountingReader(reader)

	r, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String(meta.ContentType),
		Bucket:      aws.String(a.bucketName),
		Key:         aws.String(id.String()),
		Body:        countingReader,
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return nil, err
	}

	log.Infow("saved", "url", r.Location, "bytes_read", countingReader.bytesRead)

	ss := &ArchivedFile{
		ID:        id.String(),
		URL:       r.Location,
		BytesRead: countingReader.bytesRead,
	}

	return ss, err
}
