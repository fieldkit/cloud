package main

import ()

type SqsMessageHeaders struct {
	UserAgent   string `json:"User-Agent"`
	ContentType string `json:"Content-Type"`
}

type SqsMessageParams struct {
	Headers     SqsMessageHeaders `json:"header"`
	QueryString map[string]string `json:"querystring"`
}

type SqsMessageContext struct {
	UserAgent string `json:"user-agent"`
	RequestId string `json:"request-id"`
}

type SqsMessage struct {
	RawBody string            `json:"body-raw"`
	Params  SqsMessageParams  `json:"params"`
	Context SqsMessageContext `json:"context"`
}
