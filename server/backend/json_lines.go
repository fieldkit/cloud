package backend

import (
	"context"
	"encoding/json"
	"io"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	JSONLinesFormatter = "json-lines"
)

type jsonLinesFormat struct {
}

type jsonLinesWriter struct {
	writer   io.Writer
	readings map[string]*float64
}

func NewJSONLinesFormatter() *jsonLinesFormat {
	return &jsonLinesFormat{}
}

func (f *jsonLinesFormat) MimeType() string {
	return "text/csv"
}

func (f *jsonLinesFormat) CreateWriter(ctx context.Context, writer io.Writer) (ExportWriter, error) {
	return &jsonLinesWriter{
		writer:   writer,
		readings: make(map[string]*float64),
	}, nil
}

func (f *jsonLinesWriter) Begin(ctx context.Context, stations []*data.Station, sensors []*repositories.SensorAndModuleMeta) error {
	return nil
}

type JSONLineStation struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type JSONLineRow struct {
	Time     int64               `db:"time"`
	Uptime   int32               `db:"uptime"`
	Station  JSONLineStation     `db:"station"`
	Location []float64           `db:"location"`
	Readings map[string]*float64 `db:"readings"`
}

func (f *jsonLinesWriter) Row(ctx context.Context, station *data.Station, sensors []*repositories.ReadingValue, row *repositories.FilteredRecord) error {
	f.readings = make(map[string]*float64)
	for _, sam := range sensors {
		if sam != nil {
			f.readings[sam.Sensor.FullKey] = &sam.Value
		}
	}

	line := JSONLineRow{
		Time:   row.Record.Time * 1000,
		Uptime: 0,
		Station: JSONLineStation{
			ID:   station.ID,
			Name: station.Name,
		},
		Location: row.Record.Location,
		Readings: f.readings,
	}

	bytes, err := json.Marshal(line)
	if err != nil {
		return err
	}

	f.writer.Write(bytes)
	f.writer.Write([]byte("\n"))

	return nil
}

func (f *jsonLinesWriter) End(ctx context.Context) error {
	return nil
}
