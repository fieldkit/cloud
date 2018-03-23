package formatting

import (
	"github.com/fieldkit/cloud/server/backend/ingestion"
)

const (
	TwilioProviderName = "TWILIO"
	TwilioFormFrom     = "From"
	TwilioFormData     = "Data"
	TwilioFormSmsSid   = "SmsSid"
)

type TwilioMessageProvider struct {
	MessageProviderBase
}

func (i *TwilioMessageProvider) CanFormatMessage(raw *RawMessage) bool {
	return false
}

func (i *TwilioMessageProvider) FormatMessage(raw *RawMessage) (fm *ingestion.FormattedMessage, err error) {
	return normalizeCommaSeparated(TwilioProviderName, raw.Form.Get(TwilioFormFrom), raw, raw.Form.Get(TwilioFormData))
}
