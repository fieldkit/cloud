package ingestion

import (
	"fmt"
	"regexp"
	"time"
)

const (
	FormUrlEncodedMimeType = "x-www-form-urlencoded"
	JsonMimeType           = "json"
)

type MessageId string

type DeviceId struct {
	Provider string
	Id       string
}

func (s DeviceId) ToString() string {
	return fmt.Sprintf("%s-%s", s.Provider, s.Id)
}

type SchemaId struct {
	Device DeviceId
	Stream string
}

func (s SchemaId) ToString() string {
	if s.Stream == "" {
		return s.Device.ToString()
	}
	return fmt.Sprintf("%s-%s", s.Device.ToString(), s.Stream)
}

func NewSchemaId(provider string, device string, stream string) SchemaId {
	return SchemaId{
		Device: DeviceId{
			Provider: provider,
			Id:       device,
		},
		Stream: stream,
	}
}

func (s SchemaId) Matches(o SchemaId) bool {
	if s.Device.Provider != o.Device.Provider {
		return false
	}
	re := regexp.MustCompile(s.Device.Id)
	if !re.MatchString(o.Device.Id) {
		return false
	}

	if s.Stream != "" && s.Stream != o.Stream {
		return false
	}

	return true
}

type ProcessedMessage struct {
	MessageId   MessageId
	SchemaId    SchemaId
	Time        *time.Time
	Location    []float32
	ArrayValues []string
	MapValues   map[string]string
}

type MessageProvider interface {
	CanProcessMessage(raw *RawMessage) bool
	ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error)
}

type MessageProviderBase struct {
}

var AllProviders = []MessageProvider{
	&RockBlockMessageProvider{},
	&TwilioMessageProvider{},
	&ParticleMessageProvider{},
	&HttpMessageProvider{},
}

func IdentifyMessageProvider(raw *RawMessage) (t MessageProvider, err error) {
	for _, provider := range AllProviders {
		if provider.CanProcessMessage(raw) {
			return provider, nil
		}
	}
	return
}
