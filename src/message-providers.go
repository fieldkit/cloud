package main

import (
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	FormUrlEncodedMimeType       = "x-www-form-urlencoded"
	RockBlockProviderName        = "ROCKBLOCK"
	RockBlockFormSerial          = "serial"
	RockBlockFormData            = "data"
	RockBlockFormDeviceType      = "device_type"
	RockBlockFormDeviceTypeValue = "ROCKBLOCK"

	TwilioProviderName = "TWILIO"
	TwilioFormFrom     = "From"
	TwilioFormData     = "Data"
	TwilioFormSmsSid   = "SmsSid"

	ParticleProviderName          = "PARTICLE"
	ParticleFormCoreId            = "coreid"
	ParticleFormData              = "data"
	ParticleFormPublishedAt       = "published_at"
	ParticleFormPublishedAtLayout = "2006-01-02T15:04:05Z"

	LegacyStationNamePattern = "[A-Z][A-Z0-9]"
	LegacyMessageTypePattern = "(ST|WE|LO|AT|SO)"
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
	ProcessMessage(rmd *RawMessageData) (pm *ProcessedMessage, err error)
}

type MessageProviderBase struct {
	Form url.Values
}

type ParticleMessageProvider struct {
	MessageProviderBase
}

func (i *ParticleMessageProvider) ProcessMessage(rmd *RawMessageData) (pm *ProcessedMessage, err error) {
	coreId := strings.TrimSpace(i.Form.Get(ParticleFormCoreId))
	trimmed := strings.TrimSpace(i.Form.Get(ParticleFormData))
	fields := strings.Split(trimmed, ",")

	publishedAt, err := time.Parse(ParticleFormPublishedAtLayout, strings.TrimSpace(i.Form.Get(ParticleFormPublishedAt)))
	if err != nil {
		return nil, err
	}

	pm = &ProcessedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId(ParticleProviderName, coreId, ""),
		Time:        &publishedAt,
		ArrayValues: fields,
	}

	return
}

type TwilioMessageProvider struct {
	MessageProviderBase
}

func (i *TwilioMessageProvider) ProcessMessage(rmd *RawMessageData) (pm *ProcessedMessage, err error) {
	return normalizeCommaSeparated(TwilioProviderName, i.Form.Get(TwilioFormFrom), rmd, i.Form.Get(TwilioFormData))
}

type RockBlockMessageProvider struct {
	MessageProviderBase
}

// This will go away eventually.
var LegacyMessageTypeRe = regexp.MustCompile(LegacyMessageTypePattern)
var LegacyStationNameRe = regexp.MustCompile(LegacyStationNamePattern)

func normalizeCommaSeparated(provider string, schemaPrefix string, rmd *RawMessageData, text string) (pm *ProcessedMessage, err error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("%s(Empty message)", provider)
	}

	trimmed := strings.TrimSpace(text)
	fields := strings.Split(trimmed, ",")
	if len(fields) < 2 {
		return nil, fmt.Errorf("%s(Not enough fields: '%s')", provider, text)
	}

	maybeTime := fields[0]
	maybeStationName := fields[1]
	maybeMessageType := fields[2]

	if _, timeGood := strconv.ParseInt(maybeTime, 10, 32); timeGood != nil {
		return nil, fmt.Errorf("%s(Invalid legacy CSV time: '%s' from '%s')", provider, maybeTime, text)
	}
	if !LegacyMessageTypeRe.MatchString(maybeMessageType) {
		return nil, fmt.Errorf("%s(Invalid legacy CSV type: '%s' from '%s')", provider, maybeMessageType, text)
	}
	if !LegacyStationNameRe.MatchString(maybeStationName) || !LegacyMessageTypeRe.MatchString(maybeMessageType) {
		return nil, fmt.Errorf("%s(Invalid legacy CSV name: '%s' from '%s')", provider, maybeStationName, text)
	}

	pm = &ProcessedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId(provider, schemaPrefix, maybeMessageType),
		ArrayValues: fields,
	}

	return
}

func normalizeBinary(provider string, schemaPrefix string, rmd *RawMessageData, bytes []byte) (pm *ProcessedMessage, err error) {
	// This is a protobuf message or some other kind of similar low level binary.
	buffer := proto.NewBuffer(bytes)

	// NOTE: Right now we're only dealing with the binary format we
	// came up with during the 'NatGeo demo phase' Eventually this
	// will be a property message we can just slurp up. Though, maybe this
	// will be a great RB message going forward?
	id, err := buffer.DecodeVarint()
	if err != nil {
		return nil, err
	}
	unix, err := buffer.DecodeVarint()
	if err != nil {
		return nil, err
	}

	// HACK
	values := make([]string, 0)

	for {
		f64, err := buffer.DecodeFixed32()
		if err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			return nil, err
		}

		value := math.Float32frombits(uint32(f64))
		values = append(values, fmt.Sprintf("%f", value))
	}

	time := time.Unix(int64(unix), 0)

	pm = &ProcessedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId(provider, schemaPrefix, strconv.Itoa(int(id))),
		Time:        &time,
		ArrayValues: values,
	}

	return
}

// TODO: We should annotate incoming messages with information about their failure for logging/debugging.
func (i *RockBlockMessageProvider) ProcessMessage(rmd *RawMessageData) (pm *ProcessedMessage, err error) {
	serial := i.Form.Get(RockBlockFormSerial)
	if len(serial) == 0 {
		return
	}

	data := i.Form.Get(RockBlockFormData)
	if len(data) == 0 {
		return
	}

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return
	}

	if unicode.IsPrint(rune(bytes[0])) {
		return normalizeCommaSeparated(RockBlockProviderName, serial, rmd, string(bytes))
	}

	return normalizeBinary(RockBlockProviderName, serial, rmd, bytes)
}

func IdentifyMessageProvider(rmd *RawMessageData) (t MessageProvider, err error) {
	if strings.Contains(rmd.Params.Headers.ContentType, FormUrlEncodedMimeType) {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, nil
		}

		if form.Get(RockBlockFormDeviceType) == RockBlockFormDeviceTypeValue {
			t = &RockBlockMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get(ParticleFormCoreId) != "" {
			t = &ParticleMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get(TwilioFormSmsSid) != "" {
			t = &TwilioMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		}
	}

	return t, nil

}
