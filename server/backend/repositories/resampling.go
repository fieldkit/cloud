package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

type Resampler struct {
	Summary      *DataSummary
	MetaFactory  *MetaFactory
	NumberOfBins int32
	bin          *ResamplingBin
}

type Resampled struct {
	Time            time.Time
	Location        []float64
	NumberOfSamples int32
	D               map[string]interface{}
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
		Records: make([]*data.DataRecord, 0, 50),
	}

	r = &Resampler{
		Summary:      summary,
		MetaFactory:  metaFactory,
		NumberOfBins: int32(opts.Resolution),
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

			if r.bin.Number == r.NumberOfBins {
				return nil, fmt.Errorf("sample beyond range (%v)", data.Time)
			}
		}
	}

	r.bin.Insert(data)

	return
}

func (r *Resampler) Close(ctx context.Context) (d *Resampled, err error) {
	log := Logger(ctx).Sugar()

	records := r.bin.Records

	log.Infow("record", "bin", r.bin.Number, "records", len(records))

	for _, record := range records {
		_ = record
	}

	d = &Resampled{
		Time:            r.bin.Time(),
		NumberOfSamples: int32(len(records)),
		Location:        []float64{},
		D:               make(map[string]interface{}),
	}

	r.bin.Clear()

	return
}

type ResamplingBin struct {
	Number  int32
	Start   time.Time
	End     time.Time
	Step    time.Duration
	Records []*data.DataRecord
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

func (b *ResamplingBin) Insert(data *data.DataRecord) {
	b.Records = append(b.Records, data)
}

func (b *ResamplingBin) Clear() {
	b.Records = b.Records[:0]
}

func timeInBetween(start, end time.Time) time.Time {
	d := end.Sub(start).Nanoseconds()
	return start.Add(time.Duration(d / 2))
}
