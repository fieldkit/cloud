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
			&SensorRangeFilter{},
			&MultipleFilteredReadingsFilter{
				Threshold: 3,
			},
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

type MultipleFilteredReadingsFilter struct {
	Threshold int
}

func (f *MultipleFilteredReadingsFilter) Name() string {
	return "multiple"
}

func (f *MultipleFilteredReadingsFilter) Apply(ctx context.Context, row *FullDataRow, filters *MatchedFilters) {
	if filters.NumberOfReadingsFiltered() >= f.Threshold {
		filters.AddRecord(f.Name())
	}
}
