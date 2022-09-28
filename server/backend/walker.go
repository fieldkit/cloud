package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/hashicorp/go-multierror"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

const (
	InitialBatchSize = 50000
	MaximumBatchSize = 200000
)

type WalkerProgressFunc func(ctx context.Context, progress float64) error

func WalkerProgressNoop(ctx context.Context, progress float64) error {
	return nil
}

type WalkStatistics struct {
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Records int64      `db:"records"`
}

type RecordWalker struct {
	metas       map[int64]*data.MetaRecord
	provisions  map[int64]*data.Provision
	db          *sqlxcache.DB
	statistics  *WalkStatistics
	started     time.Time
	metaRecords int64
	dataRecords int64
	wg          sync.WaitGroup
	queue       chan []*data.DataRecord
	maximumID   int64
}

func NewRecordWalker(db *sqlxcache.DB) (rw *RecordWalker) {
	return &RecordWalker{
		db:         db,
		provisions: make(map[int64]*data.Provision),
		metas:      make(map[int64]*data.MetaRecord),
		queue:      make(chan []*data.DataRecord, 10),
		maximumID:  -1,
	}
}

type WalkParameters struct {
	StationIDs []int32
	Start      time.Time
	End        time.Time
}

type WalkInfo struct {
	MetaRecords int64
	DataRecords int64
}

func (rw *RecordWalker) Info(ctx context.Context) (*WalkInfo, error) {
	return &WalkInfo{
		DataRecords: rw.dataRecords,
		MetaRecords: rw.metaRecords,
	}, nil
}

func (rw *RecordWalker) WalkStation(ctx context.Context, handler RecordHandler, progress WalkerProgressFunc, params *WalkParameters) error {
	log := Logger(ctx).Sugar()

	rw.started = time.Now()

	ws, err := rw.queryStatistics(ctx, params)
	if err != nil {
		return err
	}

	rw.statistics = ws

	if ws.Records == 0 || ws.Start == nil || ws.End == nil {
		log.Infow("empty")
	} else {
		log.Infow("statistics", "start", *ws.Start, "end", *ws.End, "records", ws.Records)
	}

	var errors *multierror.Error

	rw.wg.Add(2)

	querierFunc := func() {
		offset := int64(0)
		batchSize := int64(InitialBatchSize)

		for {
			if rows, err := rw.queryBatch(ctx, params, offset, batchSize); err != nil {
				if err == sql.ErrNoRows {
					break
				}
				errors = multierror.Append(errors, err)
				break
			} else {
				if len(rw.queue) == 0 && batchSize < MaximumBatchSize {
					batchSize *= 2
				}

				rw.queue <- rows

				offset += int64(len(rows))
			}
		}

		close(rw.queue)

		rw.wg.Done()
	}

	processorFunc := func() {
		for {
			rows, valid := <-rw.queue
			if !valid {
				break
			}

			if err := rw.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
				for _, row := range rows {
					if err := rw.handleRecord(txCtx, handler, progress, row); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				errors = multierror.Append(errors, err)
				break
			}
		}

		rw.wg.Done()
	}

	go querierFunc()

	go processorFunc()

	rw.wg.Wait()

	elapsed := time.Now().Sub(rw.started)

	if err := handler.OnDone(ctx); err != nil {
		return err
	}

	log.Infow("done", "station_ids", params.StationIDs, "records", rw.dataRecords, "rps", float64(rw.dataRecords)/elapsed.Seconds())

	return errors.ErrorOrNil()
}

func (rw *RecordWalker) processBatchInTransaction(ctx context.Context, handler RecordHandler, progress WalkerProgressFunc, params *WalkParameters, offset, batchSize int64) error {
	return rw.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		return rw.processBatch(txCtx, handler, progress, params, offset, batchSize)
	})
}

