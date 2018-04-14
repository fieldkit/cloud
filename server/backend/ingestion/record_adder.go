package ingestion

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

type RecordChange struct {
	ID        int64
	Location  data.Location
	Timestamp *time.Time
}

type SourceChange struct {
	QueuedAt time.Time
	SourceID int64
	Records  []RecordChange
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

func (ra *RecordAdder) AddRecord(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage) (*data.Record, error) {
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

	id, err := ra.repository.AddRecord(ctx, record)
	if err != nil {
		return nil, err
	}

	record.ID = id

	if pm.LocationUpdated {
		if err := ra.repository.UpdateLocation(ctx, ds.Stream, pm.Location); err != nil {
			return nil, err
		}
	}

	return record, nil
}
