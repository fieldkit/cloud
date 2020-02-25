package backend

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/common"
)

func SignS3URL(svc *s3.S3, url string) (signed string, err error) {
	bak, err := common.GetBucketAndKey(url)
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
