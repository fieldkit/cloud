package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/logging"
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

func (br *FormattedMessageSaver) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) error {
	log := logging.Logger(ctx).Sugar()

	ds, err := br.Resolver.ResolveDeviceAndSchemas(ctx, fm.SchemaId)
	if err != nil {
		return err
	}

	pm, err := br.SchemaApplier.ApplySchemas(ds, fm)
	if err != nil {
		return err
	}

	change, err := br.RecordAdder.AddRecord(ctx, ds, pm)
	if err != nil {
		return err
	}

	br.Changes[change.ID] = change

	log.Infof("(%s)(%s)[Success] %v, %d values (location = %t), %v", fm.MessageId, fm.SchemaId, fm.Modules, len(fm.MapValues), pm.LocationUpdated, fm.Location)

	return nil
}

func (br *FormattedMessageSaver) EmitChanges(ctx context.Context, sourceChanges ingestion.SourceChangesPublisher) {
	sources := make(map[int64][]*ingestion.RecordChange)
	for _, change := range br.Changes {
		if sources[change.SourceID] == nil {
			sources[change.SourceID] = make([]*ingestion.RecordChange, 0)
		}
		sources[change.SourceID] = append(sources[change.SourceID], change)
	}
	for id, changes := range sources {
		sourceChanges.SourceChanged(ctx, ingestion.NewSourceChange(id, changes))
	}
}
