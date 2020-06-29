package main

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/data"
)

type RecordWalker struct {
	metas       map[int64]*data.MetaRecord
	db          *sqlxcache.DB
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

func (rw *RecordWalker) WalkStation(ctx context.Context, id int32, visitor RecordVisitor) error {
	rows, err := rw.db.QueryxContext(ctx, `
		SELECT r.id, r.provision_id, r.time, r.number, r.meta_record_id, ST_AsBinary(r.location) AS location, r.raw, r.pb
		FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1))
        ORDER BY r.time`, id)
	if err != nil {
		return err
	}

	defer rows.Close()

	return rw.walkQuery(ctx, rows, visitor)
}

func (rw *RecordWalker) walkQuery(ctx context.Context, rows *sqlx.Rows, visitor RecordVisitor) error {
	data := &data.DataRecord{}

	for rows.Next() {
		if err := rows.StructScan(data); err != nil {
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
	}

	if err := visitor.VisitEnd(ctx); err != nil {
		return err
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
