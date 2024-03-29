package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

var (
	AggregateNames = []string{"24h", "12h", "6h", "1h", "30m", "10m", "1m", "10s"}
)

var (
	AggregateTimeGroupThresholds = map[string]int{
		"24h": 86400 * 2,
		"12h": 86400,
		"6h":  86400,
		"1h":  3600 * 4,
		"30m": 3600 * 2,
		"10m": 30 * 60,
		"1m":  2 * 60,
		"10s": 30,
	}
	AggregateIntervals = map[string]int32{
		"24h": 86400,
		"12h": 86400 / 2,
		"6h":  86400 / 4,
		"1h":  3600,
		"30m": 30 * 60,
		"10m": 10 * 60,
		"1m":  60,
		"10s": 10,
	}
)

type AggregationFunction interface {
	Apply(values []float64) (float64, error)
}

type AverageFunction struct {
}

func (f *AverageFunction) Apply(values []float64) (float64, error) {
	return stats.Mean(values)
}

type MaximumFunction struct {
}

func (f *MaximumFunction) Apply(values []float64) (float64, error) {
	return stats.Max(values)
}

type AggregateSensorKey struct {
	SensorKey string
	ModuleID  int64
}

type aggregation struct {
	opened   time.Time
	interval time.Duration
	name     string
	table    string
	values   map[AggregateSensorKey][]float64
	location []float64
	config   AggregatorConfig
}

type Aggregated struct {
	Time          time.Time
	Location      []float64
	Values        map[AggregateSensorKey]float64
	NumberSamples int32
}

func MakeAggregateTableName(suffix string, name string) string {
	return fmt.Sprintf("%s%s_%s", "fieldkit.aggregated", suffix, name)
}

func (a *aggregation) getTime(original time.Time) time.Time {
	return original.UTC().Truncate(a.interval)
}

func (a *aggregation) canAdd(t time.Time) bool {
	return a.opened.IsZero() || a.opened == t
}

func (a *aggregation) add(t time.Time, key AggregateSensorKey, location []float64, value float64) error {
	if !a.canAdd(t) {
		return fmt.Errorf("unable to add to aggregation")
	}

	if a.opened != t && a.opened.IsZero() {
		a.opened = t
	}

	if a.values[key] == nil {
		a.values[key] = make([]float64, 0, 1)
	}

	a.values[key] = append(a.values[key], value)
	if len(location) == 3 && location[0] != 0 && location[1] != 0 {
		a.location = location
	}

	return nil
}

func (a *aggregation) close() (*Aggregated, error) {
	agg := make(map[AggregateSensorKey]float64)

	nsamples := 0

	for key, values := range a.values {
		if len(values) > 0 {
			aggregated_value, err := a.config.Apply(key, values)
			if err != nil {
				return nil, fmt.Errorf("error taking mean: %w", err)
			}

			agg[key] = aggregated_value

			if len(values) > nsamples {
				nsamples = len(values)
			}

			a.values[key] = a.values[key][:0]
		}
	}

	aggregated := &Aggregated{
		Time:          a.opened,
		Location:      a.location,
		Values:        agg,
		NumberSamples: int32(nsamples),
	}

	a.opened = time.Time{}

	return aggregated, nil
}

type Aggregator struct {
	db           *sqlxcache.DB
	stationID    int32
	aggregations []*aggregation
	sensors      map[string]*data.Sensor
	samples      int64
	batchSize    int
	batches      [][]*Aggregated
}

type AggregatorConfig interface {
	Apply(sensorKey AggregateSensorKey, values []float64) (float64, error)
}

type defaultAggregatorConfig struct{}

func NewDefaultAggregatorConfig() *defaultAggregatorConfig {
	return &defaultAggregatorConfig{}
}

func (c *defaultAggregatorConfig) Apply(sensorKey AggregateSensorKey, values []float64) (float64, error) {
	return stats.Mean(values)
}

