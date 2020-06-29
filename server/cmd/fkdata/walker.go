package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type RecordWalker struct {
	metas       map[int64]*data.MetaRecord
	db          *sqlxcache.DB
	started     time.Time
	metaRecords int64
	dataRecords int64
}

func NewRecordWalker(db *sqlxcache.DB) (rw *RecordWalker) {
	return &RecordWalker{
		metas: make(map[int64]*data.MetaRecord),
		db:    db,
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

func (rw *RecordWalker) WalkStation(ctx context.Context, stationID int32, visitor RecordVisitor) error {
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

		if err := rw.walkQuery(ctx, visitor); err != nil {
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

	if err := visitor.VisitEnd(ctx); err != nil {
		return err
	}

	return nil
}

func (rw *RecordWalker) walkQuery(ctx context.Context, visitor RecordVisitor) error {
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

		if rw.metas[data.MetaRecordID] == nil {
			if meta, err := rw.loadMeta(ctx, data.MetaRecordID); err != nil {
				return err
			} else {
				if err := visitor.VisitMeta(ctx, meta); err != nil {
					return err
				}
			}

			rw.metaRecords += 1
		}

		meta := rw.metas[data.MetaRecordID]
		if err := visitor.VisitData(ctx, meta, data); err != nil {
			return err
		}

		rw.dataRecords += 1

		if rw.dataRecords%1000 == 0 {
			elapsed := time.Now().Sub(rw.started)
			log.Printf("%v records (%v rps)", rw.dataRecords, float64(rw.dataRecords)/elapsed.Seconds())
		}

		batchSize += 1
	}

	if batchSize == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (rw *RecordWalker) loadMeta(ctx context.Context, id int64) (meta *data.MetaRecord, err error) {
	records := []*data.MetaRecord{}
	if err := rw.db.SelectContext(ctx, &records, `SELECT * FROM fieldkit.meta_record WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("error querying for meta: %v", err)
	}

	if len(records) != 1 {
		return nil, fmt.Errorf("error finding meta: %v", id)
	}

	rw.metas[id] = records[0]

	return records[0], nil
}
