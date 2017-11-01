package main

import (
	"encoding/hex"
	"unicode"
)

const (
	RockBlockProviderName        = "ROCKBLOCK"
	RockBlockFormSerial          = "serial"
	RockBlockFormData            = "data"
	RockBlockFormDeviceType      = "device_type"
	RockBlockFormDeviceTypeValue = "ROCKBLOCK"
)

type RockBlockMessageProvider struct {
	MessageProviderBase
}

func (i *RockBlockMessageProvider) CanProcessMessage(raw *RawMessage) bool {
	return raw.Form.Get(RockBlockFormDeviceType) == RockBlockFormDeviceTypeValue
}

// TODO: We should annotate incoming messages with information about their failure for logging/debugging.
func (i *RockBlockMessageProvider) ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error) {
	serial := raw.Form.Get(RockBlockFormSerial)
	if len(serial) == 0 {
		return
	}

	data := raw.Form.Get(RockBlockFormData)
	if len(data) == 0 {
		return
	}

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return
	}

	if unicode.IsPrint(rune(bytes[0])) {
		return normalizeCommaSeparated(RockBlockProviderName, serial, raw, string(bytes))
	}

	return normalizeBinary(RockBlockProviderName, serial, raw, bytes)
}
