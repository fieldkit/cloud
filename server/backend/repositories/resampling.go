package repositories

import (
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

type ResamplingBin struct {
	Index    int32
	Start    time.Time
	End      time.Time
	BinStart time.Time
	BinEnd   time.Time
	Step     time.Duration
	Records  []*data.DataRecord
	Summary  *DataSummary
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
		Summary:  s,
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
		// NOTE Should this be the middle of the bin's time range?
		Time:            b.Start,
		Location:        []float64{},
		NumberOfSamples: int32(len(b.Records)),
		D:               make(map[string]interface{}),
	}

	b.Records = b.Records[:0]

	return
}

type Downsampled struct {
	Time            time.Time
	Location        []float64
	NumberOfSamples int32
	D               map[string]interface{}
}
