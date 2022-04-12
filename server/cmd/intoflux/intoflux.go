package main

import (
	"context"
	"flag"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/fieldkit/cloud/server/backend"
)

type Options struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
}

type InfluxDbHandler struct {
	cli   influxdb2.Client
	write api.WriteAPI
}

func NewInfluxDbHandler(address, token, org, bucket string) (h *InfluxDbHandler) {
	cli := influxdb2.NewClient(address, token)

	// Use blocking write client for writes to desired bucket
	write := cli.WriteAPI(org, bucket)

	return &InfluxDbHandler{
		cli:   cli,
		write: write,
	}
}

func (h *InfluxDbHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	log := logging.Logger(ctx).Sugar()
	log.Infow("meta")
	return nil
}

func (h *InfluxDbHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	log := logging.Logger(ctx).Sugar()
	log.Infow("data")
	return nil
}

func (h *InfluxDbHandler) OnDone(ctx context.Context) error {
	return nil
}

func (h *InfluxDbHandler) Close() {
	h.cli.Close()
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	stationIds := []int32{110}

	for _, id := range stationIds {
		walkParams := &backend.WalkParameters{
			Start:      time.Time{},
			End:        time.Now(),
			StationIDs: []int32{id},
		}

		rw := backend.NewRecordWalker(db)
		handler := &InfluxDbHandler{}
		if err := rw.WalkStation(ctx, handler, backend.WalkerProgressNoop, walkParams); err != nil {
			return err
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
