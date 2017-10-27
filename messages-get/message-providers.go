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

type MessageId string

type SchemaId string

func MakeSchemaId(provider string, device string, stream string) SchemaId {
	if len(stream) > 0 {
		return SchemaId(fmt.Sprintf("%s-%s-%s", provider, device, stream))
	}
	return SchemaId(fmt.Sprintf("%s-%s", provider, device))
}

type NormalizedMessage struct {
	MessageId   MessageId
	SchemaId    SchemaId
	Time        time.Time
	ArrayValues []string
}

type MessageProvider interface {
	NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error)
}

type MessageProviderBase struct {
	Form url.Values
}

type ParticleMessageProvider struct {
	MessageProviderBase
}

func (i *ParticleMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	coreId := strings.TrimSpace(i.Form.Get("coreid"))
	trimmed := strings.TrimSpace(i.Form.Get("data"))
	fields := strings.Split(trimmed, ",")

	publishedAt, err := time.Parse("2006-01-02T15:04:05Z", strings.TrimSpace(i.Form.Get("published_at")))
	if err != nil {
		return nil, err
	}

	nm = &NormalizedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId("PARTICLE", coreId, ""),
		Time:        publishedAt,
		ArrayValues: fields,
	}

	return
}

type TwilioMessageProvider struct {
	MessageProviderBase
}

func (i *TwilioMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	return normalizeCommaSeparated("TWILIO", i.Form.Get("From"), rmd, i.Form.Get("Body"))
}

type RockBlockMessageProvider struct {
	MessageProviderBase
}

func normalizeCommaSeparated(provider string, schemaPrefix string, rmd *RawMessageData, text string) (nm *NormalizedMessage, err error) {
	stationNameRe := regexp.MustCompile("[A-Z][A-Z]")
	trimmed := strings.TrimSpace(text)
	fields := strings.Split(trimmed, ",")
	if len(fields) < 2 {
		return nil, fmt.Errorf("Not enough fields in comma separated message.")
	}
	maybeStationName := fields[2]
	if !stationNameRe.MatchString(maybeStationName) {
		return nil, fmt.Errorf("Invalid name: %s", maybeStationName)
	}

	nm = &NormalizedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId(provider, schemaPrefix, maybeStationName),
		ArrayValues: fields,
	}

	return
}

func normalizeBinary(provider string, schemaPrefix string, rmd *RawMessageData, bytes []byte) (nm *NormalizedMessage, err error) {
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

	nm = &NormalizedMessage{
		MessageId:   MessageId(rmd.Context.RequestId),
		SchemaId:    MakeSchemaId(provider, schemaPrefix, strconv.Itoa(int(id))),
		Time:        time.Unix(int64(unix), 0),
		ArrayValues: values,
	}

	return
}

// TODO: We should annotate incoming messages with information about their failure for logging/debugging.
func (i *RockBlockMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	serial := i.Form.Get("serial")
	if len(serial) == 0 {
		return
	}

	data := i.Form.Get("data")
	if len(data) == 0 {
		return
	}

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return
	}

	if unicode.IsPrint(rune(bytes[0])) {
		return normalizeCommaSeparated("ROCKBLOCK", serial, rmd, string(bytes))
	}

	return normalizeBinary("ROCKBLOCK", serial, rmd, bytes)
}

func IdentifyMessageProvider(rmd *RawMessageData) (t MessageProvider, err error) {
	if strings.Contains(rmd.Params.Headers.ContentType, "x-www-form-urlencoded") {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, nil
		}

		if form.Get("device_type") == "ROCKBLOCK" {
			t = &RockBlockMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get("coreid") != "" {
			t = &ParticleMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get("SmsSid") != "" {
			t = &TwilioMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		}
	}

	return t, nil

}
