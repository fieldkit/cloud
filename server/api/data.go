package api

import (
	"context"
	"encoding/hex"

	"goa.design/goa/v3/security"

	"github.com/fieldkit/cloud/server/data"

	datas "github.com/fieldkit/cloud/server/api/gen/data"
)

type DataService struct {
	options *ControllerOptions
}

func NewDataService(ctx context.Context, options *ControllerOptions) *DataService {
	return &DataService{
		options: options,
	}
}

type ProvisionSummaryRow struct {
	GenerationID []byte          `db:"generation"`
	Type         string          `db:"type"`
	Blocks       data.Int64Range `db:"blocks"`
}

type BlocksSummaryRow struct {
	Size   int64
	Blocks data.Int64Range
}

func (s *DataService) DeviceSummary(ctx context.Context, payload *datas.DeviceSummaryPayload) (*datas.DeviceDataSummaryResponse, error) {
	log := Logger(ctx).Sugar()

	deviceIdBytes, err := data.DecodeBinaryString(payload.DeviceID)
	if err != nil {
		return nil, err
	}

	p, err := NewPermissions(ctx, s.options).ForStationByDeviceID(deviceIdBytes)
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	log.Infow("summary", "device_id", deviceIdBytes)

	provisions := make([]*data.Provision, 0)
	if err := s.options.Database.SelectContext(ctx, &provisions, `SELECT * FROM fieldkit.provision WHERE (device_id = $1) ORDER BY updated DESC`, deviceIdBytes); err != nil {
		return nil, err
	}

	var rows = make([]ProvisionSummaryRow, 0)
	if err := s.options.Database.SelectContext(ctx, &rows, `SELECT generation, type, range_merge(blocks) AS blocks FROM fieldkit.ingestion WHERE (device_id = $1) GROUP BY type, generation`, deviceIdBytes); err != nil {
		return nil, err
	}

	var blockSummaries = make(map[string]map[string]BlocksSummaryRow)
	for _, p := range provisions {
		blockSummaries[hex.EncodeToString(p.GenerationID)] = make(map[string]BlocksSummaryRow)
	}

	for _, row := range rows {
		g := hex.EncodeToString(row.GenerationID)
		blockSummaries[g][row.Type] = BlocksSummaryRow{
			Blocks: row.Blocks,
		}
	}

	provisionVms := make([]*datas.DeviceProvisionSummary, len(provisions))
	for i, p := range provisions {
		g := hex.EncodeToString(p.GenerationID)
		byType := blockSummaries[g]
		provisionVms[i] = &datas.DeviceProvisionSummary{
			Generation: g,
			Created:    p.Created.Unix(),
			Updated:    p.Updated.Unix(),
			Meta: &datas.DeviceMetaSummary{
				Size:  int64(0),
				First: int64(byType[data.MetaTypeName].Blocks[0]),
				Last:  int64(byType[data.MetaTypeName].Blocks[1]),
			},
			Data: &datas.DeviceDataSummary{
				Size:  int64(0),
				First: int64(byType[data.DataTypeName].Blocks[0]),
				Last:  int64(byType[data.DataTypeName].Blocks[1]),
			},
		}
	}

	return &datas.DeviceDataSummaryResponse{
		Provisions: provisionVms,
	}, nil
}

func (s *DataService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return datas.NotFound(m) },
		Unauthorized: func(m string) error { return datas.Unauthorized(m) },
		Forbidden:    func(m string) error { return datas.Forbidden(m) },
	})
}
