package ingestion

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

type RecordChange struct {
	ID        int64
	SourceID  int64
	Timestamp time.Time
	Location  *data.Location
}

type SourceChange struct {
	QueuedAt time.Time
	SourceID int64
}

func NewSourceChange(sourceId int64) SourceChange {
	return SourceChange{
		SourceID: sourceId,
		QueuedAt: time.Now(),
	}
}

type SourceChangesPublisher interface {
	SourceChanged(ctx context.Context, sourceChange SourceChange)
}

type SourceChangesBroadcaster struct {
	Publishers []SourceChangesPublisher
}

func NewSourceChangesBroadcaster(publishers []SourceChangesPublisher) *SourceChangesBroadcaster {
	return &SourceChangesBroadcaster{
		Publishers: publishers,
	}
}

func (b *SourceChangesBroadcaster) SourceChanged(ctx context.Context, sourceChange SourceChange) {
	for _, p := range b.Publishers {
		p.SourceChanged(ctx, sourceChange)
	}
}

type RecordAdder struct {
	repository *Repository
}

func NewRecordAdder(repository *Repository) *RecordAdder {
	return &RecordAdder{
		repository: repository,
	}
}

func (ra *RecordAdder) AddRecord(ctx context.Context, ds *DeviceStream, pm *ProcessedMessage) (*RecordChange, error) {
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

	if pm.LocationUpdated {
		if err := ra.repository.UpdateLocation(ctx, ds.Stream, pm.Location); err != nil {
			return nil, err
		}
	}

	change := &RecordChange{
		ID:        id,
		SourceID:  ids.DeviceID,
		Location:  record.Location,
		Timestamp: record.Timestamp,
	}

	return change, nil
}
