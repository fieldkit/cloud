package common

import (
	"fmt"
	"net/url"
	"strings"
)

type BucketAndKey struct {
	Bucket string
	Key    string
}

func GetBucketAndKey(s3Url string) (*BucketAndKey, error) {
	if len(s3Url) == 0 {
		return nil, fmt.Errorf("bucket and key: url is required")
	}

	u, err := url.Parse(s3Url)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("malformed s3 url: %s", s3Url)
	}

	parts := strings.Split(u.Host, ".")

	return &BucketAndKey{
		Bucket: parts[0],
		Key:    u.Path[1:],
	}, nil
}
