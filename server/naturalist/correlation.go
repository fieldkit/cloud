package main

import (
	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
)

type INaturalistCorrelator struct {
	Database *sqlxcache.DB
}

func NewINaturalistCorrelator(url string) (nc *INaturalistCorrelator, err error) {
	db, err := backend.OpenDatabase(url)
	if err != nil {
		return nil, err
	}

	nc = &INaturalistCorrelator{
		Database: db,
	}

	return
}

func (nc *INaturalistCorrelator) Correlate(co *CachedObservation) (err error) {
	return nil
}
