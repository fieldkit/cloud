package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

type QueryParams struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Sensors    []int64   `json:"sensors"`
	Stations   []int64   `json:"stations"`
	Resolution int32     `json:"resolution"`
	Aggregate  string    `json:"aggregate"`
	Tail       int32     `json:"tail"`
	Complete   bool      `json:"complete"`
}

type RawQueryParams struct {
	Start      *int64
	End        *int64
	Resolution *int32
	Stations   *string
	Sensors    *string
	Aggregate  *string
	Tail       *int32
	Complete   *bool
}

func (raw *RawQueryParams) BuildQueryParams() (qp *QueryParams, err error) {
	start := time.Time{}
	if raw.Start != nil {
		start = time.Unix(0, *raw.Start*int64(time.Millisecond))
	}

	end := time.Now()
	if raw.End != nil {
		end = time.Unix(0, *raw.End*int64(time.Millisecond))
	}

	resolution := int32(0)
	if raw.Resolution != nil {
		resolution = *raw.Resolution
	}

	stations := make([]int64, 0)
	if raw.Stations != nil {
		parts := strings.Split(*raw.Stations, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				stations = append(stations, int64(i))
			}
		}
	}

	sensors := make([]int64, 0)
	if raw.Sensors != nil {
		parts := strings.Split(*raw.Sensors, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				sensors = append(sensors, int64(i))
			}
		}
	}

	aggregate := handlers.AggregateNames[0]
	if raw.Aggregate != nil {
		found := false
		for _, name := range handlers.AggregateNames {
			if name == *raw.Aggregate {
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid aggregate: %v", *raw.Aggregate)
		}

		aggregate = *raw.Aggregate
	}

	tail := int32(0)
	if raw.Tail != nil {
		tail = *raw.Tail
	}

	complete := raw.Complete != nil && *raw.Complete

	qp = &QueryParams{
		Start:      start,
		End:        end,
		Resolution: resolution,
		Stations:   stations,
		Sensors:    sensors,
		Aggregate:  aggregate,
		Tail:       tail,
		Complete:   complete,
	}

	return
}
