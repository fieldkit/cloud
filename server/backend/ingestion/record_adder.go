package ingestion

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

type SourceChange struct {
	QueuedAt   time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	SourceID   int64
}

func NewSourceChange(sourceId int64) SourceChange {
	return SourceChange{
		SourceID: sourceId,
		QueuedAt: time.Now(),
	}
}

type SourceChangesPublisher interface {
	SourceChanged(sourceChange SourceChange)
}

type RecordAdder struct {
	repository *Repository
}

func NewRecordAdder(repository *Repository) *RecordAdder {
	return &RecordAdder{
		repository: repository,
	}
}

func (da *RecordAdder) AddRecord(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage) error {
	ids := pm.Schema.Ids
	record := &data.Record{
		SchemaID:  int32(ids.SchemaID),
		SourceID:  int32(ids.DeviceID),
		TeamID:    nil,
		UserID:    nil,
		Timestamp: *pm.Time,
		Location:  data.NewLocation(pm.Location.Coordinates),
		Fixed:     pm.Fixed,
		Visible:   true,
	}

	record.SetData(pm.Fields)

	if err := da.repository.AddRecord(ctx, record); err != nil {
		return err
	}

	if pm.LocationUpdated {
		if err := da.repository.UpdateLocation(ctx, ds.Stream, pm.Location); err != nil {
			return err
		}
	}

	return nil
}
