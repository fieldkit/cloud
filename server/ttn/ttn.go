package ttn

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type ThingsNetworkMessage struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	SchemaID  *int32    `db:"schema_id"`
	Headers   *string   `db:"headers"`
	Body      []byte    `db:"body"`
}

type ParsedMessage struct {
	original   *ThingsNetworkMessage
	deviceID   []byte
	deviceName string
	data       map[string]interface{}
	receivedAt time.Time
	schema     *ThingsNetworkSchema
	schemaID   int32
	ownerID    int32
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

func (m *ThingsNetworkMessage) Parse(ctx context.Context, schemas map[int32]*ThingsNetworkSchemaRegistration) (p *ParsedMessage, err error) {
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

	receivedAt, err := time.Parse("2006-01-02T15:04:05.999999999Z", raw.UplinkMessage.ReceivedAt)
	if err != nil {
		return nil, fmt.Errorf("malformed received at (%s)", raw.UplinkMessage.ReceivedAt)
	}

	if m.SchemaID == nil {
		return nil, fmt.Errorf("missing schema")
	}

	schemaRegistration := schemas[*m.SchemaID]
	schema, err := schemaRegistration.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %s", err)

	}

	return &ParsedMessage{
		original:   m,
		deviceID:   deviceID,
		deviceName: raw.EndDeviceIDs.DeviceID,
		data:       raw.UplinkMessage.DecodedPayload,
		receivedAt: receivedAt,
		ownerID:    schemaRegistration.OwnerID,
		schemaID:   schemaRegistration.ID,
		schema:     schema,
	}, nil
}

type ThingsNetworkSchemaSensor struct {
	DecodedKey string `json:"decoded_key"`
}

type ThingsNetworkSchemaModule struct {
	Key     string                               `json:"key"`
	Name    string                               `json:"name"`
	Sensors map[string]ThingsNetworkSchemaSensor `json:"sensors"`
}

type ThingsNetworkSchemaStation struct {
	Key     string                      `json:"key"`
	Model   string                      `json:"model"`
	Modules []ThingsNetworkSchemaModule `json:"modules"`
}

type ThingsNetworkSchema struct {
	Station ThingsNetworkSchemaStation `json:"station"`
}

type ThingsNetworkSchemaRegistration struct {
	ID              int32      `db:"id"`
	OwnerID         int32      `db:"owner_id"`
	Name            string     `db:"name"`
	Token           []byte     `db:"token"`
	Body            []byte     `db:"body"`
	ReceivedAt      *time.Time `db:"received_at"`
	ProcessedAt     *time.Time `db:"processed_at"`
	ProcessInterval *int32     `db:"process_interval"`
}

func (r *ThingsNetworkSchemaRegistration) Parse() (*ThingsNetworkSchema, error) {
	s := &ThingsNetworkSchema{}
	if err := json.Unmarshal(r.Body, s); err != nil {
		return nil, err
	}
	return s, nil
}

type ThingsNetworkMessageReceived struct {
	SchemaID  int32 `json:"ttn_schema_id"`
	MessageID int64 `json:"ttn_message_id"`
}

type ProcessSchema struct {
	SchemaID int32 `json:"schema_id"`
}
