package querying

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type StationTailInfo struct {
	BucketSize int `json:"bucketSize"`
}

type SensorTailData struct {
	Data     []*backend.DataRow         `json:"data"`
	Stations map[int32]*StationTailInfo `json:"stations"`
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

type StationLastTime struct {
	Last *data.NumericWireTime `json:"last"`
}

type RecentlyAggregated struct {
	Windows  map[time.Duration][]*backend.DataRow `json:"windows"`
	Stations map[int32]*StationLastTime           `json:"stations"`
}

type DataBackend interface {
	QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error)
	QueryTail(ctx context.Context, stationIDs []int32) (*SensorTailData, error)
	QueryRecentlyAggregated(ctx context.Context, stationIDs []int32, windows []time.Duration) (*RecentlyAggregated, error)
}
