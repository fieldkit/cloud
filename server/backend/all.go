package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"
	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

func NewAllHandlers(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig, publisher jobs.MessagePublisher) RecordHandler {
	return NewHandlerCollectionHandler(
		[]RecordHandler{
			handlers.NewStationModelRecordHandler(db),
			handlers.NewTsDbHandler(db, tsConfig, publisher),
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

func (v *HandlerCollectionHandler) OnMeta(ctx context.Context, provision *data.Provision, rawMeta *pb.DataRecord, meta *data.MetaRecord) error {
	for _, h := range v.handlers {
		if err := h.OnMeta(ctx, provision, rawMeta, meta); err != nil {
			return err
		}
	}
	return nil
}

func (v *HandlerCollectionHandler) OnData(ctx context.Context, provision *data.Provision, rawData *pb.DataRecord, rawMeta *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	for _, h := range v.handlers {
		if err := h.OnData(ctx, provision, rawData, rawMeta, db, meta); err != nil {
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
