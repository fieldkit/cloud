package main

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
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
	return original.Truncate(a.interval)
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

	for key, _ := range a.values {
		mean, err := stats.Mean(a.values[key])
		if err != nil {
			return nil, err
		}

		agg[key] = mean

		a.values[key] = a.values[key][:0]
	}

	aggregated := &Aggregated{
		Time:   a.opened,
		Values: agg,
	}

	a.opened = time.Time{}

	return aggregated, nil
}

type AggregatingVisitor struct {
	db           *sqlxcache.DB
	metaFactory  *repositories.MetaFactory
	aggregations []*aggregation
	sensors      map[string]*Sensor
	stationID    int32
}

func NewAggregatingVisitor(db *sqlxcache.DB, stationID int32) *AggregatingVisitor {
	minutely := &aggregation{
		interval: time.Minute * 1,
		table:    "fieldkit.aggregated_minutely",
		values:   make(map[string][]float64),
	}
	hourly := &aggregation{
		interval: time.Hour * 1,
		table:    "fieldkit.aggregated_hourly",
		values:   make(map[string][]float64),
	}
	daily := &aggregation{
		interval: time.Hour * 24,
		table:    "fieldkit.aggregated_daily",
		values:   make(map[string][]float64),
	}
	return &AggregatingVisitor{
		db:          db,
		stationID:   stationID,
		metaFactory: repositories.NewMetaFactory(),
		aggregations: []*aggregation{
			minutely,
			hourly,
			daily,
		},
	}
}

func (v *AggregatingVisitor) refreshSensors(ctx context.Context) error {
	sensors := []*Sensor{}
	if err := v.db.SelectContext(ctx, &sensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return err
	}

	v.sensors = make(map[string]*Sensor)

	for _, sensor := range sensors {
		v.sensors[sensor.Key] = sensor
	}

	return nil
}

func (v *AggregatingVisitor) upsertAggregated(ctx context.Context, a *aggregation, d *Aggregated) error {
	var location *data.Location

	if d.Location != nil {
		location = data.NewLocation(d.Location)
	}

	if v.sensors == nil {
		if err := v.refreshSensors(ctx); err != nil {
			return err
		}
	}

	for key, value := range d.Values {
		if v.sensors[key] == nil {
			newSensor := &Sensor{
				Key: key,
			}
			if err := v.db.NamedGetContext(ctx, newSensor, `INSERT INTO fieldkit.aggregated_sensor (key) VALUES (:key)`, newSensor); err != nil {
				return err
			}
			v.sensors[key] = newSensor
		}

		row := &AggregatedReading{
			StationID: v.stationID,
			SensorID:  v.sensors[key].ID,
			Time:      d.Time,
			Location:  location,
			Value:     value,
		}

		if err := v.db.NamedGetContext(ctx, row, fmt.Sprintf(`
			INSERT INTO %s (station_id, sensor_id, time, location, value)
			VALUES (:station_id, :sensor_id, :time, ST_SetSRID(ST_GeomFromText(:location), 4326), :value)
			ON CONFLICT (time, station_id, sensor_id) DO UPDATE SET value = EXCLUDED.value, location = EXCLUDED.location
			RETURNING id
			`, a.table), row); err != nil {
			return err
		}

		_ = row
	}

	return nil
}

func (v *AggregatingVisitor) VisitMeta(ctx context.Context, meta *data.MetaRecord) error {
	_, err := v.metaFactory.Add(ctx, meta, true)
	if err != nil {
		return err
	}

	return nil
}

func (v *AggregatingVisitor) VisitData(ctx context.Context, meta *data.MetaRecord, data *data.DataRecord) error {
	filtered, err := v.metaFactory.Resolve(ctx, data, false, true)
	if err != nil {
		return err
	}

	for _, aggregation := range v.aggregations {
		time := aggregation.getTime(data.Time)

		if !aggregation.canAdd(time) {
			if values, err := aggregation.close(); err != nil {
				return err
			} else {
				if err := v.upsertAggregated(ctx, aggregation, values); err != nil {
					return err
				}
				if false {
					fmt.Printf("%v\n", values)
				}
			}
		}

		for key, value := range filtered.Record.Readings {
			if err := aggregation.add(time, key, filtered.Record.Location, value.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *AggregatingVisitor) VisitEnd(ctx context.Context) error {
	for _, aggregation := range v.aggregations {
		if values, err := aggregation.close(); err != nil {
			return err
		} else {
			if err := v.upsertAggregated(ctx, aggregation, values); err != nil {
				return err
			}
			if false {
				fmt.Printf("%v\n", values)
			}
		}
	}
	return nil
}

type Sensor struct {
	ID  int64  `db:"id"`
	Key string `db:"key"`
}

type AggregatedReading struct {
	ID        int64          `db:"id"`
	StationID int32          `db:"station_id"`
	SensorID  int64          `db:"sensor_id"`
	Time      time.Time      `db:"time"`
	Location  *data.Location `db:"location"`
	Value     float64        `db:"value"`
}
