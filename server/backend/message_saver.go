package backend

import (
	"context"
	"log"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type FormattedMessageSaver struct {
	SchemaApplier *ingestion.SchemaApplier
	Repository    *ingestion.Repository
	Resolver      *ingestion.Resolver
	RecordAdder   *ingestion.RecordAdder
	SourceIDs     map[int64]bool
}

func NewFormattedMessageSaver(b *Backend) *FormattedMessageSaver {
	r := ingestion.NewRepository(b.db)

	return &FormattedMessageSaver{
		SourceIDs:     make(map[int64]bool),
		Repository:    r,
		SchemaApplier: ingestion.NewSchemaApplier(),
		Resolver:      ingestion.NewResolver(r),
		RecordAdder:   ingestion.NewRecordAdder(r),
	}
}

func (br *FormattedMessageSaver) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) error {
	ds, err := br.Resolver.ResolveDeviceAndSchemas(ctx, fm.SchemaId)
	if err != nil {
		return err
	}

	pm, err := br.SchemaApplier.ApplySchemas(ds, fm)
	if err != nil {
		return err
	}

	err = br.RecordAdder.AddRecord(ctx, ds, pm)
	if err != nil {
		return err
	}

	br.SourceIDs[pm.Schema.Ids.DeviceID] = true

	log.Printf("(%s)(%s)[Success] %v, %d values (location = %t), %v", fm.MessageId, fm.SchemaId, fm.Modules, len(fm.MapValues), pm.LocationUpdated, fm.Location)

	return nil
}

func (br *FormattedMessageSaver) EmitChanges(sourceChanges ingestion.SourceChangesPublisher) {
	for id, _ := range br.SourceIDs {
		sourceChanges.SourceChanged(ingestion.NewSourceChange(id))
	}
}
