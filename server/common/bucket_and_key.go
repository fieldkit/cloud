package common

import (
	"net/url"
	"strings"
)

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
