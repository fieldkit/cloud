package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type HttpMessageProvider struct {
	MessageProviderBase
}

const (
	HttpProviderName            = "HTTP"
	HttpProviderJsonContentType = "application/vnd.fk.message+json"
	HttpProviderFormContentType = "application/vnd.fk.message+x-www-form-urlencoded"
)

type HttpJsonMessage struct {
	Location []float32         `json:"location"`
	Time     int64             `json:"time"`
	Device   string            `json:"device"`
	Stream   string            `json:"stream"`
	Values   map[string]string `json:"values"`
}

func (i *HttpMessageProvider) CanProcessMessage(raw *RawMessage) bool {
	if raw.Data.Params.Headers.ContentType == HttpProviderJsonContentType {
		if raw.Data.Params.QueryString["token"] == "" {
			return false
		}
		return true
	}
	return false
}

func (i *HttpMessageProvider) ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error) {
	if raw.Data.Params.Headers.ContentType == HttpProviderJsonContentType {
		message := HttpJsonMessage{}
		err = json.Unmarshal([]byte(raw.Data.RawBody), &message)
		if err != nil {
			return nil, err
		}

		if message.Device == "" {
			return nil, fmt.Errorf("Malformed HttpJsonMessage. Device is required.")
		}

		messageTime := time.Unix(message.Time, 0)

		pm = &ProcessedMessage{
			MessageId: MessageId(raw.Data.Context.RequestId),
			SchemaId:  NewSchemaId(HttpProviderName, message.Device, message.Stream),
			Time:      &messageTime,
			Location:  message.Location,
			MapValues: message.Values,
		}

		return
	}

	return nil, fmt.Errorf("Unknown ContentType: %s", raw.Data.Params.Headers.ContentType)
}
