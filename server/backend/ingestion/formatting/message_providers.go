package formatting

import (
	"github.com/fieldkit/cloud/server/backend/ingestion"
)

const (
	FormUrlEncodedMimeType = "x-www-form-urlencoded"
	JsonMimeType           = "json"
)

type MessageProvider interface {
	CanFormatMessage(raw *RawMessage) bool
	FormatMessage(raw *RawMessage) (pm *ingestion.FormattedMessage, err error)
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
		if provider.CanFormatMessage(raw) {
			return provider, nil
		}
	}
	return
}
