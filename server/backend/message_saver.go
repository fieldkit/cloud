package backend

import (
	"context"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type FormattedMessageSaver struct {
	SchemaApplier *ingestion.SchemaApplier
	Repository    *ingestion.Repository
	Resolver      *ingestion.Resolver
	RecordAdder   *ingestion.RecordAdder
	Changes       map[int64]*ingestion.RecordChange
}

func NewFormattedMessageSaver(b *Backend) *FormattedMessageSaver {
	r := ingestion.NewRepository(b.db)

	return &FormattedMessageSaver{
		Changes:       make(map[int64]*ingestion.RecordChange),
		Repository:    r,
		SchemaApplier: ingestion.NewSchemaApplier(),
		Resolver:      ingestion.NewResolver(r),
		RecordAdder:   ingestion.NewRecordAdder(r),
	}
}

func (fms *FormattedMessageSaver) HandleRecord(ctx context.Context, r *pb.DataRecord) (error) {
	return nil
}

func (fms *FormattedMessageSaver) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	log := Logger(ctx).Sugar()

	ds, err := fms.Resolver.ResolveDeviceAndSchemas(ctx, fm.SchemaId)
	if err != nil {
		return nil, err
	}

	pm, err := fms.SchemaApplier.ApplySchemas(ds, fm)
	if err != nil {
		return nil, err
	}

	change, err := fms.RecordAdder.AddRecord(ctx, ds, pm)
	if err != nil {
		return nil, err
	}

	fms.Changes[change.ID] = change

	log.Infow("Record", "device_id", fm.SchemaId.Device, "stream_id", fm.SchemaId.Stream, "modules", fm.Modules, "location", fm.Location, "source_id", change.SourceID)

	return change, nil
}

func (fms *FormattedMessageSaver) EmitChanges(ctx context.Context, sourceChanges ingestion.SourceChangesPublisher) {
	sources := make(map[int64][]*ingestion.RecordChange)
	for _, change := range fms.Changes {
		if sources[change.SourceID] == nil {
			sources[change.SourceID] = make([]*ingestion.RecordChange, 0)
		}
		sources[change.SourceID] = append(sources[change.SourceID], change)
	}
	for id, _ := range sources {
		sourceChanges.SourceChanged(ctx, ingestion.NewSourceChange(id))
	}
}
