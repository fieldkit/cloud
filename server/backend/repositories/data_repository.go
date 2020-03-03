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
	Interval   int64
}

func (r *DataRepository) queryMetaRecords(ctx context.Context, deviceIdBytes []byte, start, end time.Time) (map[int64]*data.MetaRecord, error) {
	mrs := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &mrs, `
	    SELECT m.* FROM fieldkit.meta_record AS m WHERE (m.id IN (
	      SELECT DISTINCT q.meta FROM (
			SELECT r.meta FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id) WHERE (p.device_id = $1) AND (timezone('UTC', r.time) BETWEEN $2 AND $3)
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

	start := time.Unix(0, opts.Start*int64(time.Millisecond))
	end := time.Unix(0, opts.End*int64(time.Millisecond))

	summaries := make([]*DataSummary, 0)
	if err := r.Database.SelectContext(ctx, &summaries, `
		    SELECT
				MIN(r.time) AS start,
				MAX(r.time) AS end,
				COUNT(*              ) AS number_of_data_records,
				COUNT(DISTINCT r.meta) AS number_of_meta_records
			FROM
				fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		    WHERE (p.device_id = $1) AND (timezone('UTC', r.time) BETWEEN $2 AND $3)`, deviceIdBytes, start, end); err != nil {
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

	start := time.Unix(0, opts.Start*int64(time.Millisecond))
	end := time.Unix(0, opts.End*int64(time.Millisecond))

	if start.After(end) {
		return nil, fmt.Errorf("malformed query: times")
	}

	log.Infow("summarizing", "device_id", opts.DeviceID, "page_number", opts.Page, "page_size", opts.PageSize, "internal", opts.Internal, "start_ms", opts.Start, "end_ms", opts.End, "start", start, "end", end, "interval", opts.Interval)

	summary, err := r.querySummary(ctx, opts)
	if err != nil {
		return nil, err
	}

	if summary.NumberOfDataRecords == 0 {
		log.Infow("empty")
		modulesAndData = &ModulesAndData{
			Modules:    make([]*DataMetaModule, 0),
			Data:       make([]*DataRow, 0),
			Statistics: nil,
		}
		return
	}

	log.Infow("summary", "start", summary.Start, "end", summary.End, "number_data_records", summary.NumberOfDataRecords, "number_meta", summary.NumberOfMetaRecords)

	if opts.Interval > 0 {
		start = summary.End.Add(time.Duration(-opts.Interval) * time.Minute)
		end = *summary.End
	}

	log.Infow("querying", "start", start, "end", end)

	dbMetas, err := r.queryMetaRecords(ctx, deviceIdBytes, start, end)
	if err != nil {
		return nil, err
	}

	log.Infow("querying", "number_metas", len(dbMetas))

	rows, err := r.Database.QueryxContext(ctx, `
		SELECT
			r.id, r.provision_id, r.time, r.time, r.number, r.meta, ST_AsBinary(r.location) AS location, r.raw
		FROM
            fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id = $1) AND (timezone('UTC', r.time) BETWEEN $2 AND $3)
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
	if d != nil {
		resampled = append(resampled, d)
	}

	modulesAndData, err = metaFactory.ToModulesAndData(resampled, summary)
	if err != nil {
		return nil, err
	}

	log.Infow("resampling", "start", summary.Start, "end", summary.End, "resolution", opts.Resolution, "number_metas", len(dbMetas), "nsamples", len(resampled))

	return
}
