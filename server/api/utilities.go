package api

import (
	"time"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getBucketAndKey(s3Url string) (bucket, key string, err error) {
	u, err := url.Parse(s3Url)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(u.Host, ".")

	return parts[0], u.Path[1:], nil
}

func SignS3URL(svc *s3.S3, url string) (signed string, err error) {
	bucket, key, err := getBucketAndKey(url)
	if err != nil {
		return "", err
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	signed, err = req.Presign(1 * time.Hour)
	if err != nil {
		return "", err
	}

	return 
}
