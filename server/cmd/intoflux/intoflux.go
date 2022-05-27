package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/webhook"
)

type Options struct {
	PostgresURL          string `split_words:"true" required:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable"`
	InfluxDbURL          string `split_words:"true" required:"true" default:"http://127.0.0.1:8086"`
	InfluxDbToken        string `split_words:"true" required:"true"`
	InfluxDbOrg          string `split_words:"true" required:"true" default:"fk"`
	InfluxDbBucket       string `split_words:"true" required:"true" default:"sensors"`
	SchemaID             int
	Verbose              bool
	BinaryRecords        bool
	JsonRecords          bool
	FailOnMissingSensors bool
}

type Influx struct {
	url    string
	token  string
	org    string
	bucket string
	cli    influxdb2.Client
	write  api.WriteAPI
}

type InfluxReading struct {
	Time      time.Time
	StationID int32
	DeviceID  []byte
	ModuleID  []byte
	SensorID  int64
	SensorKey string
	Value     float64
	Longitude *float64
	Latitude  *float64
	Altitude  *float64
	Tags      map[string]string
}

func (r *InfluxReading) ToPoint() *write.Point {
	tags := make(map[string]string)
	tags["device_id"] = hex.EncodeToString(r.DeviceID)
	tags["module_id"] = hex.EncodeToString(r.ModuleID)
	tags["station_id"] = fmt.Sprintf("%v", r.StationID)
	tags["sensor_id"] = fmt.Sprintf("%v", r.SensorID)
	tags["sensor_key"] = r.SensorKey

	fields := make(map[string]interface{})
	fields["value"] = r.Value

	if r.Longitude != nil && r.Latitude != nil {
		fields["lon"] = *r.Longitude
		fields["lat"] = *r.Latitude
	}

	if r.Altitude != nil {
		fields["altitude"] = *r.Altitude
	}

	return influxdb2.NewPoint("reading", tags, fields, r.Time)
}

func NewInflux(url, token, org, bucket string, db *sqlxcache.DB) (h *Influx) {
	cli := influxdb2.NewClient(url, token)

	// Use blocking write client for writes to desired bucket
	write := cli.WriteAPI(org, bucket)

	return &Influx{
		url:    url,
		token:  token,
		org:    org,
		bucket: bucket,
		cli:    cli,
		write:  write,
	}
}

func (h *Influx) Open(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("opening")

	if health, err := h.cli.Health(ctx); err != nil {
		return err
	} else {
		log.Infow("health", "heatlh", health)
	}

	return nil
}

func (h *Influx) Close() error {
	h.write.Flush()

	h.cli.Close()

	return nil
}

