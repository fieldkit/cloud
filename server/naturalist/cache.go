package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx/types"

	"github.com/Conservify/gonaturalist"
	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

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

type CachedObservation struct {
	ID        int64          `db:"id,omitempty"`
	UpdatedAt time.Time      `db:"updated_at"`
	Timestamp time.Time      `db:"timestamp"`
	Location  *data.Location `db:"location"`
	Data      types.JSONText `db:"data"`
}

func (co *CachedObservation) Valid() bool {
	return !co.Timestamp.IsZero() && !co.Location.IsZero()
}

func (co *CachedObservation) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	co.Data = jsonData
	return nil
}

func (in *INaturalistCache) AddOrUpdateObservation(ctx context.Context, o *gonaturalist.SimpleObservation) (bool, error) {
	co := &CachedObservation{
		ID:        o.Id,
		UpdatedAt: time.Now(),
		Timestamp: o.TimeObservedAtUtc,
		Location:  data.NewLocation([]float64{o.Longitude, o.Latitude}),
	}

	if !co.Valid() {
		return false, nil
	}

	co.SetData(o)

	_, err := in.Database.NamedExecContext(ctx, `
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

		log.Printf("Refresh(%v)", page)

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

func (in *INaturalistCache) RefreshRecentlyAdded(ctx context.Context, updatedSince time.Time) error {
	orderBy := "created_at"
	hasGeo := true

	options := &gonaturalist.GetObservationsOpt{
		OrderBy:      &orderBy,
		HasGeo:       &hasGeo,
		UpdatedSince: &updatedSince,
	}

	return in.refreshUntilEmptyPage(ctx, options)
}

func (in *INaturalistCache) RefreshObservedOn(ctx context.Context, on time.Time) error {
	hasGeo := true

	options := &gonaturalist.GetObservationsOpt{
		On:     &on,
		HasGeo: &hasGeo,
	}

	return in.refreshUntilEmptyPage(ctx, options)
}
