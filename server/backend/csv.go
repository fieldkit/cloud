package backend

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	CSVFormatter    = "csv"
	NumberOfColumns = 6
)

type csvFormat struct {
}

type csvWriter struct {
	writer *csv.Writer
}

func NewCsvFormatter() *csvFormat {
	return &csvFormat{}
}

func (f *csvFormat) MimeType() string {
	return "text/csv"
}

func (f *csvFormat) CreateWriter(ctx context.Context, writer io.Writer) (ExportWriter, error) {
	return &csvWriter{
		writer: csv.NewWriter(writer),
	}, nil
}

func (f *csvWriter) Begin(ctx context.Context, stations []*data.Station, sensors []*repositories.SensorAndModuleMeta) error {
	cols := make([]string, NumberOfColumns)
	cols[0] = "time_js"
	cols[1] = "uptime"
	cols[2] = "station_id"
	cols[3] = "station_name"
	cols[4] = "latitude"
	cols[5] = "longitude"

	for _, sensor := range sensors {
		cols = append(cols, sensor.Sensor.FullKey)
	}

	f.writer.Write(cols)

	return nil
}

func (f *csvWriter) Row(ctx context.Context, station *data.Station, sensors []*repositories.ReadingValue, row *repositories.FilteredRecord) error {
	cols := make([]string, NumberOfColumns)
	cols[0] = fmt.Sprintf("%v", row.Record.Time*1000)
	cols[1] = fmt.Sprintf("%v", 0)
	cols[2] = fmt.Sprintf("%v", station.ID)
	cols[3] = fmt.Sprintf("%v", station.Name)
	if row.Record.Location != nil {
		l := row.Record.Location
		cols[4] = fmt.Sprintf("%v", l[0])
		cols[5] = fmt.Sprintf("%v", l[1])
	}

	for _, sam := range sensors {
		if sam != nil {
			cols = append(cols, fmt.Sprintf("%v", sam.Value))
		} else {
			cols = append(cols, "")
		}
	}

	f.writer.Write(cols)

	return nil
}

func (f *csvWriter) End(ctx context.Context) error {
	f.writer.Flush()
	return nil
}
