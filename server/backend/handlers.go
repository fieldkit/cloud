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

type muxRecordHandler struct {
	handlers []RecordHandler
}

func NewMuxRecordHandler(handlers []RecordHandler) *muxRecordHandler {
	return &muxRecordHandler{
		handlers: handlers,
	}
}

func (h *muxRecordHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	for _, h := range h.handlers {
		if err := h.OnMeta(ctx, p, r, db); err != nil {
			return err
		}
	}
	return nil
}

func (h *muxRecordHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	for _, h := range h.handlers {
		if err := h.OnData(ctx, p, r, db, meta); err != nil {
			return err
		}
	}
	return nil
}

func (h *muxRecordHandler) OnDone(ctx context.Context) error {
	for _, h := range h.handlers {
		if err := h.OnDone(ctx); err != nil {
			return err
		}
	}
	return nil
}
