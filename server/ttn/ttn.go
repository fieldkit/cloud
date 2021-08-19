package ttn

import (
	"context"
	"encoding/json"
	"time"
)

type ThingsNetworkMessage struct {
	ID        int32     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Headers   *string   `db:"headers"`
	Body      []byte    `db:"body"`
}

type parsedMessage struct {
	original   *ThingsNetworkMessage
	deviceID   string
	deviceName string
	data       map[string]interface{}
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

func (m *ThingsNetworkMessage) Parse(ctx context.Context) (p *parsedMessage, err error) {
	// My plan here is to grow the number of structs that we're attempting to
	// parse against as a means for versioning these. So, eventually there may
	// be a loop here or a series of Unmarshal attempts.
	raw := LatestThingsNetwork{}
	if err := json.Unmarshal(m.Body, &raw); err != nil {
		return nil, err
	}

	return &parsedMessage{
		original:   m,
		deviceID:   raw.EndDeviceIDs.DeviceEUI,
		deviceName: raw.EndDeviceIDs.DeviceID,
		data:       raw.UplinkMessage.DecodedPayload,
	}, nil
}