func (h *Influx) DoesStationHaveData(ctx context.Context, stationID int32) (bool, error) {
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

type Resolver struct {
	queryStations *repositories.StationRepository
	querySensors  *repositories.SensorsRepository
	sensors       map[string]int64
	stations      map[string]*data.Station
}

func NewResolver(db *sqlxcache.DB) *Resolver {
	return &Resolver{
		queryStations: repositories.NewStationRepository(db),
		querySensors:  repositories.NewSensorsRepository(db),
	}
}

func (r *Resolver) Open(ctx context.Context) error {
	if sensors, err := r.querySensors.QueryAllSensors(ctx); err != nil {
		return err
	} else {
		r.sensors = make(map[string]int64)
		for _, sensor := range sensors {
			r.sensors[sensor.Key] = sensor.ID
		}
	}

	return nil
}

func (r *Resolver) LookupStationID(ctx context.Context, deviceID []byte) (int32, error) {
	cacheKey := hex.EncodeToString(deviceID)
	if station, ok := r.stations[cacheKey]; ok {
		return station.ID, nil
	}

	station, err := r.queryStations.QueryStationByDeviceID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	r.stations[cacheKey] = station

	return station.ID, nil
}

type stationInfo struct {
	provision *data.Provision
	meta      *data.MetaRecord
	station   *data.Station
}

type importErrors struct {
	ValueNaN      map[string]int64
	MissingMeta   map[string]int64
	MalformedMeta map[string]int64
}

type InfluxDbHandler struct {
	influx       *Influx
	resolve      *Resolver
	metaFactory  *repositories.MetaFactory
	stations     *repositories.StationRepository
	querySensors *repositories.SensorsRepository
	byProvision  map[int64]*stationInfo
	skipping     map[int64]bool
	errors       *importErrors
}

func NewInfluxDbHandler(influx *Influx, resolve *Resolver, db *sqlxcache.DB) (h *InfluxDbHandler) {
	return &InfluxDbHandler{
		influx:       influx,
		resolve:      resolve,
		metaFactory:  repositories.NewMetaFactory(db),
		stations:     repositories.NewStationRepository(db),
		querySensors: repositories.NewSensorsRepository(db),
		byProvision:  make(map[int64]*stationInfo),
		skipping:     make(map[int64]bool),
		errors: &importErrors{
			ValueNaN:      make(map[string]int64),
			MissingMeta:   make(map[string]int64),
			MalformedMeta: make(map[string]int64),
		},
	}
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
		if _, ok := err.(*repositories.MalformedMetaError); ok {
			log.Infow("malformed-meta", "meta_record_id", meta.ID)
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

		sensorID, ok := h.resolve.sensors[key.SensorKey]
		if !ok {
			h.errors.MissingMeta[key.SensorKey] += 1
			continue
		}

		moduleID, err := hex.DecodeString(rv.Module.ID)
		if err != nil {
			return err
		}

		var latitude *float64
		var longitude *float64
		var altitude *float64

		if len(filtered.Record.Location) >= 2 {
			longitude = &filtered.Record.Location[0]
			latitude = &filtered.Record.Location[1]
			if len(filtered.Record.Location) >= 3 {
				altitude = &filtered.Record.Location[2]
			}
		}

		// TODO Should/can we reuse maps for this?
		tags := make(map[string]string)
		tags["provision_id"] = fmt.Sprintf("%v", p.ID)

		reading := InfluxReading{
			Time:      time.Unix(filtered.Record.Time, 0),
			DeviceID:  p.DeviceID,
			ModuleID:  moduleID,
			StationID: stationInfo.station.ID,
			SensorID:  sensorID,
			SensorKey: key.SensorKey,
			Value:     rv.Value,
			Tags:      tags,
			Longitude: longitude,
			Latitude:  latitude,
			Altitude:  altitude,
		}

		dp := reading.ToPoint()

		h.influx.write.WritePoint(dp)
	}

	_ = log

	return nil
}

func (h *InfluxDbHandler) OnDone(ctx context.Context) error {
	return nil
}

func processBinary(ctx context.Context, options *Options, db *sqlxcache.DB, influx *Influx, resolver *Resolver) error {
	log := logging.Logger(ctx).Sugar()

	allStationIDs := []int32{}
	if err := db.SelectContext(ctx, &allStationIDs, "SELECT id FROM fieldkit.station ORDER BY ingestion_at DESC"); err != nil {
		return err
	}

	stationIds := allStationIDs

	handler := NewInfluxDbHandler(influx, resolver, db)

	for _, id := range stationIds {
		if yes, err := influx.DoesStationHaveData(ctx, id); err != nil {
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

	return nil
}

func toFloatArray(x interface{}) ([]float64, bool) {
	if arrayValue, ok := x.([]interface{}); ok {
		values := make([]float64, 0)
		for _, opaque := range arrayValue {
			if v, ok := toFloat(opaque); ok {
				values = append(values, v)
			}
		}
		return values, true
	}
	return nil, false
}

func toFloat(x interface{}) (float64, bool) {
	switch x := x.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	case *big.Int:
		f, err := strconv.ParseFloat(x.String(), 64)
		return f, err == nil
	default:
		return 0.0, false
	}
}

func processJson(ctx context.Context, options *Options, db *sqlxcache.DB, influx *Influx, resolver *Resolver) error {
	schemas := webhook.NewMessageSchemaRepository(db)

	source := webhook.NewDatabaseMessageSource(db, int32(options.SchemaID))

	jqCache := &webhook.JqCache{}

	batch := &webhook.MessageBatch{
		StartTime: time.Time{},
	}

	for {
		batchLog := logging.Logger(ctx).Sugar().With("batch_start_time", batch.StartTime)

		if err := source.NextBatch(ctx, batch); err != nil {
			if err == sql.ErrNoRows || err == io.EOF {
				batchLog.Infow("eof")
				break
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		_, err := schemas.QuerySchemas(ctx, batch)
		if err != nil {
			return fmt.Errorf("message schemas (%v)", err)
		}

		for _, row := range batch.Messages {
			rowLog := logging.Logger(ctx).Sugar().With("schema_id", row.SchemaID).With("message_id", row.ID)

			parsed, err := row.Parse(ctx, jqCache, batch.Schemas)
			if err != nil {
				rowLog.Infow("wh:skipping", "reason", err)
			} else if parsed != nil {
				var latitude *float64
				var longitude *float64

				for _, attribute := range parsed.Attributes {
					if attribute.Location && attribute.JSONValue != nil {
						if coordinates, ok := toFloatArray(attribute.JSONValue); ok {
							if len(coordinates) >= 2 {
								longitude = &coordinates[0]
								latitude = &coordinates[1]
								break
							}
						}
					}
				}

				for _, parsedSensor := range parsed.Data {
					key := parsedSensor.Key
					if key == "" {
						return fmt.Errorf("parsed-sensor has no sensor key")
					}

					if !parsedSensor.Transient {
						sensorID, ok := resolver.sensors[key]
						if !ok {
							if options.FailOnMissingSensors {
								return fmt.Errorf("parsed-sensor for unknown sensor: %v", key)
							}
						} else {
							stationID, err := resolver.LookupStationID(ctx, parsed.DeviceID)
							if err != nil {
								return err
							}

							tags := make(map[string]string)
							tags["schema_id"] = fmt.Sprintf("%v", row.SchemaID)

							reading := InfluxReading{
								Time:      parsed.ReceivedAt,
								DeviceID:  parsed.DeviceID,
								ModuleID:  parsed.DeviceID, // HACK This is what we do in model_adapter.go
								StationID: stationID,
								SensorID:  sensorID,
								SensorKey: key,
								Value:     parsedSensor.Value,
								Longitude: longitude,
								Latitude:  latitude,
								Tags:      tags,
							}

							dp := reading.ToPoint()

							influx.write.WritePoint(dp)

							if options.Verbose {
								rowLog.Infow("wh:parsed", "received_at", parsed.ReceivedAt, "device_name", parsed.DeviceName, "data", parsed.Data)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	influx := NewInflux(options.InfluxDbURL, options.InfluxDbToken, options.InfluxDbOrg, options.InfluxDbBucket, db)

	if err := influx.Open(ctx); err != nil {
		return err
	}

	resolver := NewResolver(db)

	if err := resolver.Open(ctx); err != nil {
		return err
	}

	if options.BinaryRecords {
		if err := processBinary(ctx, options, db, influx, resolver); err != nil {
			return err
		}
	}

	if options.JsonRecords {
		if err := processJson(ctx, options, db, influx, resolver); err != nil {
			return err
		}
	}

	if err := influx.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.BoolVar(&options.BinaryRecords, "binary", false, "process binary records")
	flag.BoolVar(&options.JsonRecords, "json", false, "process json records")
	flag.BoolVar(&options.Verbose, "verbose", false, "increase verbosity")
	flag.IntVar(&options.SchemaID, "schema-id", -1, "schema id to process, -1 (default) for all")

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
