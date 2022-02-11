package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/lib/pq"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx/types"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
)

type Options struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	mmr := repositories.NewModuleMetaRepository(db)

	all, err := mmr.FindAllModulesMeta(ctx)
	if err != nil {
		return err
	}

	for _, mm := range all {
		fmt.Printf("%v\n", mm)

		kinds := make([]int32, len(mm.Header.AllKinds))
		for i, _ := range mm.Header.AllKinds {
			kinds[i] = int32(mm.Header.AllKinds[i])
		}

		pmm := &repositories.PersistedModuleMeta{
			Key:          mm.Key,
			Manufacturer: mm.Header.Manufacturer,
			Kinds:        pq.Int32Array(kinds),
			Version:      pq.Int32Array([]int32{int32(mm.Header.Version)}),
			Internal:     mm.Internal,
		}

		if err := db.NamedGetContext(ctx, pmm, `INSERT INTO fieldkit.module_meta (key, manufacturer, kinds, version, internal) VALUES (:key, :manufacturer, :kinds, :version, :internal) RETURNING id`, pmm); err != nil {
			return err
		}

		for _, sm := range mm.Sensors {
			stringsSerialized, err := json.Marshal(sm.Strings)
			if err != nil {
				return err
			}

			vizConfigs := sm.VizConfigs

			if vizConfigs == nil {
				vizConfigs = make([]repositories.VizConfig, 0)
			}

			vizSerialized, err := json.Marshal(vizConfigs)
			if err != nil {
				return err
			}

			rangesSerialized, err := json.Marshal(sm.Ranges)
			if err != nil {
				return err
			}

			psm := &repositories.PersistedSensorMeta{
				ModuleID:      pmm.ID,
				Ordering:      int32(sm.Order),
				SensorKey:     sm.Key,
				FullKey:       sm.FullKey,
				UnitOfMeasure: sm.UnitOfMeasure,
				Internal:      sm.Internal,
				Strings:       types.JSONText(stringsSerialized),
				Viz:           types.JSONText(vizSerialized),
				Ranges:        types.JSONText(rangesSerialized),
			}

			if err := db.NamedGetContext(ctx, psm, `
			INSERT INTO fieldkit.sensor_meta (module_id, ordering, sensor_key, full_key, internal, uom, strings, viz, ranges)
			VALUES (:module_id, :ordering, :sensor_key, :full_key, :internal, :uom, :strings, :viz, :ranges) RETURNING id
			`, psm); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.Parse()

	logging.Configure(false, "scratch")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
