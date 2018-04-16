package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx/types"

	"github.com/Conservify/gonaturalist"
	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type CachedObservation struct {
	ID        int64          `db:"id,omitempty"`
	UpdatedAt time.Time      `db:"updated_at"`
	Timestamp time.Time      `db:"timestamp"`
	Location  *data.Location `db:"location"`
	Data      types.JSONText `db:"data"`
}

func NewCachedObservation(o *gonaturalist.SimpleObservation) (co *CachedObservation, err error) {
	co = &CachedObservation{
		ID:        o.Id,
		UpdatedAt: time.Now(),
		Timestamp: o.TimeObservedAtUtc,
		Location:  data.NewLocation([]float64{o.Longitude, o.Latitude}),
	}

	jsonData, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	co.Data = jsonData

	return
}

func (co *CachedObservation) Valid() bool {
	return !co.Timestamp.IsZero() && !co.Location.IsZero()
}

type INaturalistCache struct {
	Database         *sqlxcache.DB
	NaturalistClient *gonaturalist.Client
}

func NewINaturalistCache(config *INaturalistConfig, url string) (in *INaturalistCache, err error) {
	var authenticator = gonaturalist.NewAuthenticatorAtCustomRoot(config.ApplicationId, config.Secret, config.RedirectUrl, config.RootUrl)

	c := authenticator.NewClientWithAccessToken(iNaturalistConfig.AccessToken)

	db, err := backend.OpenDatabase(url)
	if err != nil {
		return nil, err
	}

	in = &INaturalistCache{
		Database:         db,
		NaturalistClient: c,
	}

	return
}

func (in *INaturalistCache) AddOrUpdateObservation(ctx context.Context, o *gonaturalist.SimpleObservation) (bool, error) {
	co, err := NewCachedObservation(o)
	if err != nil {
		return false, err
	}

	if !co.Valid() {
		return false, nil
	}

	_, err = in.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.inaturalist_observations (id, updated_at, timestamp, location, data)
		VALUES (:id, :updated_at, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326), :data)
		ON CONFLICT (id)
		DO UPDATE SET updated_at = excluded.updated_at, timestamp = excluded.timestamp, location = excluded.location, data = excluded.data
		`, co)

	return true, err
}

func (in *INaturalistCache) refreshUntilEmptyPage(ctx context.Context, options *gonaturalist.GetObservationsOpt) error {
	perPage := 100
	page := 0

	for {
		options.Page = &page
		options.PerPage = &perPage

		log.Printf("Refresh(%v)", spew.Sdump("%v", options))

		observations, err := in.NaturalistClient.GetObservations(options)
		if err != nil {
			return err
		}

		if len(observations.Observations) == 0 {
			break
		}

		for _, o := range observations.Observations {
			if _, err := in.AddOrUpdateObservation(ctx, o); err != nil {
				return err
			}
		}

		page += 1
	}

	return nil
}

func (in *INaturalistCache) RefreshRecentlyUpdated(ctx context.Context, updatedSince time.Time) error {
	orderBy := "created_at"
	hasGeo := true

	options := &gonaturalist.GetObservationsOpt{
		HasGeo:       &hasGeo,
		OrderBy:      &orderBy,
		UpdatedSince: &updatedSince,
	}

	return in.refreshUntilEmptyPage(ctx, options)
}

func (in *INaturalistCache) RefreshObservedOn(ctx context.Context, on time.Time) error {
	hasGeo := true

	options := &gonaturalist.GetObservationsOpt{
		HasGeo: &hasGeo,
		On:     &on,
	}

	return in.refreshUntilEmptyPage(ctx, options)
}
