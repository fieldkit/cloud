package ingestion

import (
	"net/url"
)

type RawMessage struct {
	RequestId   string
	RawBody     []byte
	ContentType string
	QueryString *url.Values
	Form        *url.Values
}