func NewAggregator(db *sqlxcache.DB, tableSuffix string, stationID int32, batchSize int, config AggregatorConfig) *Aggregator {
	return &Aggregator{
		db:        db,
		stationID: stationID,
		batchSize: batchSize,
		batches: [][]*Aggregated{
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
			make([]*Aggregated, 0, batchSize),
		},
		aggregations: []*aggregation{
			&aggregation{
				interval: time.Second * 10,
				name:     "10s",
				table:    fmt.Sprintf("fieldkit.aggregated%s_10s", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Minute * 1,
				name:     "1m",
				table:    fmt.Sprintf("fieldkit.aggregated%s_1m", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Minute * 10,
				name:     "10m",
				table:    fmt.Sprintf("fieldkit.aggregated%s_10m", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Minute * 30,
				name:     "30m",
				table:    fmt.Sprintf("fieldkit.aggregated%s_30m", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Hour * 1,
				name:     "1h",
				table:    fmt.Sprintf("fieldkit.aggregated%s_1h", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Hour * 6,
				name:     "6h",
				table:    fmt.Sprintf("fieldkit.aggregated%s_6h", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Hour * 12,
				name:     "12h",
				table:    fmt.Sprintf("fieldkit.aggregated%s_12h", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
			&aggregation{
				interval: time.Hour * 24,
				name:     "24h",
				table:    fmt.Sprintf("fieldkit.aggregated%s_24h", tableSuffix),
				values:   make(map[AggregateSensorKey][]float64),
				config:   config,
			},
		},
	}
}

func (v *Aggregator) upsertBatch(ctx context.Context, a *aggregation, batch []*Aggregated) error {
	return v.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		for _, d := range batch {
			if err := v.upsertSingle(txCtx, a, d); err != nil {
				return err
			}
		}
		log := Logger(ctx).Sugar().With("station_id", v.stationID)
		log.Infow("batch", "aggregate", a.name, "records", len(batch))
		return nil
	})
}

func (v *Aggregator) upsertSingle(ctx context.Context, a *aggregation, d *Aggregated) error {
	var location *data.Location

	if d.Location != nil && len(d.Location) == 3 {
		location = data.NewLocation(d.Location)
	}
	if v.sensors == nil {
		if err := v.refreshSensors(ctx); err != nil {
			return err
		}
	}

	for key, value := range d.Values {
		if v.sensors[key.SensorKey] == nil {
			newSensor := &data.Sensor{
				Key: key.SensorKey,
			}
			if err := v.db.NamedGetContext(ctx, newSensor, `
				INSERT INTO fieldkit.aggregated_sensor (key) VALUES (:key) RETURNING id
				`, newSensor); err != nil {
				return fmt.Errorf("error adding aggregated sensor: %w", err)
			}
			v.sensors[key.SensorKey] = newSensor
		}

		row := &data.AggregatedReading{
			StationID:     v.stationID,
			SensorID:      v.sensors[key.SensorKey].ID,
			ModuleID:      key.ModuleID,
			Time:          data.NumericWireTime(d.Time),
			Location:      location,
			Value:         value,
			NumberSamples: d.NumberSamples,
		}

		if err := v.db.NamedGetContext(ctx, row, fmt.Sprintf(`
			INSERT INTO %s (station_id, sensor_id, module_id, time, location, value, nsamples)
			VALUES (:station_id, :sensor_id, :module_id, :time, ST_SetSRID(ST_GeomFromText(:location), 4326), :value, :nsamples)
			ON CONFLICT (time, station_id, sensor_id, module_id)
			DO UPDATE SET value = EXCLUDED.value, location = EXCLUDED.location, nsamples = EXCLUDED.nsamples
			RETURNING id
			`, a.table), row); err != nil {
			return fmt.Errorf("error upserting sensor reading: %w", err)
		}
	}

	return nil
}

func (v *Aggregator) ClearNumberSamples(ctx context.Context) error {
	log := Logger(ctx).Sugar().With("station_id", v.stationID)

	log.Infow("nsamples-clear")

	for _, child := range v.aggregations {
		if _, err := v.db.ExecContext(ctx, fmt.Sprintf(`UPDATE %s SET nsamples = 0 WHERE station_id = $1`, child.table), v.stationID); err != nil {
			return fmt.Errorf("error updating aggregate nsamples: %w", err)
		}
	}

	return nil
}

func (v *Aggregator) DeleteEmptyAggregates(ctx context.Context) error {
	log := Logger(ctx).Sugar().With("station_id", v.stationID)

	for _, child := range v.aggregations {
		total := int32(0)
		if err := v.db.GetContext(ctx, &total, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE nsamples = 0 AND station_id = $1`, child.table), v.stationID); err != nil {
			return fmt.Errorf("error deleting empty aggregates: %w", err)
		}

		log.Infow("nsamples-delete", "table", child.table, "records", total)

		if _, err := v.db.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE nsamples = 0 AND station_id = $1`, child.table), v.stationID); err != nil {
			return fmt.Errorf("error deleting empty aggregates: %w", err)
		}
	}

	return nil
}

func (v *Aggregator) AddSample(ctx context.Context, sampled time.Time, location []float64, sensorKey AggregateSensorKey, value float64) error {
	for _, child := range v.aggregations {
		time := child.getTime(sampled)

		if !child.canAdd(time) {
			return fmt.Errorf("wow, NextTime call failed or missing")
		}

		if err := child.add(time, sensorKey, location, value); err != nil {
			return fmt.Errorf("error adding: %w", err)
		}
	}

	v.samples += 1

	return nil
}

func (v *Aggregator) addMap(ctx context.Context, sampled time.Time, location []float64, data map[AggregateSensorKey]float64) error {
	for _, child := range v.aggregations {
		time := child.getTime(sampled)

		if !child.canAdd(time) {
			return fmt.Errorf("wow, NextTime call failed or missing")
		}

		for key, value := range data {
			if err := child.add(time, key, location, value); err != nil {
				return fmt.Errorf("error adding: %w", err)
			}
		}
	}

	v.samples += 1

	return nil
}

func (v *Aggregator) NextTime(ctx context.Context, sampled time.Time) error {
	for index, child := range v.aggregations {
		time := child.getTime(sampled)

		if !child.canAdd(time) {
			if err := v.closeChild(ctx, index, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Aggregator) closeChild(ctx context.Context, index int, tail bool) error {
	aggregation := v.aggregations[index]

	if values, err := aggregation.close(); err != nil {
		return fmt.Errorf("error closing aggregation: %w", err)
	} else {
		if v.batchSize > 1 {
			if len(v.batches[index]) == v.batchSize || tail {
				if tail {
					v.batches[index] = append(v.batches[index], values)
				}
				if err := v.upsertBatch(ctx, aggregation, v.batches[index]); err != nil {
					return err
				}
				v.batches[index] = v.batches[index][:0]
			}
			if !tail {
				v.batches[index] = append(v.batches[index], values)
			}
		} else {
			if err := v.upsertSingle(ctx, aggregation, values); err != nil {
				return fmt.Errorf("error upserting aggregated: %w", err)
			}
		}
	}
	return nil
}

func (v *Aggregator) Close(ctx context.Context) error {
	for index, _ := range v.aggregations {
		if err := v.closeChild(ctx, index, true); err != nil {
			return err
		}
	}
	return nil
}

func (v *Aggregator) refreshSensors(ctx context.Context) error {
	querySensors := repositories.NewSensorsRepository(v.db)

	queried, err := querySensors.QueryAllSensors(ctx)
	if err != nil {
		return err
	}

	v.sensors = queried

	return nil
}
