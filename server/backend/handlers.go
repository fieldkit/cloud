package backend

import (
	"context"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type RecordHandler interface {
	OnMeta(ctx context.Context, r *pb.DataRecord, db *data.MetaRecord, p *data.Provision) error
	OnData(ctx context.Context, r *pb.DataRecord, db *data.DataRecord, p *data.Provision) error
}

type noopRecordHandler struct {
}

func NewNoopRecordHandler() *noopRecordHandler {
	return &noopRecordHandler{}
}

func (h *noopRecordHandler) OnMeta(ctx context.Context, r *pb.DataRecord, db *data.MetaRecord, p *data.Provision) error {
	return nil
}

func (h *noopRecordHandler) OnData(ctx context.Context, r *pb.DataRecord, db *data.DataRecord, p *data.Provision) error {
	return nil
}
