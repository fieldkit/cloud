package naturalist

import (
	"context"
	"log"
	"time"

	"github.com/Conservify/sqlxcache"
	"github.com/davecgh/go-spew/spew"

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
	Verbose  bool
}

func NewINaturalistCorrelator(db *sqlxcache.DB, be *backend.Backend) (nc *INaturalistCorrelator, err error) {
	nc = &INaturalistCorrelator{
		Database: db,
		Backend:  be,
	}

	return
}

func (nc *INaturalistCorrelator) Correlate(ctx context.Context, co *CachedObservation) (err error) {
	offset := 1 * time.Minute
	options := backend.FindRecordsOpt{
		StartTime: co.Timestamp.Add(-1 * offset),
		EndTime:   co.Timestamp.Add(offset),
	}

	records, err := nc.Backend.FindRecords(ctx, &options)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		if nc.Verbose {
			log.Printf("No contemporaneous records.")
		}
		return nil
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
		if nc.Verbose {
			log.Printf("No local records.")
		}
		return nil
	}

	if nc.Verbose {
		log.Printf("%s", spew.Sdump(recordsByTier))
	} else {
		log.Printf("Observation #%d has %d records across %d tier(s)", co.ID, numberOfRecords, len(recordsByTier))
	}

	return nil
}

func (nc *INaturalistCorrelator) Handle(ctx context.Context, co *CachedObservation) error {
	return nc.Correlate(ctx, co)
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
