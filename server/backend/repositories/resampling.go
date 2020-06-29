package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/fieldkit/cloud/server/data"
)

const (
	ResampledKey = "_resampled"
	FiltersKey   = "_filters"
)

type Resampler struct {
	summary      *DataSummary
	metaFactory  *MetaFactory
	numberOfBins int32
	bin          *ResamplingBin
}

func NewResampler(summary *DataSummary, metaFactory *MetaFactory, opts *SummaryQueryOpts) (r *Resampler, err error) {
	// Add to ensure final record ends up in the final bin. Not sure
	// if there's a more elegant way. Millisecond is here because
	// that's closer to the granularity of these samples. Nanosecond
	// is too small.
	epsilon := 1 * time.Millisecond
	end := summary.End.Add(epsilon)
	start := *summary.Start
	step := time.Duration(end.Sub(start).Nanoseconds() / int64(opts.Resolution))

	bin := &ResamplingBin{
		Number:  0,
		Start:   start,
		End:     start.Add(step),
		Step:    step,
		Records: make([]*FilteredRecord, 0, 1000),
	}

	r = &Resampler{
		summary:      summary,
		metaFactory:  metaFactory,
		numberOfBins: int32(opts.Resolution),
		bin:          bin,
	}

	return
}

func (r *Resampler) Insert(ctx context.Context, data *data.DataRecord) (d *Resampled, err error) {
	if !r.bin.Contains(data.Time) {
		if len(r.bin.Records) > 0 {
			d, err = r.Close(ctx)
			if err != nil {
				return nil, err
			}
		}

		for {
			r.bin.Next()

			if r.bin.Contains(data.Time) {
				break
			}

			if r.bin.Number == r.numberOfBins {
				return nil, fmt.Errorf("sample beyond range (%v)", data.Time)
			}
		}
	}

	row, err := r.metaFactory.Resolve(ctx, data, false, false)
	if err != nil {
		return nil, err
	}

	if row != nil {
		r.bin.Insert(row)
	}

	return
}

func (r *Resampler) Close(ctx context.Context) (d *Resampled, err error) {
	records := FilteredRecordCollection{
		records: r.bin.Records,
	}
	metaIDs := r.bin.MetaIDs()
	location := records.firstLocation()
	data, err := records.mean()
	if err != nil {
		return nil, err
	}

	d = &Resampled{
		Time:            r.bin.Time(),
		NumberOfSamples: int32(len(r.bin.Records)),
		MetaIDs:         metaIDs,
		Location:        location,
		D:               data,
	}

	if false {
		log := Logger(ctx).Sugar()
		log.Infow("record", "bin", r.bin.Number, "records", len(r.bin.Records), "location", location, "data", data)
	}

	r.bin.Clear()

	return
}

type ResamplingBin struct {
	Number  int32
	Start   time.Time
	End     time.Time
	Step    time.Duration
	Records []*FilteredRecord
}

func (b *ResamplingBin) MetaIDs() []int64 {
	idsMap := make(map[int64]bool)
	for _, filteredRecord := range b.Records {
		for _, reading := range filteredRecord.Record.Readings {
			idsMap[reading.MetaRecordID] = true
		}
	}
	ids := make([]int64, 0)
	for k, _ := range idsMap {
		ids = append(ids, k)
	}
	return ids
}

func (b *ResamplingBin) Time() time.Time {
	return timeInBetween(b.Start, b.End)
}

func (b *ResamplingBin) Contains(i time.Time) bool {
	return (i.After(b.Start) || i == b.Start) && i.Before(b.End)
}

func (b *ResamplingBin) Next() {
	if len(b.Records) != 0 {
		panic("iterating non-empty bin")
	}

	b.Start = b.End
	b.End = b.Start.Add(b.Step)
	b.Number += 1
}

func (b *ResamplingBin) Insert(record *FilteredRecord) {
	b.Records = append(b.Records, record)
}

func (b *ResamplingBin) Clear() {
	b.Records = b.Records[:0]
}

type FilteredRecordCollection struct {
	records []*FilteredRecord
}

func (c *FilteredRecordCollection) firstLocation() []float64 {
	for _, r := range c.records {
		if r.Record.Location != nil && len(r.Record.Location) > 0 {
			return r.Record.Location
		}
	}
	return nil
}

func (c *FilteredRecordCollection) mean() (map[string]interface{}, error) {
	filterAuditLog := NewFilterAuditLog()
	all := make(map[string][]float64)
	ids := make([]int64, 0)

	for _, r := range c.records {
		filterAuditLog.Include(r)

		for k, reading := range r.Record.Readings {
			if !r.Filters.IsFiltered(k) {
				if all[k] == nil {
					all[k] = make([]float64, 0, len(c.records))
				}
				all[k] = append(all[k], reading.Value)
			}
		}

		ids = append(ids, r.Record.ID)
	}

	d := make(map[string]interface{})
	for k, v := range all {
		mean, err := stats.Mean(v)
		if err != nil {
			return nil, err
		}
		d[k] = mean
	}

	d[FiltersKey] = filterAuditLog
	d[ResampledKey] = ResampleInfo{
		Size: int32(len(c.records)),
		IDs:  ids,
	}

	return d, nil
}

func timeInBetween(start, end time.Time) time.Time {
	d := end.Sub(start).Nanoseconds()
	return start.Add(time.Duration(d / 2))
}
