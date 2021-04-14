package repositories

import (
	"context"
	"time"
)

type Filter interface {
	Name() string
	Apply(ctx context.Context, record *ResolvedRecord, filters *MatchedFilters)
}

type Filtering struct {
	filters []Filter
}

func NewFiltering() (f *Filtering) {
	return &Filtering{
		filters: []Filter{
			&timeFilter{
				epoch: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			&sensorRangeFilter{},
			&multipleFilteredReadingsFilter{
				threshold: 3,
			},
			&emptyFilter{},
		},
	}
}

func (f *Filtering) Apply(ctx context.Context, record *ResolvedRecord) *FilteredRecord {
	filters := &MatchedFilters{
		Record:   make([]string, 0),
		Readings: make(map[SensorKey][]string),
	}

	for _, filter := range f.filters {
		filter.Apply(ctx, record, filters)
	}

	return &FilteredRecord{
		Record:  record,
		Filters: filters,
	}
}

type timeFilter struct {
	epoch time.Time
}

func (f *timeFilter) Name() string {
	return "time"
}

func (f *timeFilter) Apply(ctx context.Context, record *ResolvedRecord, filters *MatchedFilters) {
	if record.Time < f.epoch.Unix() {
		filters.AddRecord(f.Name())
	}
}

type sensorRangeFilter struct {
}

func (f *sensorRangeFilter) Name() string {
	return "range"
}

func (f *sensorRangeFilter) Apply(ctx context.Context, record *ResolvedRecord, filters *MatchedFilters) {
	for key, reading := range record.Readings {
		for _, r := range reading.Sensor.Ranges {
			if r.Minimum == r.Maximum {
				if reading.Value == r.Minimum {
					filters.AddReading(key, f.Name())
				}
			} else {
				if reading.Value < r.Minimum || reading.Value > r.Maximum {
					filters.AddReading(key, f.Name())
				}
			}
		}
		_ = key
	}
}

type emptyFilter struct {
}

func (f *emptyFilter) Name() string {
	return "empty"
}

func (f *emptyFilter) Apply(ctx context.Context, record *ResolvedRecord, filters *MatchedFilters) {
	if len(record.Readings) == 0 {
		filters.AddRecord(f.Name())
	}
}

type multipleFilteredReadingsFilter struct {
	threshold int
}

func (f *multipleFilteredReadingsFilter) Name() string {
	return "multiple"
}

func (f *multipleFilteredReadingsFilter) Apply(ctx context.Context, record *ResolvedRecord, filters *MatchedFilters) {
	if filters.NumberOfReadingsFiltered() >= f.threshold {
		filters.AddRecord(f.Name())
	}
}