func (rw *RecordWalker) queryStatistics(ctx context.Context, params *WalkParameters) (*WalkStatistics, error) {
	log := Logger(ctx).Sugar()

	log.Infow("query-statistics", "station_ids", params.StationIDs, "start", params.Start, "end", params.End)

	provisionIDs, err := rw.queryProvisionIDs(ctx, params.StationIDs)
	if err != nil {
		return nil, err
	}

	if len(provisionIDs) == 0 {
		return &WalkStatistics{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			MIN(r.time) AS start,
			MAX(r.time) AS end,
			COUNT(r.id) AS records
		FROM fieldkit.data_record AS r
		WHERE r.provision_id IN (?)
	 	AND r.time > ?
		`, provisionIDs, params.Start)
	if err != nil {
		return nil, err
	}

	ws := &WalkStatistics{}
	err = rw.db.GetContext(ctx, ws, rw.db.Rebind(query), args...)
	if err != nil {
		log.Errorw("error", "sql", rw.db.Rebind(query))
		return nil, err
	}

	return ws, nil
}

func (rw *RecordWalker) queryProvisionIDs(ctx context.Context, stationIDs []int32) ([]int64, error) {
	query, args, err := sqlx.In(`
		SELECT p.id FROM fieldkit.provision AS p WHERE p.device_id IN (SELECT device_id FROM fieldkit.station WHERE id IN (?))
		`, stationIDs)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = rw.db.SelectContext(ctx, &ids, rw.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (rw *RecordWalker) queryBatch(ctx context.Context, params *WalkParameters, offset int64, batchSize int64) ([]*data.DataRecord, error) {
	log := Logger(ctx).Sugar()

	log.Infow("query-rows", "station_ids", params.StationIDs, "offset", offset, "batch_size", batchSize, "start", params.Start, "end", params.End)

	rows := make([]*data.DataRecord, 0)

	provisionIDs, err := rw.queryProvisionIDs(ctx, params.StationIDs)
	if err != nil {
		return nil, err
	}

	if len(provisionIDs) == 0 {
		return nil, sql.ErrNoRows
	}

	query, args, err := sqlx.In(`
		SELECT
			r.id, r.provision_id, r.time, r.number, r.meta_record_id, r.raw AS raw, r.pb,
			ST_AsBinary(r.location) AS location
		FROM fieldkit.data_record AS r
		WHERE 
			r.provision_id IN (?)
			AND id > ?
			AND r.time > ?
		ORDER BY r.id
		LIMIT ?
		`, provisionIDs, rw.maximumID, params.Start, batchSize)
	if err != nil {
		return nil, err
	}

	err = rw.db.SelectContext(ctx, &rows, rw.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}

	rw.maximumID = rows[len(rows)-1].ID

	return rows, nil
}

func (rw *RecordWalker) processBatch(ctx context.Context, handler RecordHandler, progress WalkerProgressFunc, params *WalkParameters, offset, batchSize int64) error {
	rows, err := rw.queryBatch(ctx, params, offset, batchSize)
	if err != nil {
		return err
	}

	for _, row := range rows {
		if err := rw.handleRecord(ctx, handler, progress, row); err != nil {
			return err
		}
	}

	return nil
}

func (rw *RecordWalker) handleRecord(ctx context.Context, handler RecordHandler, progress WalkerProgressFunc, record *data.DataRecord) error {
	if provision, err := rw.loadProvision(ctx, record.ProvisionID); err != nil {
		return fmt.Errorf("error loading provision: %w", err)
	} else {
		if meta, err := rw.loadMeta(ctx, provision, record.MetaRecordID, handler); err != nil {
			return fmt.Errorf("(load-meta) error handling row: %w", err)
		} else {
			if err := handler.OnData(ctx, provision, nil, nil, record, meta); err != nil {
				return fmt.Errorf("(on-data) error handling row: %w", err)
			}
		}
	}

	rw.dataRecords += 1

	if rw.dataRecords%1000 == 0 {
		elapsed := time.Now().Sub(rw.started)
		percentage := float64(rw.dataRecords) / float64(rw.statistics.Records) * 100.0
		rps := float64(rw.dataRecords) / elapsed.Seconds()
		log := Logger(ctx).Sugar()
		log.Infow("progress", "records", rw.dataRecords, "rps", rps, "progress", percentage)
		if err := progress(ctx, percentage); err != nil {
			return err
		}
	}

	return nil
}

func (rw *RecordWalker) loadProvision(ctx context.Context, id int64) (*data.Provision, error) {
	if rw.provisions[id] != nil {
		return rw.provisions[id], nil
	}

	records := []*data.Provision{}
	if err := rw.db.SelectContext(ctx, &records, `SELECT * FROM fieldkit.provision WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error querying for provision: %w", err)
	}

	if len(records) != 1 {
		return nil, fmt.Errorf("error finding provision: %v", id)
	}

	rw.provisions[id] = records[0]

	return records[0], nil
}

func (rw *RecordWalker) loadMeta(ctx context.Context, provision *data.Provision, id int64, handler RecordHandler) (meta *data.MetaRecord, err error) {
	if rw.metas[id] != nil {
		// We used to skip this call only we can actually flip flop
		// between them in rare cases so better be safe and give the
		// handler a chance to make sure things are still expecting
		// this meta information.
		if err := handler.OnMeta(ctx, provision, nil, rw.metas[id]); err != nil {
			return nil, fmt.Errorf("(on-meta) error: %w", err)
		}
		return rw.metas[id], nil
	}

	records := []*data.MetaRecord{}
	if err := rw.db.SelectContext(ctx, &records, `SELECT * FROM fieldkit.meta_record WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error querying for meta: %w", err)
	}

	if len(records) != 1 {
		return nil, fmt.Errorf("error finding meta: %v", id)
	}

	if err := handler.OnMeta(ctx, provision, nil, records[0]); err != nil {
		return nil, fmt.Errorf("(on-meta) error: %w", err)
	}

	rw.metas[id] = records[0]
	rw.metaRecords += 1

	return records[0], nil
}
