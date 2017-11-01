package main

import (
	"fmt"
	"time"
)

const (
	FormUrlEncodedMimeType = "x-www-form-urlencoded"
)

type MessageId string

type SchemaId string

func MakeSchemaId(provider string, device string, stream string) SchemaId {
	if len(stream) > 0 {
		return SchemaId(fmt.Sprintf("%s-%s-%s", provider, device, stream))
	}
	return SchemaId(fmt.Sprintf("%s-%s", provider, device))
}

type ProcessedMessage struct {
	MessageId   MessageId
	SchemaId    SchemaId
	Time        *time.Time
	ArrayValues []string
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
}

func IdentifyMessageProvider(raw *RawMessage) (t MessageProvider, err error) {
	for _, provider := range AllProviders {
		if provider.CanProcessMessage(raw) {
			return provider, nil
		}
	}
	return
}
