package inaturalist

import (
	"context"
	"fmt"
	"time"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type CorrelationTier struct {
	Threshold float64
	Readings  []string
}

var Tiers = []*CorrelationTier{
	&CorrelationTier{
		Threshold: 1.0,
		Readings:  []string{},
	},
}

type INaturalistCorrelator struct {
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type Correlation struct {
	RecordsByTier map[*CorrelationTier][]*data.Record
}

func NewINaturalistCorrelator(db *sqlxcache.DB, be *backend.Backend) (nc *INaturalistCorrelator, err error) {
	nc = &INaturalistCorrelator{
		Database: db,
		Backend:  be,
	}

	return
}

func (nc *INaturalistCorrelator) Correlate(ctx context.Context, co *CachedObservation) (c *Correlation, err error) {
	log := Logger(ctx).Sugar()
	verbose := true

	offset := 1 * time.Minute
	options := backend.FindRecordsOpt{
		StartTime: co.Timestamp.Add(-1 * offset),
		EndTime:   co.Timestamp.Add(offset),
	}

	records, err := nc.Backend.FindRecords(ctx, &options)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		if verbose {
			log.Infow(fmt.Sprintf("No contemporaneous records (%v).", co.Timestamp), "observationId", co.ID, "siteId", co.SiteID)
		}
		return nil, nil
	}

	numberOfRecords := 0
	recordsByTier := make(map[*CorrelationTier][]*data.Record)

	for _, r := range records {
		distance := r.Location.Distance(co.Location)
		tiers := nc.findTiers(distance)
		if len(tiers) > 0 {
			for _, tier := range tiers {
				if recordsByTier[tier] == nil {
					recordsByTier[tier] = make([]*data.Record, 0)
				}
				recordsByTier[tier] = append(recordsByTier[tier], r)
				numberOfRecords += 1
			}
		}
	}

	if len(recordsByTier) == 0 {
		if verbose {
			log.Infow("No local records.", "observationId", co.ID, "siteId", co.SiteID)
		}
		return nil, nil
	}

	log.Infow("Observation correlated with records.", "observationId", co.ID, "siteId", co.SiteID, "numberOfRecords", numberOfRecords, "numberOfTiers", len(recordsByTier))

	c = &Correlation{
		RecordsByTier: recordsByTier,
	}

	return
}

func (nc *INaturalistCorrelator) findTiers(distance float64) []*CorrelationTier {
	matches := []*CorrelationTier{}
	for _, tier := range Tiers {
		if distance < tier.Threshold {
			matches = append(matches, tier)
		}
	}
	return matches
}
