package ingestion

const (
	TwilioProviderName = "TWILIO"
	TwilioFormFrom     = "From"
	TwilioFormData     = "Data"
	TwilioFormSmsSid   = "SmsSid"
)

type TwilioMessageProvider struct {
	MessageProviderBase
}

func (i *TwilioMessageProvider) CanProcessMessage(raw *RawMessage) bool {
	return false
}

func (i *TwilioMessageProvider) ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error) {
	return normalizeCommaSeparated(TwilioProviderName, raw.Form.Get(TwilioFormFrom), raw, raw.Form.Get(TwilioFormData))
}
