package webhook

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

const (
	CoverageSchemaID = 9
)

type MapsRepository struct {
	db *sqlxcache.DB
}

func NewMapsRepository(db *sqlxcache.DB) (rr *MapsRepository) {
	return &MapsRepository{db: db}
}

func (r *MapsRepository) QueryCoverageMap(ctx context.Context) ([]*data.CoveragePoint, error) {
	messages := []*WebHookMessage{}
	if err := r.db.SelectContext(ctx, &messages,
		`SELECT id, created_at, schema_id, headers, body
		FROM fieldkit.ttn_messages
		WHERE schema_id = $1
		ORDER BY created_at`, CoverageSchemaID); err != nil {
		return nil, err
	}

	coverage := make([]*data.CoveragePoint, 0)

	for _, m := range messages {
		ttnMessage := TheThingsNetworkV3{}
		if err := json.Unmarshal(m.Body, &ttnMessage); err != nil {
			return nil, err
		}

		receivedAt, err := time.Parse("2006-01-02T15:04:05.999999999Z", ttnMessage.Uplink.ReceivedAt)
		if err != nil {
			return nil, err
		}

		coverage = append(coverage, &data.CoveragePoint{
			Time:        receivedAt,
			Coordinates: []float64{ttnMessage.Uplink.DecodedPayload.Longitude, ttnMessage.Uplink.DecodedPayload.Latitude, ttnMessage.Uplink.DecodedPayload.Altitude},
			Satellites:  ttnMessage.Uplink.DecodedPayload.Satellites,
			HDOP:        ttnMessage.Uplink.DecodedPayload.HDOP,
		})
	}

	return coverage, nil
}

type TheThingsNetworkV3 struct {
	Uplink TheThingsNetworkUplink `json:"uplink_message"`
}

type TheThingsNetworkUplink struct {
	ReceivedAt     string                         `json:"received_at"`
	DecodedPayload TheThingsNetworkDecodedPayload `json:"decoded_payload"`
}

type TheThingsNetworkDecodedPayload struct {
	Altitude   float64 `json:"altitude"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Satellites int32   `json:"sats"`
	HDOP       float64 `json:"hdop"`
}
