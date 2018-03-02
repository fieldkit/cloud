package ingestion

import (
	"fmt"
	"time"
)

const (
	FormUrlEncodedMimeType = "x-www-form-urlencoded"
	JsonMimeType           = "json"
)

type MessageId string

type DeviceId string

func NewDeviceId(id string) DeviceId {
	return DeviceId(id)
}

func NewProviderDeviceId(provider, id string) DeviceId {
	return DeviceId(fmt.Sprintf("%s-%s", provider, id))
}

func (s DeviceId) String() string {
	return fmt.Sprintf("%s", string(s))
}

type SchemaId struct {
	Device DeviceId
	Stream string
}

func (s SchemaId) String() string {
	if s.Stream == "" {
		return s.Device.String()
	}
	return fmt.Sprintf("%s-%s", s.Device.String(), s.Stream)
}

func NewSchemaId(deviceId DeviceId, stream string) SchemaId {
	return SchemaId{
		Device: deviceId,
		Stream: stream,
	}
}

type ProcessedMessage struct {
	MessageId   MessageId
	SchemaId    SchemaId
	Time        *time.Time
	Location    []float64
	Fixed       bool
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
