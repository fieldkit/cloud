package backend

import (
	"context"
	"fmt"
	"strings"
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
	DeviceID     string
	FileTypeID   string
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

	meta := sanitizeMeta(obj.Metadata)
	deviceID := ""
	if value, ok := meta[FkDeviceIdHeaderName]; ok {
		deviceID = *value
	}
	fileTypeID := ""
	if value, ok := meta[FkFileIdHeaderName]; ok {
		fileTypeID = *value
	}

	fi = &FileInfo{
		Key:          key,
		DeviceID:     deviceID,
		FileTypeID:   fileTypeID,
		URL:          fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fr.Bucket, key),
		LastModified: *obj.LastModified,
		Size:         *obj.ContentLength,
		Meta:         meta,
	}

	return
}

func sanitizeMeta(m map[string]*string) map[string]*string {
	ci := make(map[string]*string)
	for key, value := range m {
		ci[strings.ToLower(key)] = value
	}
	newM := make(map[string]*string)
	for _, key := range []string{FkDeviceIdHeaderName, FkFileIdHeaderName, FkBuildHeaderName, FkFileNameHeaderName, FkVersionHeaderName} {
		if value, ok := ci[strings.ToLower(key)]; ok {
			newM[key] = value
		}
	}
	return newM
}
