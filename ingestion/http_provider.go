package ingestion

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
	HttpProviderTokenKey        = "token"
)

type HttpJsonMessage struct {
	Location []float32         `json:"location"`
	Time     int64             `json:"time"`
	Device   string            `json:"device"`
	Stream   string            `json:"stream"`
	Values   map[string]string `json:"values"`
}

func (i *HttpMessageProvider) CanProcessMessage(raw *RawMessage) bool {
	if raw.ContentType == HttpProviderJsonContentType {
		if raw.QueryString.Get(HttpProviderTokenKey) == "" {
			return false
		}
		return true
	}
	return false
}

func (i *HttpMessageProvider) ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error) {
	if raw.ContentType == HttpProviderJsonContentType {
		message := HttpJsonMessage{}
		err = json.Unmarshal([]byte(raw.RawBody), &message)
		if err != nil {
			return nil, err
		}

		if message.Device == "" {
			return nil, fmt.Errorf("Malformed HttpJsonMessage. Device is required.")
		}

		if len(message.Location) < 2 {
			return nil, fmt.Errorf("Malformed HttpJsonMessage. Location is required.")
		}

		messageTime := time.Unix(message.Time, 0)

		pm = &ProcessedMessage{
			MessageId: MessageId(raw.RequestId),
			SchemaId:  NewSchemaId(HttpProviderName, message.Device, message.Stream),
			Time:      &messageTime,
			Location:  message.Location,
			MapValues: message.Values,
		}

		return
	}

	return nil, fmt.Errorf("Unknown ContentType: %s", raw.ContentType)
}
