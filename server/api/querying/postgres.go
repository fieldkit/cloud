package querying

import (
	"context"
	"fmt"
	"math"

	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type PostgresBackend struct {
	db *sqlxcache.DB
}

func NewPostgresBackend(db *sqlxcache.DB) *PostgresBackend {
	return &PostgresBackend{db: db}
}

func (pgb *PostgresBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	log := Logger(ctx).Sugar()

	dq := backend.NewDataQuerier(pgb.db)

	summaries, selectedAggregateName, err := dq.SelectAggregate(ctx, qp)
	if err != nil {
		return nil, err
	}

	aqp, err := backend.NewAggregateQueryParams(qp, selectedAggregateName, summaries[selectedAggregateName])
	if err != nil {
		return nil, err
	}

	rows := make([]*backend.DataRow, 0)
	if aqp.ExpectedRecords > 0 {
		queried, err := dq.QueryAggregate(ctx, aqp)
		if err != nil {
			return nil, err
		}

		defer queried.Close()

		for queried.Next() {
			row := &backend.DataRow{}
			if err = scanRow(queried, row); err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	} else {
		log.Infow("empty summary")
	}

	outerRows, err := dq.QueryOuterValues(ctx, aqp)
	if err != nil {
		return nil, err
	}

	if qp.Complete {
		hackedRows := make([]*backend.DataRow, 0, len(rows)+len(outerRows))
		if outerRows[0] != nil {
			outerRows[0].Time = data.NumericWireTime(aqp.Start)
			hackedRows = append(hackedRows, outerRows[0])
		}
		for i := 0; i < len(rows); i += 1 {
			hackedRows = append(hackedRows, rows[i])
		}
		if outerRows[1] != nil {
			outerRows[1].Time = data.NumericWireTime(aqp.End)
			hackedRows = append(hackedRows, outerRows[1])
		}

		rows = hackedRows
	}

	data := &QueriedData{
		Summaries: summaries,
		Aggregate: AggregateInfo{
			Name:     aqp.AggregateName,
			Interval: aqp.Interval,
			Complete: aqp.Complete,
			Start:    aqp.Start,
			End:      aqp.End,
		},
		Data:  rows,
		Outer: outerRows,
	}

	return data, nil
}

func (pgb *PostgresBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
		id,
		time,
		station_id,
		sensor_id,
		module_id,
		value
		FROM (
			SELECT
				ROW_NUMBER() OVER (PARTITION BY agg.sensor_id ORDER BY time DESC) AS r,
				agg.*
			FROM %s AS agg
			WHERE agg.station_id IN (?)
		) AS q
		WHERE
		q.r <= ?
		`, "fieldkit.aggregated_1m"), qp.Stations, qp.Tail)
	if err != nil {
		return nil, err
	}

	queried, err := pgb.db.QueryxContext(ctx, pgb.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*backend.DataRow, 0)

	for queried.Next() {
		row := &backend.DataRow{}
		if err = scanRow(queried, row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return &SensorTailData{
		Data: rows,
	}, nil
}

func scanRow(queried *sqlx.Rows, row *backend.DataRow) error {
	if err := queried.StructScan(row); err != nil {
		return fmt.Errorf("error scanning row: %v", err)
	}

	if row.Value != nil && math.IsNaN(*row.Value) {
		row.Value = nil
	}

	return nil
}
