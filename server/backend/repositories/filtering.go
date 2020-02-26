package repositories

import (
	"context"
	"time"
)

type Filter interface {
	Name() string
	Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters)
}

type Filtering struct {
	Filters []Filter
}

func NewFiltering() (f *Filtering) {
	return &Filtering{
		Filters: []Filter{
			&TimeFilter{
				Epoch: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			&MissingWeatherSensorFilter{},
			&SensorRangeFilter{},
			&EmptyFilter{},
		},
	}
}

func (f *Filtering) Apply(ctx context.Context, record *FullDataRow) *FilteredRecord {
	filters := &MatchedFilters{
		Record:   make([]string, 0),
		Readings: make(map[string][]string),
	}

	for _, filter := range f.Filters {
		filter.Apply(ctx, record, filters)
	}

	return &FilteredRecord{
		Record:  record,
		Filters: filters,
	}
}

type TimeFilter struct {
	Epoch time.Time
}

func (f *TimeFilter) Name() string {
	return "time"
}

func (f *TimeFilter) Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters) {
	if row.Time < f.Epoch.Unix() {
		filters.AddRecord(f.Name())
	}
}

type SensorRangeFilter struct {
}

func (f *SensorRangeFilter) Name() string {
	return "range"
}

func (f *SensorRangeFilter) Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters) {
	for key, reading := range row.Readings {
		for _, r := range reading.Meta.Ranges {
			if reading.Value < r.Minimum || reading.Value > r.Maximum {
				filters.AddReading(key, f.Name())
			}
		}
		_ = key
	}
}

type EmptyFilter struct {
}

func (f *EmptyFilter) Name() string {
	return "empty"
}

func (f *EmptyFilter) Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters) {
	if len(row.Readings) == 0 {
		filters.AddRecord(f.Name())
	}
}

type MissingWeatherSensorFilter struct {
}

func (f *MissingWeatherSensorFilter) Name() string {
	return "no-sensor"
}

func (f *MissingWeatherSensorFilter) Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters) {
	pressure, ok := row.Readings["pressure"]
	if !ok {
		return
	}

	temperature1, ok := row.Readings["temperature1"]
	if !ok {
		return
	}

	if pressure.Value == 0 && temperature1.Value == -45 {
		filters.AddRecord(f.Name())
	}
}
