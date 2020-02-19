package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type DataSummary struct {
	Start               *time.Time `db:"start"`
	End                 *time.Time `db:"end"`
	NumberOfDataRecords int64      `db:"number_of_data_records"`
	NumberOfMetaRecords int64      `db:"number_of_meta_records"`
}

type DataRepository struct {
	Database *sqlxcache.DB
}

func NewDataRepository(database *sqlxcache.DB) (rr *DataRepository, err error) {
	return &DataRepository{Database: database}, nil
}

type SummaryQueryOpts struct {
	DeviceID   string
	Internal   bool
	Start      int64
	End        int64
	Resolution int
	Page       int
	PageSize   int
}

func (r *DataRepository) queryMetaRecords(ctx context.Context, opts *SummaryQueryOpts) (map[int64]*data.MetaRecord, error) {
	deviceIdBytes, err := data.DecodeBinaryString(opts.DeviceID)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, opts.Start*1000)
	end := time.Unix(0, opts.End*1000)

	mrs := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &mrs, `
	    SELECT m.* FROM fieldkit.meta_record AS m WHERE (m.id IN (
	      SELECT DISTINCT q.meta FROM (
			SELECT r.meta FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id) WHERE (p.device_id = $1) AND (r.time BETWEEN $2 AND $3)
	      ) AS q
	    ))`, deviceIdBytes, start, end); err != nil {
		return nil, err
	}

	metas := make(map[int64]*data.MetaRecord)
	for _, m := range mrs {
		metas[m.ID] = m
	}

	return metas, nil
}

func (r *DataRepository) querySummary(ctx context.Context, opts *SummaryQueryOpts) (*DataSummary, error) {
	deviceIdBytes, err := data.DecodeBinaryString(opts.DeviceID)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, opts.Start*1000)
	end := time.Unix(0, opts.End*1000)

	summaries := make([]*DataSummary, 0)
	if err := r.Database.SelectContext(ctx, &summaries, `
		    SELECT
				MIN(r.time) AS start,
				MAX(r.time) AS end,
				COUNT(*                    ) AS number_of_data_records,
				COUNT(DISTINCT provision_id) AS number_of_meta_records
			FROM
				fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		    WHERE (p.device_id = $1) AND (r.time BETWEEN $2 AND $3)`, deviceIdBytes, start, end); err != nil {
		return nil, err
	}

	if len(summaries) != 1 {
		return nil, fmt.Errorf("unexpected number of summary rows")
	}

	summary := summaries[0]

	return summary, nil
}

func (r *DataRepository) QueryDeviceModulesAndData(ctx context.Context, opts *SummaryQueryOpts) (modulesAndData *ModulesAndData, err error) {
	log := Logger(ctx).Sugar()

	deviceIdBytes, err := data.DecodeBinaryString(opts.DeviceID)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, opts.Start*1000)
	end := time.Unix(0, opts.End*1000)

	log.Infow("summarizing", "device_id", opts.DeviceID, "page_number", opts.Page, "page_size", opts.PageSize, "internal", opts.Internal, "start_ms", opts.Start, "end_ms", opts.End, "start", start, "end", end)

	summary, err := r.querySummary(ctx, opts)
	if err != nil {
		return nil, err
	}
	if summary.NumberOfDataRecords == 0 {
		modulesAndData = &ModulesAndData{
			Modules: make([]*DataMetaModule, 0),
			Data:    make([]*DataRow, 0),
		}
		return
	}

	log.Infow("querying for meta")

	dbMetas, err := r.queryMetaRecords(ctx, opts)
	if err != nil {
		return nil, err
	}

	log.Infow("querying for data")

	rows, err := r.Database.QueryxContext(ctx, `
		SELECT
			r.id, r.provision_id, r.time, r.time, r.number, r.meta, ST_AsBinary(r.location) AS location, r.raw
		FROM
            fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id = $1) AND (r.time BETWEEN $2 AND $3)
        ORDER BY r.time`, deviceIdBytes, start, end)
	if err != nil {
		return nil, err
	}

	metaFactory := NewMetaFactory()
	for _, dbMeta := range dbMetas {
		_, err := metaFactory.Add(dbMeta)
		if err != nil {
			return nil, err
		}
	}

	resampler, err := NewResampler(summary, metaFactory, opts)
	if err != nil {
		return nil, err
	}

	resampled := make([]*Resampled, 0, opts.Resolution)

	for rows.Next() {
		data := &data.DataRecord{}
		if err = rows.StructScan(data); err != nil {
			return nil, err
		}

		d, err := resampler.Insert(ctx, data)
		if err != nil {
			return nil, err
		}

		if d != nil {
			resampled = append(resampled, d)
		}
	}

	d, err := resampler.Close(ctx)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Warnw("final resample was empty")
	} else {
		resampled = append(resampled, d)
	}

	modulesAndData, err = metaFactory.ToModulesAndData(resampled)
	if err != nil {
		return nil, err
	}

	log.Infow("resampling", "start", summary.Start, "end", summary.End, "resolution", opts.Resolution, "number_metas", len(dbMetas), "nsamples", len(resampled))

	return
}
