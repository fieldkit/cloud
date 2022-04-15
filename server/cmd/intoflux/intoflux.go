package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
)

type Options struct {
	PostgresURL    string `split_words:"true" required:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable"`
	InfluxDbURL    string `split_words:"true" required:"true" default:"http://127.0.0.1:8086"`
	InfluxDbToken  string `split_words:"true" required:"true"`
	InfluxDbOrg    string `split_words:"true" required:"true" default:"fk"`
	InfluxDbBucket string `split_words:"true" required:"true" default:"sensors"`
}

type stationInfo struct {
	provision *data.Provision
	meta      *data.MetaRecord
	station   *data.Station
}

type ImportErrors struct {
	MissingMeta map[string]int64
	ValueNaN    map[string]int64
}

type InfluxDbHandler struct {
	url          string
	token        string
	org          string
	bucket       string
	cli          influxdb2.Client
	write        api.WriteAPI
	metaFactory  *repositories.MetaFactory
	stations     *repositories.StationRepository
	querySensors *repositories.SensorRepository
	byProvision  map[int64]*stationInfo
	skipping     map[int64]bool
	sensors      map[string]int64
	errors       *ImportErrors
}

func NewInfluxDbHandler(url, token, org, bucket string, db *sqlxcache.DB) (h *InfluxDbHandler) {
	cli := influxdb2.NewClient(url, token)

	// Use blocking write client for writes to desired bucket
	write := cli.WriteAPI(org, bucket)

	return &InfluxDbHandler{
		url:          url,
		token:        token,
		org:          org,
		bucket:       bucket,
		cli:          cli,
		write:        write,
		metaFactory:  repositories.NewMetaFactory(db),
		stations:     repositories.NewStationRepository(db),
		querySensors: repositories.NewSensorRepository(db),
		byProvision:  make(map[int64]*stationInfo),
		skipping:     make(map[int64]bool),
		errors: &ImportErrors{
			ValueNaN:    make(map[string]int64),
			MissingMeta: make(map[string]int64),
		},
	}
}

func (h *InfluxDbHandler) Open(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("opening")

	if health, err := h.cli.Health(ctx); err != nil {
		return err
	} else {
		log.Infow("health", "heatlh", health)
	}

	if sensors, err := h.querySensors.QueryAllSensors(ctx); err != nil {
		return err
	} else {
		h.sensors = make(map[string]int64)
		for _, sensor := range sensors {
			h.sensors[sensor.Key] = sensor.ID
		}
	}

	return nil
}

func (h *InfluxDbHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	if v, ok := h.skipping[p.ID]; ok && v {
		return nil
	}

	log := logging.Logger(ctx).Sugar()

	if _, ok := h.byProvision[p.ID]; !ok {
		station, err := h.stations.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil {
			return err
		}

		h.byProvision[p.ID] = &stationInfo{
			meta:      meta,
			provision: p,
			station:   station,
		}
	}

	_, err := h.metaFactory.Add(ctx, meta, true)
	if err != nil {
		if _, ok := err.(*repositories.MissingSensorMetaError); ok {
			log.Infow("missing-meta", "meta_record_id", meta.ID)
			h.skipping[p.ID] = true
			return nil
		}
		return err
	}

	_ = log

	return nil
}

func (h *InfluxDbHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	if v, ok := h.skipping[p.ID]; ok && v {
		return nil
	}

	log := logging.Logger(ctx).Sugar()

	stationInfo := h.byProvision[p.ID]
	if stationInfo == nil {
		panic("ASSERT: missing station in data handler")
	}

	filtered, err := h.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("resolving: %v", err)
	}
	if filtered == nil {
		return nil
	}

	for key, rv := range filtered.Record.Readings {
		if math.IsNaN(rv.Value) {
			h.errors.ValueNaN[key.SensorKey] += 1
			continue
		}

		sensorID, ok := h.sensors[key.SensorKey]
		if !ok {
			h.errors.MissingMeta[key.SensorKey] += 1
			continue
		}

		// TODO Should/can we reuse maps for this?
		tags := make(map[string]string)
		fields := make(map[string]interface{})

		tags["station_id"] = fmt.Sprintf("%v", stationInfo.station.ID)
		tags["provision_id"] = fmt.Sprintf("%v", p.ID)
		tags["device_id"] = hex.EncodeToString(p.DeviceID)
		tags["sensor_id"] = fmt.Sprintf("%v", sensorID)
		tags["sensor_key"] = key.SensorKey

		fields["value"] = rv.Value

		dp := influxdb2.NewPoint("reading", tags, fields, time.Unix(filtered.Record.Time, 0))

		h.write.WritePoint(dp)
	}

	_ = log

	return nil
}

func (h *InfluxDbHandler) OnDone(ctx context.Context) error {
	return nil
}

func (h *InfluxDbHandler) Close() error {
	h.write.Flush()

	h.cli.Close()

	return nil
}

func (h *InfluxDbHandler) DoesStationHaveData(ctx context.Context, stationID int32) (bool, error) {
	queryAPI := h.cli.QueryAPI(h.org)

	query := fmt.Sprintf(`
		from(bucket:"%s")
		|> range(start: -10y, stop: now())
		|> filter(fn: (r) => r._measurement == "reading")
		|> filter(fn: (r) => r._field == "value")
		|> filter(fn: (r) => r.station_id == "%v")
	`, h.bucket, stationID)
	rows, err := queryAPI.Query(ctx, query)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		return true, nil
	}

	if rows.Err() != nil {
		return false, rows.Err()
	}

	return false, nil
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	allStationIDs := []int32{}
	if err := db.SelectContext(ctx, &allStationIDs, "SELECT id FROM fieldkit.station ORDER BY ingestion_at DESC"); err != nil {
		return err
	}

	handler := NewInfluxDbHandler(options.InfluxDbURL, options.InfluxDbToken, options.InfluxDbOrg, options.InfluxDbBucket, db)

	if err := handler.Open(ctx); err != nil {
		return err
	}

	stationIds := allStationIDs

	for _, id := range stationIds {
		if yes, err := handler.DoesStationHaveData(ctx, id); err != nil {
			return err
		} else if yes {
			log.Infow("skipping:has-data", "station_id", id)
			continue
		}

		walkParams := &backend.WalkParameters{
			Start:      time.Time{},
			End:        time.Now(),
			StationIDs: []int32{id},
		}

		rw := backend.NewRecordWalker(db)
		if err := rw.WalkStation(ctx, handler, backend.WalkerProgressNoop, walkParams); err != nil {
			return err
		}
	}

	if err := handler.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.Parse()

	logging.Configure(false, "influx")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
