package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type RawMessageRow struct {
	SqsId string
	Data  string
	Time  uint64
}

type HandlerFunc func(raw *RawMessage) error

func (f HandlerFunc) HandleMessage(raw *RawMessage) error {
	return f(raw)
}

type Handler interface {
	HandleMessage(raw *RawMessage) error
}

type RawMessage struct {
	Row  *RawMessageRow
	Data *SqsMessage
	Form *url.Values
}

func CreateRawMessageFromRow(row *RawMessageRow) (raw *RawMessage, err error) {
	rmd := SqsMessage{}
	err = json.Unmarshal([]byte(row.Data), &rmd)
	if err != nil {
		return nil, fmt.Errorf("Malformed RawMessage: %s", row.Data)
	}

	if rmd.Context.RequestId == "" {
		return nil, fmt.Errorf("Malformed RawMessage: Context missing RequestId.")
	}

	if strings.Contains(rmd.Params.Headers.ContentType, FormUrlEncodedMimeType) {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, fmt.Errorf("Malformed RawMessage: RawBody invalid for %s.", rmd.Params.Headers.ContentType)
		}

		raw = &RawMessage{
			Row:  row,
			Data: &rmd,
			Form: &form,
		}

		return raw, nil
	}

	if strings.Contains(rmd.Params.Headers.ContentType, JsonMimeType) {
		raw = &RawMessage{
			Row:  row,
			Data: &rmd,
			Form: &url.Values{},
		}

		return raw, nil
	}

	return nil, fmt.Errorf("Unexpected ContentType: %s", rmd.Params.Headers.ContentType)
}
