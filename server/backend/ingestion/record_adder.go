package ingestion

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type SourceChange struct {
	Time     time.Time
	SourceID int64
}

func NewSourceChange(sourceId int64) SourceChange {
	return SourceChange{
		SourceID: sourceId,
		Time:     time.Now(),
	}
}

type SourceChangesPublisher interface {
	SourceChanged(sourceChange SourceChange)
}

type RecordAdder struct {
	db *sqlxcache.DB
}

func NewRecordAdder(db *sqlxcache.DB) *RecordAdder {
	return &RecordAdder{
		db: db,
	}
}

func (da *RecordAdder) AddRecord(ctx context.Context, pm *ProcessedMessage) error {
	ids := pm.Schema.Ids
	record := data.Record{
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

	_, err := da.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.record (schema_id, source_id, team_id, user_id, timestamp, location, data, fixed, visible)
		VALUES (:schema_id, :source_id, :team_id, :user_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326), :data, :fixed, :visible)
		`, record)
	if err != nil {
		return err
	}

	return nil
}
