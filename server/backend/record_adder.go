package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/data"
)

type RecordAdder struct {
	backend *Backend
}

func NewRecordAdder(backend *Backend) *RecordAdder {
	return &RecordAdder{
		backend: backend,
	}
}

func (da *RecordAdder) EmitSourceChanged(id int64) {
	da.backend.SourceChanges <- SourceChange{
		SourceID: id,
	}
}

func (da *RecordAdder) AddRecord(ctx context.Context, im *ingestion.IngestedMessage) error {
	ids := im.Schema.Ids
	d := data.Record{
		SchemaID:  int32(ids.SchemaID),
		SourceID:  int32(ids.DeviceID),
		TeamID:    nil,
		UserID:    nil,
		Timestamp: *im.Time,
		Location:  data.NewLocation(im.Location.Coordinates),
		Fixed:     im.Fixed,
		Visible:   true,
	}
	d.SetData(im.Fields)
	err := da.backend.AddRecord(ctx, &d)
	if err != nil {
		return err
	}

	return nil
}
