package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/fieldkit/cloud/server/common/logging"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Influx struct {
	url    string
	token  string
	org    string
	bucket string
	cli    influxdb2.Client
	write  api.WriteAPI
}

func ToPoint(r *MovedReading) *write.Point {
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

func NewInflux(url, token, org, bucket string) (h *Influx) {
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

type MoveDataToInfluxHandler struct {
	influx *Influx
}

func NewMoveDataIntoInfluxHandler(influx *Influx) *MoveDataToInfluxHandler {
	return &MoveDataToInfluxHandler{
		influx: influx,
	}
}

func (h *MoveDataToInfluxHandler) MoveReadings(ctx context.Context, readings []*MovedReading) error {
	for _, r := range readings {
		dp := ToPoint(r)
		h.influx.write.WritePoint(dp)
	}
	return nil
}

func (h *MoveDataToInfluxHandler) Close(ctx context.Context) error {
	return nil
}
