package main

import (
	"context"
	"log"
	"time"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
)

type INaturalistCorrelator struct {
	Database *sqlxcache.DB
	Backend  *backend.Backend
	Verbose  bool
}

func NewINaturalistCorrelator(url string) (nc *INaturalistCorrelator, err error) {
	db, err := backend.OpenDatabase(url)
	if err != nil {
		return nil, err
	}

	be, err := backend.New(url)
	if err != nil {
		return nil, err
	}

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

	for _, r := range records {
		distance := r.Location.Distance(co.Location)
		log.Printf("Record distance=%fm", distance)
	}

	return nil
}

func (nc *INaturalistCorrelator) Handle(ctx context.Context, co *CachedObservation) error {
	return nc.Correlate(ctx, co)
}
