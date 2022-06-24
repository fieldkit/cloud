package querying

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/fieldkit/cloud/server/backend"
)

type TimeScaleDBConfig struct {
	Url string
}

type TimeScaleDBBackend struct {
	config   *TimeScaleDBConfig
	cli      influxdb2.Client
	queryAPI api.QueryAPI
}

func NewTimeScaleDBBackend(config *TimeScaleDBConfig) (*TimeScaleDBBackend, error) {
	return &TimeScaleDBBackend{
		config: config,
	}, nil
}

func (tsdb *TimeScaleDBBackend) QueryRanges(ctx context.Context, qp *backend.QueryParams) (*DataRange, error) {
	log := Logger(ctx).Sugar()

	_ = log

	return nil, nil
}

func (tsdb *TimeScaleDBBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	if len(qp.Stations) != 1 {
		return nil, fmt.Errorf("multiple stations unsupported")
	}

	log := Logger(ctx).Sugar()

	_ = log

	_, err := tsdb.QueryRanges(ctx, qp)
	if err != nil {
		return nil, err
	}

	dataRows := make([]*backend.DataRow, 0)

	queriedData := &QueriedData{
		Summaries: make(map[string]*backend.AggregateSummary),
		Aggregate: AggregateInfo{
			Name:     "",
			Interval: 0,
			Complete: qp.Complete,
			Start:    qp.Start,
			End:      qp.End,
		},
		Data:  dataRows,
		Outer: make([]*backend.DataRow, 0),
	}

	return queriedData, nil
}

func (tddb *TimeScaleDBBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	return &SensorTailData{
		Data: make([]*backend.DataRow, 0),
	}, nil
}
