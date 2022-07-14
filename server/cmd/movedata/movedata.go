package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/webhook"
)

type Options struct {
	PostgresURL string `split_words:"true" required:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable"`

	BinaryRecords        bool
	JsonRecords          bool
	SchemaID             int
	Verbose              bool
	FailOnMissingSensors bool

	InfluxDbURL    string `split_words:"true" required:"true" default:"http://127.0.0.1:8086"`
	InfluxDbToken  string `split_words:"true" required:"true"`
	InfluxDbOrg    string `split_words:"true" required:"true" default:"fk"`
	InfluxDbBucket string `split_words:"true" required:"true" default:"sensors"`
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
		stations:      make(map[string]*data.Station),
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

type MovedReading struct {
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

type MoveDataHandler interface {
	MoveReadings(ctx context.Context, readings []*MovedReading) error
	Close(ctx context.Context) error
}

func processBinary(ctx context.Context, options *Options, db *sqlxcache.DB, handler backend.RecordHandler) error {
	log := logging.Logger(ctx).Sugar()

	allStationIDs := []int32{}
	if err := db.SelectContext(ctx, &allStationIDs, "SELECT id FROM fieldkit.station ORDER BY ingestion_at DESC"); err != nil {
		return err
	}

	for _, id := range allStationIDs {
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

	_ = log

	return nil
}

func processJsonSchema(ctx context.Context, options *Options, db *sqlxcache.DB, schemaID int32, resolver *Resolver, handler MoveDataHandler) error {
	source := webhook.NewDatabaseMessageSource(db, schemaID, 0, false)

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

		for _, row := range batch.Messages {
			rowLog := logging.Logger(ctx).Sugar().With("schema_id", row.SchemaID).With("message_id", row.ID)

			allParsed, err := row.Parse(ctx, jqCache, batch.Schemas)
			if err != nil {
				rowLog.Infow("wh:skipping", "reason", err)
			} else {
				for _, parsed := range allParsed {
					if parsed.ReceivedAt != nil {
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
								sensorKey := parsedSensor.FullSensorKey

								sensorID, ok := resolver.sensors[sensorKey]
								if !ok {
									if options.FailOnMissingSensors {
										return fmt.Errorf("parsed-sensor for unknown sensor: %v", sensorKey)
									} else {
										rowLog.Infow("wh:missing", "key", sensorKey, "sensors", resolver.sensors)
									}
								} else {
									stationID, err := resolver.LookupStationID(ctx, parsed.DeviceID)
									if err != nil {
										return err
									}

									tags := make(map[string]string)
									tags["schema_id"] = fmt.Sprintf("%v", row.SchemaID)

									reading := &MovedReading{
										Time:      *parsed.ReceivedAt,
										DeviceID:  parsed.DeviceID,
										ModuleID:  parsed.DeviceID, // HACK This is what we do in model_adapter.go
										StationID: stationID,
										SensorID:  sensorID,
										SensorKey: sensorKey,
										Value:     parsedSensor.Value,
										Longitude: longitude,
										Latitude:  latitude,
										Tags:      tags,
									}

									if err := handler.MoveReadings(ctx, []*MovedReading{reading}); err != nil {
										return err
									}

									if options.Verbose {
										rowLog.Infow("wh:parsed", "received_at", parsed.ReceivedAt, "device_name", parsed.DeviceName, "data", parsed.Data)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func processJson(ctx context.Context, options *Options, db *sqlxcache.DB, resolver *Resolver, handler MoveDataHandler) error {
	ids := make([]int32, 0)
	if options.SchemaID == -1 {
		msr := webhook.NewMessageSchemaRepository(db)
		if schemas, err := msr.QueryAllSchemas(ctx); err != nil {
			return err
		} else {
			for _, schema := range schemas {
				ids = append(ids, schema.ID)
			}
		}
	} else {
		ids = append(ids, int32(options.SchemaID))
	}

	for _, schemaID := range ids {
		if err := processJsonSchema(ctx, options, db, schemaID, resolver, handler); err != nil {
			return err
		}
	}
	return nil
}

func (options *Options) createDestinationHandler(ctx context.Context) (MoveDataHandler, error) {
	influx := NewInflux(options.InfluxDbURL, options.InfluxDbToken, options.InfluxDbOrg, options.InfluxDbBucket)
	if err := influx.Open(ctx); err != nil {
		return nil, err
	}

	handler := NewMoveDataIntoInfluxHandler(influx)

	return handler, nil
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	resolver := NewResolver(db)
	if err := resolver.Open(ctx); err != nil {
		return err
	}

	destination, err := options.createDestinationHandler(ctx)
	if err != nil {
		return err
	}

	defer destination.Close(ctx)

	handler := NewMoveBinaryDataHandler(resolver, db, destination)

	if options.BinaryRecords {
		if err := processBinary(ctx, options, db, handler); err != nil {
			return err
		}
	}

	if options.JsonRecords {
		if err := processJson(ctx, options, db, resolver, destination); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.BoolVar(&options.BinaryRecords, "binary", false, "process binary records")
	flag.BoolVar(&options.JsonRecords, "json", false, "process json records")
	flag.IntVar(&options.SchemaID, "schema-id", -1, "schema id to process, -1 (default) for all")
	flag.BoolVar(&options.Verbose, "verbose", false, "increase verbosity")

	flag.Parse()

	logging.Configure(false, "move-data")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
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
