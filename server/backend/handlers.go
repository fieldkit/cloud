package backend

import (
	"context"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type RecordHandler interface {
	OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error
	OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error
	OnDone(ctx context.Context) error
}

type noopRecordHandler struct {
}

func NewNoopRecordHandler() *noopRecordHandler {
	return &noopRecordHandler{}
}

func (h *noopRecordHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	return nil
}

func (h *noopRecordHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	return nil
}

func (h *noopRecordHandler) OnDone(ctx context.Context) error {
	return nil
}
