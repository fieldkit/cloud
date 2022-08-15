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
	Data       []*backend.DataRow `json:"data"`
	BucketSize int                `json:"bucketSize"`
}

type DataBackend interface {
	QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error)
	QueryTail(ctx context.Context, stationIDs []int32) (*SensorTailData, error)
	QueryRecentlyAggregated(ctx context.Context, stationIDs []int32, windows []time.Duration) (map[time.Duration][]*backend.DataRow, error)
}
