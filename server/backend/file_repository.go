package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileRepository struct {
	Session *session.Session
	Bucket  string
}

type FileInfo struct {
	Key          string
	URL          string
	LastModified time.Time
	Size         int64
	Meta         map[string]*string
}

func NewFileRepository(s *session.Session, bucket string) (fr *FileRepository, err error) {
	fr = &FileRepository{
		Session: s,
		Bucket:  bucket,
	}

	return
}

func (fr *FileRepository) Info(ctx context.Context, key string) (fi *FileInfo, err error) {
	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(fr.Bucket),
		Key:    aws.String(key),
	}

	svc := s3.New(fr.Session)

	obj, err := svc.HeadObject(hoi)
	if err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok {
			switch aerr.StatusCode() {
			case 404:
				return nil, nil
			}
		}

		return nil, fmt.Errorf("Error calling HeadObject(%v): %v", key, err)
	}

	fi = &FileInfo{
		Key:          key,
		URL:          fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fr.Bucket, key),
		LastModified: *obj.LastModified,
		Size:         *obj.ContentLength,
		Meta:         obj.Metadata,
	}

	return
}

func fixMeta(m map[string]*string) map[string]*string {
	newM := make(map[string]*string)
	newM["Fk-DeviceId"] = m["Fk-Deviceid"]
	newM["Fk-FileId"] = m["Fk-Fileid"]
	newM["Fk-Build"] = m["Fk-Build"]
	newM["Fk-FileName"] = m["Fk-Filename"]
	newM["Fk-Version"] = m["Fk-Version"]
	return newM
}
