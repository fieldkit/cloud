package ingestion

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LegacyStationNamePattern = "[A-Z][A-Z0-9]"
	LegacyMessageTypePattern = "(ST|WE|LO|AT|SO)"
)

// This will go away eventually.
var LegacyMessageTypeRe = regexp.MustCompile(LegacyMessageTypePattern)
var LegacyStationNameRe = regexp.MustCompile(LegacyStationNamePattern)

func normalizeCommaSeparated(provider, providerId string, raw *RawMessage, text string) (pm *ProcessedMessage, err error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("%s(Empty message)", provider)
	}

	trimmed := strings.TrimSpace(text)
	fields := strings.Split(trimmed, ",")
	if len(fields) < 2 {
		return nil, fmt.Errorf("%s(Not enough fields: '%s')", provider, StripNewLines(text))
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
		MessageId:   MessageId(raw.RequestId),
		SchemaId:    NewSchemaId(NewProviderDeviceId(provider, providerId), maybeMessageType),
		ArrayValues: fields,
	}

	return
}

func normalizeBinary(provider, providerId string, raw *RawMessage, bytes []byte) (pm *ProcessedMessage, err error) {
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
		MessageId:   MessageId(raw.RequestId),
		SchemaId:    NewSchemaId(NewProviderDeviceId(provider, providerId), strconv.Itoa(int(id))),
		Time:        &time,
		ArrayValues: values,
	}

	return
}
