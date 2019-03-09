package api

import (
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func SignS3URL(svc *s3.S3, url string) (signed string, err error) {
	bak, err := GetBucketAndKey(url)
	if err != nil {
		return "", err
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bak.Bucket),
		Key:    aws.String(bak.Key),
	})

	signed, err = req.Presign(1 * time.Hour)
	if err != nil {
		return "", err
	}

	return
}

type BucketAndKey struct {
	Bucket string
	Key    string
}

func GetBucketAndKey(s3Url string) (*BucketAndKey, error) {
	u, err := url.Parse(s3Url)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(u.Host, ".")

	return &BucketAndKey{
		Bucket: parts[0],
		Key:    u.Path[1:],
	}, nil
}
