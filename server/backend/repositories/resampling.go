package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/fieldkit/cloud/server/data"
)

type Resampler struct {
	summary      *DataSummary
	metaFactory  *MetaFactory
	numberOfBins int32
	bin          *ResamplingBin
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
	}
}

func NewResampler(summary *DataSummary, metaFactory *MetaFactory, opts *SummaryQueryOpts) (r *Resampler, err error) {
	// Add to ensure final record ends up in the final bin. Not sure
	// if there's a more elegant way. Millisecond is here because
	// that's closer to the granularity of these samples. Nanosecond
	// is too small.
	epsilon := 1 * time.Millisecond
	end := summary.End.Add(epsilon)
	start := summary.Start
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
		d, err = r.Close(ctx)
		if err != nil {
			return nil, err
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

	row, err := r.metaFactory.Resolve(data)
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
	ids := make([]int64, 0)
	for _, record := range b.Records {
		if len(record.MetaIDs) != 1 {
			panic("record with invalid meta-ids collection")
		}
		if len(ids) == 0 || ids[len(ids)-1] != record.MetaIDs[0] {
			ids = append(ids, record.MetaIDs[0])
		}
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
	return []float64{}
}

func getResampledData(records []*DataRow) (map[string]interface{}, error) {
	all := make(map[string][]float64)
	for _, r := range records {
		for k, v := range r.D {
			if all[k] == nil {
				all[k] = make([]float64, 0, len(records))
			}
			all[k] = append(all[k], float64(v.(float32)))
		}
	}

	d := make(map[string]interface{})

	for k, v := range all {
		mean, err := stats.Mean(v)
		if err != nil {
			return nil, err
		}

		d[k] = mean
	}

	return d, nil
}

func timeInBetween(start, end time.Time) time.Time {
	d := end.Sub(start).Nanoseconds()
	return start.Add(time.Duration(d / 2))
}
