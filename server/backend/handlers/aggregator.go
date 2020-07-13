package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/data"
)

var (
	AggregateNames = []string{"24h", "12h", "6h", "1h", "30m", "10m", "1m"}
)

var (
	AggregateTableNames = map[string]string{
		"24h": "fieldkit.aggregated_24h",
		"12h": "fieldkit.aggregated_12h",
		"6h":  "fieldkit.aggregated_6h",
		"1h":  "fieldkit.aggregated_1h",
		"30m": "fieldkit.aggregated_30m",
		"10m": "fieldkit.aggregated_10m",
		"1m":  "fieldkit.aggregated_1m",
	}
	AggregateTimeGroupThresholds = map[string]int{
		"24h": 86400 * 2,
		"12h": 86400,
		"6h":  86400,
		"1h":  3600 * 4,
		"30m": 3600 * 2,
		"10m": 30 * 60,
		"1m":  2 * 60,
	}
	AggregateIntervals = map[string]int32{
		"24h": 86400,
		"12h": 86400 / 2,
		"6h":  86400 / 4,
		"1h":  3600,
		"30m": 30 * 60,
		"10m": 10 * 60,
		"1m":  60,
	}
)

type aggregation struct {
	opened   time.Time
	interval time.Duration
	table    string
	values   map[string][]float64
	location []float64
}

type Aggregated struct {
	Time     time.Time
	Location []float64
	Values   map[string]float64
}

func (a *aggregation) getTime(original time.Time) time.Time {
	return original.UTC().Truncate(a.interval)
}

func (a *aggregation) canAdd(t time.Time) bool {
	return a.opened.IsZero() || a.opened == t
}

func (a *aggregation) add(t time.Time, key string, location []float64, value float64) error {
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
	agg := make(map[string]float64)

	for key, values := range a.values {
		if len(values) > 0 {
			mean, err := stats.Mean(values)
			if err != nil {
				return nil, fmt.Errorf("error taking mean: %v", err)
			}

			agg[key] = mean

			a.values[key] = a.values[key][:0]
		}
	}

	aggregated := &Aggregated{
		Time:   a.opened,
		Values: agg,
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

func NewAggregator(db *sqlxcache.DB, stationID int32, batchSize int) *Aggregator {
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
		},
		aggregations: []*aggregation{
			&aggregation{
				interval: time.Minute * 1,
				table:    "fieldkit.aggregated_1m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Minute * 10,
				table:    "fieldkit.aggregated_10m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Minute * 30,
				table:    "fieldkit.aggregated_30m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 1,
				table:    "fieldkit.aggregated_1h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 6,
				table:    "fieldkit.aggregated_6h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 12,
				table:    "fieldkit.aggregated_12h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 24,
				table:    "fieldkit.aggregated_24h",
				values:   make(map[string][]float64),
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
		log.Infow("batch", "aggregate", a.table, "records", len(batch))
		return nil
	})
}

func (v *Aggregator) upsertSingle(ctx context.Context, a *aggregation, d *Aggregated) error {
	var location *data.Location

	if d.Location != nil {
		location = data.NewLocation(d.Location)
	}

	if v.sensors == nil {
		if err := v.refreshSensors(ctx); err != nil {
			return err
		}
	}

	keeping := make([]int64, 0)

	for key, value := range d.Values {
		if v.sensors[key] == nil {
			newSensor := &data.Sensor{
				Key: key,
			}
			if err := v.db.NamedGetContext(ctx, newSensor, `
				INSERT INTO fieldkit.aggregated_sensor (key) VALUES (:key) RETURNING id
				`, newSensor); err != nil {
				return fmt.Errorf("error adding aggregated sensor: %v", err)
			}
			v.sensors[key] = newSensor
		}

		row := &data.AggregatedReading{
			StationID: v.stationID,
			SensorID:  v.sensors[key].ID,
			Time:      data.NumericWireTime(d.Time),
			Location:  location,
			Value:     value,
		}

		if err := v.db.NamedGetContext(ctx, row, fmt.Sprintf(`
			INSERT INTO %s (station_id, sensor_id, time, location, value)
			VALUES (:station_id, :sensor_id, :time, ST_SetSRID(ST_GeomFromText(:location), 4326), :value)
			ON CONFLICT (time, station_id, sensor_id) DO UPDATE SET value = EXCLUDED.value, location = EXCLUDED.location
			RETURNING id
			`, a.table), row); err != nil {
			return fmt.Errorf("error upserting sensor reading: %v", err)
		}

		keeping = append(keeping, row.ID)
	}

	// Delete any other rows for this station in this time.
	if len(keeping) > 0 {
		query, args, err := sqlx.In(fmt.Sprintf(`
			DELETE FROM %s WHERE station_id = ? AND time = ? AND id NOT IN (?)
		`, a.table), v.stationID, d.Time, keeping)
		if err != nil {
			return err
		}
		if _, err := v.db.ExecContext(ctx, v.db.Rebind(query), args...); err != nil {
			return fmt.Errorf("error deleting old sensor readings: %v", err)
		}
	} else {
		if _, err := v.db.ExecContext(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE station_id = $1 AND time = $2
		`, a.table), v.stationID, d.Time); err != nil {
			return fmt.Errorf("error deleting old sensor readings: %v", err)
		}
	}

	return nil
}

func (v *Aggregator) AddSample(ctx context.Context, sampled time.Time, location []float64, sensorKey string, value float64) error {
	for _, child := range v.aggregations {
		time := child.getTime(sampled)

		if !child.canAdd(time) {
			return fmt.Errorf("wow, NextTime call failed or missing")
		}

		if err := child.add(time, sensorKey, location, value); err != nil {
			return fmt.Errorf("error adding: %v", err)
		}
	}

	v.samples += 1

	return nil
}

func (v *Aggregator) AddMap(ctx context.Context, sampled time.Time, location []float64, data map[string]float64) error {
	for _, child := range v.aggregations {
		time := child.getTime(sampled)

		if !child.canAdd(time) {
			return fmt.Errorf("wow, NextTime call failed or missing")
		}

		for key, value := range data {
			if err := child.add(time, key, location, value); err != nil {
				return fmt.Errorf("error adding: %v", err)
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
		return fmt.Errorf("error closing aggregation: %v", err)
	} else {
		if v.batchSize > 1 {
			if len(v.batches[index]) == v.batchSize || tail {
				if err := v.upsertBatch(ctx, aggregation, v.batches[index]); err != nil {
					return err
				}
				v.batches[index] = v.batches[index][:0]
			}
			v.batches[index] = append(v.batches[index], values)
		} else {
			if err := v.upsertSingle(ctx, aggregation, values); err != nil {
				return fmt.Errorf("error upserting aggregated: %v", err)
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
	sensors := []*data.Sensor{}
	if err := v.db.SelectContext(ctx, &sensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return err
	}

	v.sensors = make(map[string]*data.Sensor)

	for _, sensor := range sensors {
		v.sensors[sensor.Key] = sensor
	}

	return nil
}
