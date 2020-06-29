package main

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type aggregation struct {
	opened   time.Time
	interval time.Duration
	table    string
	values   map[string][]float64
}

type Aggregated struct {
	Time   time.Time
	Values map[string]float64
}

func (a *aggregation) getTime(original time.Time) time.Time {
	return original.Truncate(a.interval)
}

func (a *aggregation) canAdd(t time.Time) bool {
	return a.opened.IsZero() || a.opened == t
}

func (a *aggregation) add(t time.Time, key string, value float64) error {
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
	metaFactory  *repositories.MetaFactory
	aggregations []*aggregation
}

func NewAggregatingVisitor() *AggregatingVisitor {
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
		metaFactory: repositories.NewMetaFactory(),
		aggregations: []*aggregation{
			minutely,
			hourly,
			daily,
		},
	}
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
				if false {
					fmt.Printf("%v\n", values)
				}
			}
		}

		for key, value := range filtered.Record.Readings {
			if err := aggregation.add(time, key, value.Value); err != nil {
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
			if false {
				fmt.Printf("%v\n", values)
			}
		}
	}
	return nil
}

type Sensor struct {
	ID        int64  `db:"id"`
	ModuleKey string `db:"module_key"`
	Key       string `db:"key"`
}

type AggregatedReading struct {
	ID        int64          `db:"id"`
	StationID int32          `db:"station_id"`
	SensorID  int64          `db:"sensor_id"`
	Time      time.Time      `db:"time"`
	Location  *data.Location `db:"location"`
	Value     float64        `db:"value"`
}
