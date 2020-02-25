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
		Records: make([]*DataRow, 0, 50),
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

	row, err := r.metaFactory.Resolve(ctx, data)
	if err != nil {
		return nil, err
	}

	r.bin.Insert(row)

	return
}

func (r *Resampler) Close(ctx context.Context) (d *Resampled, err error) {
	log := Logger(ctx).Sugar()

	records := r.bin.Records
	metaIDs := r.bin.MetaIDs()
	location := getResampledLocation(records)
	data, err := getResampledData(records)
	if err != nil {
		return nil, err
	}

	d = &Resampled{
		Time:            r.bin.Time(),
		NumberOfSamples: int32(len(records)),
		MetaIDs:         metaIDs,
		Location:        location,
		D:               data,
	}

	if false {
		log.Infow("record", "bin", r.bin.Number, "records", len(records), "location", location, "data", data)
	}

	r.bin.Clear()

	return
}

type ResamplingBin struct {
	Number  int32
	Start   time.Time
	End     time.Time
	Step    time.Duration
	Records []*DataRow
}

func (b *ResamplingBin) MetaIDs() []int64 {
	idsMap := make(map[int64]bool)
	for _, record := range b.Records {
		if len(record.MetaIDs) != 1 {
			panic("record with invalid meta-ids collection")
		}
		for _, id := range record.MetaIDs {
			idsMap[id] = true
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

func (b *ResamplingBin) Insert(data *DataRow) {
	b.Records = append(b.Records, data)
}

func (b *ResamplingBin) Clear() {
	b.Records = b.Records[:0]
}

func getResampledLocation(records []*DataRow) []float64 {
	for _, r := range records {
		if r.Location != nil && len(r.Location) > 0 {
			return r.Location
		}
	}
	return nil
}

func shouldInclude(row *DataRow) bool {
	return true
}

func getResampledData(records []*DataRow) (map[string]interface{}, error) {
	start := time.Time{}
	end := time.Time{}
	ids := make([]int64, 0)
	filtered := make([]int64, 0)

	all := make(map[string][]float64)
	for _, r := range records {
		if !shouldInclude(r) {
			filtered = append(filtered, r.ID)
			continue
		}

		for k, v := range r.D {
			if all[k] == nil {
				all[k] = make([]float64, 0, len(records))
			}
			all[k] = append(all[k], float64(v.(float32)))
		}

		time := time.Unix(r.Time, 0)
		if start.IsZero() || time.Before(start) {
			start = time
		}
		if end.IsZero() || time.After(end) {
			end = time
		}

		ids = append(ids, r.ID)
	}

	d := make(map[string]interface{})

	for k, v := range all {
		mean, err := stats.Mean(v)
		if err != nil {
			return nil, err
		}

		d[k] = mean
	}

	d[ResampledKey] = ResampleInfo{
		Size:     int32(len(records)),
		IDs:      ids,
		Filtered: filtered,
		Start:    start,
		End:      end,
	}

	return d, nil
}

func timeInBetween(start, end time.Time) time.Time {
	d := end.Sub(start).Nanoseconds()
	return start.Add(time.Duration(d / 2))
}

type ResampleInfo struct {
	Size     int32     `json:"size"`
	IDs      []int64   `json:"ids"`
	Filtered []int64   `json:"fids"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type Resampled struct {
	NumberOfSamples int32
	MetaIDs         []int64
	Time            time.Time
	Location        []float64
	D               map[string]interface{}
}

func (r *Resampled) ToDataRow() *DataRow {
	return &DataRow{
		ID:       0,
		MetaIDs:  r.MetaIDs,
		Time:     r.Time.Unix(),
		Location: r.Location,
		D:        r.D,
		Filtered: false,
	}
}
