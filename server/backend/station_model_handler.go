package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type stationModelRecordHandler struct {
	database *sqlxcache.DB
}

func NewStationModelRecordHandler(database *sqlxcache.DB) *stationModelRecordHandler {
	return &stationModelRecordHandler{
		database: database,
	}
}

func (h *stationModelRecordHandler) OnMeta(ctx context.Context, r *pb.DataRecord, db *data.MetaRecord, p *data.Provision) error {
	return nil
}

func (h *stationModelRecordHandler) OnData(ctx context.Context, r *pb.DataRecord, db *data.DataRecord, p *data.Provision) error {
	return nil
}
