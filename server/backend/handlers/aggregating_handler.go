package handlers

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type AggregatingHandler struct {
	db          *sqlxcache.DB
	metaFactory *repositories.MetaFactory
	stations    map[int64]*Aggregator
	aggregator  *Aggregator
}

func NewAggregatingHandler(db *sqlxcache.DB) *AggregatingHandler {
	return &AggregatingHandler{
		db:          db,
		metaFactory: repositories.NewMetaFactory(),
		stations:    make(map[int64]*Aggregator),
	}
}

func (v *AggregatingHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	if _, ok := v.stations[p.ID]; !ok {
		sr := repositories.NewStationRepository(v.db)

		station, err := sr.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil || station == nil {
			// mark giving up
			return nil
		}

		v.stations[p.ID] = NewAggregator(v.db, station.ID, 100)
	}

	_, err := v.metaFactory.Add(ctx, meta, true)
	if err != nil {
		return err
	}

	return nil
}

func (v *AggregatingHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	aggregator := v.stations[p.ID]
	if aggregator == nil {
		return fmt.Errorf("no aggregator for provision: %d", p.ID)
	}

	filtered, err := v.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("error resolving: %v", err)
	}
	if filtered == nil {
		return nil
	}

	if err := aggregator.NextTime(ctx, db.Time); err != nil {
		return fmt.Errorf("error adding: %v", err)
	}

	for key, value := range filtered.Record.Readings {
		if !filtered.Filters.IsFiltered(key) {
			if err := aggregator.AddSample(ctx, db.Time, filtered.Record.Location, key, value.Value); err != nil {
				return fmt.Errorf("error adding: %v", err)
			}
		}
	}

	return nil
}

func (v *AggregatingHandler) OnDone(ctx context.Context) error {
	for _, aggregator := range v.stations {
		if err := aggregator.Close(ctx); err != nil {
			return err
		}
	}

	// TODO Delete out of range data.

	return nil
}
