package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type DataRepository struct {
	Database *sqlxcache.DB
}

func NewDataRepository(database *sqlxcache.DB) (rr *DataRepository, err error) {
	return &DataRepository{Database: database}, nil
}

type SummaryQueryOpts struct {
	DeviceID      string
	DeviceIdBytes []byte
	Internal      bool
	Start         int64
	End           int64
	Resolution    int
	Page          int
	PageSize      int
}

func (r *DataRepository) QueryDevice(ctx context.Context, opts *SummaryQueryOpts) (versions []*Version, err error) {
	log := Logger(ctx).Sugar()

	sr, err := NewStationRepository(r.Database)
	if err != nil {
		return nil, err
	}

	start := time.Unix(opts.Start, 0)
	end := time.Unix(opts.End, 0)

	log.Infow("querying", "device_id", opts.DeviceID, "page_number", opts.Page, "page_size", opts.PageSize, "internal", opts.Internal, "start_unix", opts.Start, "end_unix", opts.End, "start", start, "end", end)

	deviceIdBytes, err := data.DecodeBinaryString(opts.DeviceID)
	if err != nil {
		return nil, err
	}

	station, err := sr.QueryStationByDeviceID(ctx, deviceIdBytes)
	if err != nil {
		return nil, err
	}

	log.Infow("querying", "station_name", station.Name, "station_id", station.ID)

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

	log.Infow("querying for data records")

	rows, err := r.Database.QueryxContext(ctx, `
		SELECT
			r.id, r.provision_id, r.time, r.time, r.number, r.meta, ST_AsBinary(r.location) AS location, r.raw
		FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id = $1) AND (r.time BETWEEN $2 AND $3)
        ORDER BY r.time`, deviceIdBytes, start, end)
	if err != nil {
		return nil, err
	}

	metaIds := make([]int64, 0)
	bin := NewResamplingBin(summary, opts)

	for rows.Next() {
		data := &data.DataRecord{}
		if err = rows.StructScan(data); err != nil {
			return nil, err
		}

		nmetas := len(metaIds)
		if nmetas == 0 || metaIds[nmetas-1] != data.Meta {
			metaIds = append(metaIds, data.Meta)
		}

		d, err := bin.Insert(data)
		if err != nil {
			return nil, err
		}

		if false && d != nil {
			b := d.Bin
			log.Infow("row", "bin", b.Index, "bin_start", b.BinStart, "bin_end", b.BinEnd, "bin_size", d.NumberOfSamples)
		}
	}

	d, err := bin.Close()
	if err != nil {
		return nil, err
	}

	if false && d != nil {
		b := d.Bin
		log.Infow("row", "bin", b.Index, "bin_start", b.BinStart, "bin_end", b.BinEnd, "bin_size", d.NumberOfSamples)
	}

	log.Infow("resampling", "start", summary.Start, "end", summary.End, "resolution", opts.Resolution, "number_metas", len(metaIds))

	return
}

type ResamplingBin struct {
	Index    int32
	Start    time.Time
	End      time.Time
	BinStart time.Time
	BinEnd   time.Time
	Step     time.Duration
	Records  []*data.DataRecord
}

func NewResamplingBin(s *DataSummary, o *SummaryQueryOpts) (b *ResamplingBin) {
	start := s.Start
	end := s.End
	step := time.Duration(end.Sub(start).Nanoseconds() / int64(o.Resolution))

	return &ResamplingBin{
		Index:    0,
		Start:    start,
		End:      end,
		BinStart: start,
		BinEnd:   start.Add(step),
		Step:     step,
		Records:  make([]*data.DataRecord, 0, 1),
	}
}

func (b *ResamplingBin) Contains(i time.Time) bool {
	return (i == b.BinStart || i.After(b.BinStart)) && (i.Before(b.BinEnd) || i == b.End)
}

func (b *ResamplingBin) Next() bool {
	b.BinStart = b.BinEnd
	b.BinEnd = b.BinStart.Add(b.Step)
	b.Index += 1
	return b.BinEnd.Before(b.End) || b.BinEnd == b.End
}

func (b *ResamplingBin) FastForward(i time.Time) bool {
	for {
		if b.Contains(i) {
			return true
		}

		if len(b.Records) != 0 {
			panic("losing records")
		}

		if !b.Next() {
			return false
		}
	}
}

func (b *ResamplingBin) Insert(data *data.DataRecord) (d *Downsampled, err error) {
	if !b.Contains(data.Time) {
		d, err = b.Close()
		if err != nil {
			return nil, err
		}
	}

	if !b.FastForward(data.Time) {
		return nil, fmt.Errorf("unable to find bin")
	}

	b.Records = append(b.Records, data)

	return
}

func (b *ResamplingBin) Clone() *ResamplingBin {
	return &ResamplingBin{
		Index:    b.Index,
		Start:    b.Start,
		End:      b.End,
		BinStart: b.BinStart,
		BinEnd:   b.BinEnd,
		Step:     b.Step,
		Records:  b.Records,
	}
}

func (b *ResamplingBin) Close() (d *Downsampled, err error) {
	d = &Downsampled{
		Bin:             b.Clone(),
		NumberOfSamples: len(b.Records),
	}

	b.Records = b.Records[:0]

	return
}

type Downsampled struct {
	Bin             *ResamplingBin
	NumberOfSamples int
}

type DataSummary struct {
	Start               time.Time `db:"start"`
	End                 time.Time `db:"end"`
	NumberOfDataRecords int64     `db:"number_of_data_records"`
	NumberOfMetaRecords int64     `db:"number_of_meta_records"`
}
