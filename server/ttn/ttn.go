package ttn

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type ThingsNetworkMessage struct {
	ID        int32     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Headers   *string   `db:"headers"`
	Body      []byte    `db:"body"`
}

type ParsedMessage struct {
	original   *ThingsNetworkMessage
	deviceID   []byte
	deviceName string
	data       map[string]interface{}
	receivedAt time.Time
}

type LatestThingsNetworkApplicationIDs struct {
	ApplicationID string `json:"application_id"`
}

type LatestThingsNetworkDeviceIDs struct {
	ApplicationIDs LatestThingsNetworkApplicationIDs `json:"application_ids"`
	DeviceID       string                            `json:"device_id"`
	DeviceEUI      string                            `json:"dev_eui"`
	DeviceAddress  string                            `json:"dev_addr"`
}

type LatestThingsNetworkUplinkMessage struct {
	Port            int32                  `json:"f_port"`
	Counter         int32                  `json:"f_cnt"`
	DecodedPayload  map[string]interface{} `json:"decoded_payload"`
	ReceivedAt      string                 `json:"received_at"`
	ConsumedAirtime string                 `json:"consumed_airtime"`
}

type LatestThingsNetwork struct {
	EndDeviceIDs  LatestThingsNetworkDeviceIDs     `json:"end_device_ids"`
	UplinkMessage LatestThingsNetworkUplinkMessage `json:"uplink_message"`
}

func (m *ThingsNetworkMessage) Parse(ctx context.Context) (p *ParsedMessage, err error) {
	// My plan here is to grow the number of structs that we're attempting to
	// parse against as a means for versioning these. So, eventually there may
	// be a loop here or a series of Unmarshal attempts.
	raw := LatestThingsNetwork{}
	if err := json.Unmarshal(m.Body, &raw); err != nil {
		return nil, err
	}

	deviceID, err := hex.DecodeString(raw.EndDeviceIDs.DeviceEUI)
	if err != nil {
		return nil, fmt.Errorf("malformed device eui: %s", raw.EndDeviceIDs.DeviceEUI)
	}

	if raw.EndDeviceIDs.DeviceID == "" {
		return nil, fmt.Errorf("empty device id (station name)")
	}

	// "received_at": "2021-08-18T22:17:18.803716242Z",
	receivedAt, err := time.Parse("2006-01-02T15:04:05.999999999Z", raw.UplinkMessage.ReceivedAt)
	if err != nil {
		return nil, fmt.Errorf("malformed received at (%s)", raw.UplinkMessage.ReceivedAt)
	}

	return &ParsedMessage{
		original:   m,
		deviceID:   deviceID,
		deviceName: raw.EndDeviceIDs.DeviceID,
		data:       raw.UplinkMessage.DecodedPayload,
		receivedAt: receivedAt,
	}, nil
}
