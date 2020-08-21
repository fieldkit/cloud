package backend

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type RecordWalker struct {
	metas       map[int64]*data.MetaRecord
	provisions  map[int64]*data.Provision
	db          *sqlxcache.DB
	started     time.Time
	metaRecords int64
	dataRecords int64
}

func NewRecordWalker(db *sqlxcache.DB) (rw *RecordWalker) {
	return &RecordWalker{
		provisions: make(map[int64]*data.Provision),
		metas:      make(map[int64]*data.MetaRecord),
		db:         db,
	}
}

type WalkParameters struct {
	Page       int
	PageSize   int
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

func (rw *RecordWalker) WalkStation(ctx context.Context, handler RecordHandler, params *WalkParameters) error {
	done := false

	rw.started = time.Now()

	for params.Page = 0; !done; params.Page += 1 {
		if err := rw.processBatchInTransaction(ctx, handler, params); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}
	}

	elapsed := time.Now().Sub(rw.started)
	log := Logger(ctx).Sugar()

	if err := handler.OnDone(ctx); err != nil {
		return err
	}

	log.Infow("done", "station_ids", params.StationIDs, "records", rw.dataRecords, "rps", float64(rw.dataRecords)/elapsed.Seconds())

	return nil
}

func (rw *RecordWalker) processBatchInTransaction(ctx context.Context, handler RecordHandler, params *WalkParameters) error {
	return rw.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		return rw.processBatch(txCtx, handler, params)
	})
}

func (rw *RecordWalker) processBatch(ctx context.Context, handler RecordHandler, params *WalkParameters) error {
	log := Logger(ctx).Sugar()

	log.Infow("querying")

	query, args, err := sqlx.In(`
		DECLARE records_cursor CURSOR FOR
		SELECT
			r.id, r.provision_id, r.time, r.number, r.meta_record_id, r.raw, r.pb,
			ST_AsBinary(r.location) AS location
		FROM fieldkit.data_record AS r
		JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (
			p.device_id IN (SELECT device_id FROM fieldkit.station WHERE id IN (?))
			AND r.time >= ?
			AND r.time < ?
		)
		ORDER BY r.time OFFSET ? LIMIT ?
		`, params.StationIDs, params.Start, params.End, params.Page*params.PageSize, params.PageSize)
	if err != nil {
		return err
	}

	_, err = rw.db.ExecContext(ctx, rw.db.Rebind(query), args...)
	if err != nil {
		return err
	}

	done := false

	if err := rw.walkQuery(ctx, handler, params); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		done = true
	}

	if _, err := rw.db.ExecContext(ctx, `CLOSE records_cursor`); err != nil {
		return err
	}

	if done {
		return sql.ErrNoRows
	}

	return nil
}

func (rw *RecordWalker) walkQuery(ctx context.Context, handler RecordHandler, params *WalkParameters) error {
	log := Logger(ctx).Sugar()

	tx := rw.db.Transaction(ctx)
	batchSize := 0

	data := &data.DataRecord{}

	for {
		if err := tx.QueryRowxContext(ctx, `FETCH NEXT FROM records_cursor`).StructScan(data); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}

		if provision, err := rw.loadProvision(ctx, data.ProvisionID); err != nil {
			return fmt.Errorf("error loading provision: %v", err)
		} else {
			if meta, err := rw.loadMeta(ctx, provision, data.MetaRecordID, handler); err != nil {
				// return fmt.Errorf("error loading meta: %v", err)
				log.Infow("warning", "error", err)
			} else {
				if err := handler.OnData(ctx, provision, nil, data, meta); err != nil {
					return fmt.Errorf("error handling row: %v", err)
				}

				rw.dataRecords += 1

				if rw.dataRecords%1000 == 0 {
					elapsed := time.Now().Sub(rw.started)
					log.Infow("progress", "station_ids", params.StationIDs, "records", rw.dataRecords, "rps", float64(rw.dataRecords)/elapsed.Seconds())
				}
			}
		}

		batchSize += 1
	}

	if batchSize == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (rw *RecordWalker) loadProvision(ctx context.Context, id int64) (*data.Provision, error) {
	if rw.provisions[id] != nil {
		return rw.provisions[id], nil
	}

	records := []*data.Provision{}
	if err := rw.db.SelectContext(ctx, &records, `SELECT * FROM fieldkit.provision WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error querying for provision: %v", err)
	}

	if len(records) != 1 {
		return nil, fmt.Errorf("error finding provision: %v", id)
	}

	rw.provisions[id] = records[0]

	return records[0], nil
}

func (rw *RecordWalker) loadMeta(ctx context.Context, provision *data.Provision, id int64, handler RecordHandler) (meta *data.MetaRecord, err error) {
	if rw.metas[id] != nil {
		return rw.metas[id], nil
	}

	records := []*data.MetaRecord{}
	if err := rw.db.SelectContext(ctx, &records, `SELECT * FROM fieldkit.meta_record WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error querying for meta: %v", err)
	}

	if len(records) != 1 {
		return nil, fmt.Errorf("error finding meta: %v", id)
	}

	if err := handler.OnMeta(ctx, provision, nil, records[0]); err != nil {
		return nil, err
	}

	rw.metas[id] = records[0]
	rw.metaRecords += 1

	return records[0], nil
}
