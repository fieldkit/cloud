package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"
	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

func NewAllHandlers(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig) RecordHandler {
	return NewHandlerCollectionHandler(
		[]RecordHandler{
			handlers.NewStationModelRecordHandler(db),
			handlers.NewTsDbHandler(db, tsConfig),
		},
	)
}

type HandlerCollectionHandler struct {
	handlers []RecordHandler
}

func NewHandlerCollectionHandler(handlers []RecordHandler) *HandlerCollectionHandler {
	return &HandlerCollectionHandler{
		handlers: handlers,
	}
}

func (v *HandlerCollectionHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	for _, h := range v.handlers {
		if err := h.OnMeta(ctx, p, r, meta); err != nil {
			return err
		}
	}
	return nil
}

func (v *HandlerCollectionHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	for _, h := range v.handlers {
		if err := h.OnData(ctx, p, r, db, meta); err != nil {
			return err
		}
	}
	return nil
}

func (v *HandlerCollectionHandler) OnDone(ctx context.Context) error {
	for _, h := range v.handlers {
		if err := h.OnDone(ctx); err != nil {
			return err
		}
	}
	return nil
}
