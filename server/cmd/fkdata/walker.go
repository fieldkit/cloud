package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
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

func (rw *RecordWalker) WalkStation(ctx context.Context, stationID int32, handler backend.RecordHandler) error {
	pageSize := 1000
	done := false

	rw.started = time.Now()

	for page := 0; !done; page += 1 {
		_, err := rw.db.ExecContext(ctx, `
			DECLARE records_cursor CURSOR FOR
			SELECT r.id, r.provision_id, r.time, r.number, r.meta_record_id, ST_AsBinary(r.location) AS location, r.raw, r.pb
			FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
			WHERE (p.device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1))
			ORDER BY r.time OFFSET $2 LIMIT $3
		`, stationID, page*pageSize, pageSize)
		if err != nil {
			return err
		}

		if err := rw.walkQuery(ctx, handler); err != nil {
			if err == sql.ErrNoRows {
				done = true
				break
			}
			return err
		}

		if _, err := rw.db.ExecContext(ctx, "CLOSE records_cursor"); err != nil {
			return err
		}
	}

	if err := handler.OnDone(ctx); err != nil {
		return err
	}

	return nil
}

func (rw *RecordWalker) walkQuery(ctx context.Context, handler backend.RecordHandler) error {
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
			return err
		} else {
			if meta, err := rw.loadMeta(ctx, provision, data.MetaRecordID, handler); err != nil {
				return err
			} else {
				if err := handler.OnData(ctx, provision, nil, data, meta); err != nil {
					return err
				}

				rw.dataRecords += 1

				if rw.dataRecords%1000 == 0 {
					elapsed := time.Now().Sub(rw.started)
					log.Printf("%v records (%v rps)", rw.dataRecords, float64(rw.dataRecords)/elapsed.Seconds())
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

	return rw.provisions[0], nil
}

func (rw *RecordWalker) loadMeta(ctx context.Context, provision *data.Provision, id int64, handler backend.RecordHandler) (meta *data.MetaRecord, err error) {
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

	rw.metas[id] = records[0]
	rw.metaRecords += 1

	if err := handler.OnMeta(ctx, provision, nil, records[0]); err != nil {
		return nil, err
	}

	return records[0], nil
}
