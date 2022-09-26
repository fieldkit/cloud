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

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/fieldkit/cloud/server/webhook"
)

type Options struct {
	PostgresURL  string `split_words:"true"`
	TimeScaleURL string `split_words:"true"`

	AWSProfile     string   `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	StreamsBuckets []string `split_words:"true" default:""`
	AwsId          string   `split_words:"true" default:""`
	AwsSecret      string   `split_words:"true" default:""`

	BinaryRecords        bool
	StationID            int
	JsonRecords          bool
	SchemaID             int
	Ingestions           bool
	ID                   int
	Verbose              bool
	FailOnMissingSensors bool
	RefreshViews         bool

	InfluxDbURL    string `split_words:"true"`
	InfluxDbToken  string `split_words:"true"`
	InfluxDbOrg    string `split_words:"true" default:"fk"`
	InfluxDbBucket string `split_words:"true" default:"sensors"`
}

type Resolver struct {
	queryStations *repositories.StationRepository
	querySensors  *repositories.SensorsRepository
	sensors       map[string]int64
	stations      map[string]*data.Station
	modules       map[string]int64
}

func NewResolver(db *sqlxcache.DB) *Resolver {
	return &Resolver{
		queryStations: repositories.NewStationRepository(db),
		querySensors:  repositories.NewSensorsRepository(db),
		stations:      make(map[string]*data.Station),
		modules:       make(map[string]int64),
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

func (r *Resolver) LookupModuleID(ctx context.Context, hardwareID []byte) (int64, error) {
	cacheKey := hex.EncodeToString(hardwareID)
	if moduleID, ok := r.modules[cacheKey]; ok {
		return moduleID, nil
	}

	modules, err := r.queryStations.QueryStationModulesByHardwareID(ctx, hardwareID)
	if err != nil {
		return 0, err
	}

	if len(modules) == 0 {
		return 0, fmt.Errorf("module missing: %v", cacheKey)
	}

	if len(modules) != 1 {
		return 0, fmt.Errorf("module ambiguous: %v", cacheKey)
	}

	r.modules[cacheKey] = modules[0].ID

	return modules[0].ID, nil
}

type MovedReading struct {
	Time             time.Time
	StationID        int32
	DeviceID         []byte
	ModuleHardwareID []byte
	ModuleID         int64
	SensorID         int64
	SensorKey        string
	Value            float64
	Longitude        *float64
	Latitude         *float64
	Altitude         *float64
	Tags             map[string]string
}

type MoveDataHandler interface {
	MoveReadings(ctx context.Context, readings []*MovedReading) error
	Close(ctx context.Context) error
}

func processBinary(ctx context.Context, options *Options, db *sqlxcache.DB, handler backend.RecordHandler) error {
	log := logging.Logger(ctx).Sugar()

	allStationIDs := []int32{}
	if options.StationID > 0 {
		allStationIDs = append(allStationIDs, int32(options.StationID))
	} else {
		if err := db.SelectContext(ctx, &allStationIDs, "SELECT id FROM fieldkit.station ORDER BY ingestion_at DESC"); err != nil {
			return err
		}
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

									moduleHardwareID := parsed.DeviceID // HACK This is what we do in model_adapter.go

									moduleID, err := resolver.LookupModuleID(ctx, moduleHardwareID)
									if err != nil {
										return err
									}

									tags := make(map[string]string)
									tags["schema_id"] = fmt.Sprintf("%v", row.SchemaID)

									reading := &MovedReading{
										Time:             *parsed.ReceivedAt,
										DeviceID:         parsed.DeviceID,
										ModuleHardwareID: moduleHardwareID,
										StationID:        stationID,
										ModuleID:         moduleID,
										SensorID:         sensorID,
										SensorKey:        sensorKey,
										Value:            parsedSensor.Value,
										Longitude:        longitude,
										Latitude:         latitude,
										Tags:             tags,
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

func processStationIngestions(ctx context.Context, options *Options, db *sqlxcache.DB, dbpool *pgxpool.Pool, resolver *Resolver, handler MoveDataHandler, stationID int32) error {
	ir, err := repositories.NewIngestionRepository(db)
	if err != nil {
		return err
	}

	ingestions, err := ir.QueryByStationID(ctx, stationID)
	if err != nil {
		return err
	}

	for _, ingestion := range ingestions {
		if err := processIngestion(ctx, options, db, dbpool, resolver, handler, ingestion.ID); err != nil {
			return err
		}
	}

	return nil
}

func (config *Options) getAwsSessionOptions() session.Options {
	if config.AwsId == "" || config.AwsSecret == "" {
		return session.Options{
			Profile: config.AWSProfile,
			Config: aws.Config{
				Region:                        aws.String("us-east-1"),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}
	}

	return session.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewStaticCredentials(config.AwsId, config.AwsSecret, ""),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}
}

func processIngestion(ctx context.Context, options *Options, db *sqlxcache.DB, dbpool *pgxpool.Pool, resolver *Resolver, handler MoveDataHandler, ingestionID int64) error {
	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(ctx, publisher, nil)
	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{})

	awsSessionOptions := options.getAwsSessionOptions()

	awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		return err
	}

	reading := make([]files.FileArchive, 0)
	writing := make([]files.FileArchive, 0)

	fs := files.NewLocalFilesArchive()
	reading = append(reading, fs)
	writing = append(writing, fs)

	for _, bucketName := range options.StreamsBuckets {
		s3, err := files.NewS3FileArchive(awsSession, metrics, bucketName, files.NoPrefix)
		if err != nil {
			return err
		}
		reading = append(reading, s3)
	}

	archive := files.NewPrioritizedFilesArchive(reading, writing)

	ingestionReceived := backend.NewIngestionReceivedHandler(db, dbpool, archive, metrics, publisher, options.timeScaleConfig())

	ir, err := repositories.NewIngestionRepository(db)
	if err != nil {
		return err
	}

	if id, err := ir.Enqueue(ctx, ingestionID); err != nil {
		return err
	} else {
		if err := ingestionReceived.Start(ctx, &messages.IngestionReceived{
			QueuedID: id,
			UserID:   2,
			Verbose:  true,
			Refresh:  true,
		}, mc); err != nil {
			return err
		}
	}

	return nil
}

func refreshViews(ctx context.Context, options *Options) error {
	tsConfig := options.timeScaleConfig()

	if tsConfig == nil {
		return fmt.Errorf("refresh-views missing tsdb configuration")
	}

	return tsConfig.RefreshViews(ctx)
}

func (options *Options) timeScaleConfig() *storage.TimeScaleDBConfig {
	if options.TimeScaleURL == "" {
		return nil
	}

	return &storage.TimeScaleDBConfig{Url: options.TimeScaleURL}
}

func (options *Options) createDestinationHandler(ctx context.Context) (MoveDataHandler, error) {
	tsConfig := options.timeScaleConfig()
	if tsConfig != nil {
		handler := NewMoveDataToTimeScaleDBHandler(tsConfig)

		return handler, nil
	}

	if options.InfluxDbURL != "" {
		influx := NewInflux(options.InfluxDbURL, options.InfluxDbToken, options.InfluxDbOrg, options.InfluxDbBucket)
		if err := influx.Open(ctx); err != nil {
			return nil, err
		}

		handler := NewMoveDataIntoInfluxHandler(influx)

		return handler, nil
	}

	return nil, fmt.Errorf("invalid destination configuration")
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open(ctx, "postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	pgxcfg, err := pgxpool.ParseConfig(options.PostgresURL)
	if err != nil {
		return err
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, pgxcfg)
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

	if options.Ingestions {
		if options.ID > 0 {
			if err := processIngestion(ctx, options, db, dbpool, resolver, destination, int64(options.ID)); err != nil {
				return err
			}
		}
		if options.StationID > 0 {
			if err := processStationIngestions(ctx, options, db, dbpool, resolver, destination, int32(options.StationID)); err != nil {
				return err
			}
		}
	}

	if options.RefreshViews {
		if err := refreshViews(ctx, options); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.BoolVar(&options.BinaryRecords, "binary", false, "process binary records")
	flag.IntVar(&options.StationID, "station-id", -1, "station id to process, -1 (default) for all")
	flag.BoolVar(&options.JsonRecords, "json", false, "process json records")
	flag.IntVar(&options.SchemaID, "schema-id", -1, "schema id to process, -1 (default) for all")
	flag.BoolVar(&options.Verbose, "verbose", false, "increase verbosity")
	flag.BoolVar(&options.Ingestions, "ingestions", false, "process ingestions")
	flag.IntVar(&options.ID, "id", -1, "id to process")
	flag.BoolVar(&options.RefreshViews, "refresh-views", false, "refresh views")

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
