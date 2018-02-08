package ingestion

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type RawMessageRow struct {
	Id   string
	Data string
	Time uint64
}

type HandlerFunc func(raw *RawMessage) error

func (f HandlerFunc) HandleMessage(raw *RawMessage) error {
	return f(raw)
}

type RawMessageHandler interface {
	HandleMessage(raw *RawMessage) error
}

type RawMessage struct {
	RequestId   string
	RawBody     []byte
	ContentType string
	QueryString *url.Values
	Form        *url.Values
}

func CreateRawMessageFromRow(row *RawMessageRow) (raw *RawMessage, err error) {
	rmd := SqsMessage{}
	err = json.Unmarshal([]byte(row.Data), &rmd)
	if err != nil {
		return nil, fmt.Errorf("Malformed RawMessage: %s (%v)", row.Data, err)
	}

	if rmd.Context.RequestId == "" {
		return nil, fmt.Errorf("Malformed RawMessage: Context missing RequestId.")
	}

	form := &url.Values{}
	queryString := &url.Values{}
	for k, v := range rmd.Params.QueryString {
		queryString.Add(k, v)
	}

	if strings.Contains(rmd.Params.Headers.ContentType, FormUrlEncodedMimeType) {
		parsedForm, err := url.ParseQuery(string(rmd.RawBody))
		if err != nil {
			return nil, fmt.Errorf("Malformed RawMessage: RawBody invalid for %s.", rmd.Params.Headers.ContentType)
		}
		form = &parsedForm
	} else if strings.Contains(rmd.Params.Headers.ContentType, JsonMimeType) {

	} else {
		return nil, fmt.Errorf("Unexpected ContentType: %s", rmd.Params.Headers.ContentType)
	}

	raw = &RawMessage{
		RequestId:   rmd.Context.RequestId,
		RawBody:     rmd.RawBody,
		ContentType: rmd.Params.Headers.ContentType,
		QueryString: queryString,
		Form:        form,
	}

	return
}
