package querying

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/backend"
)

type SensorTailData struct {
	Data []*backend.DataRow `json:"data"`
}

type AggregateInfo struct {
	Name     string    `json:"name"`
	Interval int32     `json:"interval"`
	Complete bool      `json:"complete"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type QueriedData struct {
	Summaries map[string]*backend.AggregateSummary `json:"summaries"`
	Aggregate AggregateInfo                        `json:"aggregate"`
	Data      []*backend.DataRow                   `json:"data"`
	Outer     []*backend.DataRow                   `json:"outer"`
}

type DataBackend interface {
	QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error)
	QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error)
}
