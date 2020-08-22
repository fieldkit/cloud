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
	writer  *csv.Writer
	columns []string
}

func NewCsvFormatter() *csvFormat {
	return &csvFormat{}
}

func (f *csvFormat) MimeType() string {
	return "text/csv"
}

func (f *csvFormat) CreateWriter(ctx context.Context, writer io.Writer) (ExportWriter, error) {
	return &csvWriter{
		writer:  csv.NewWriter(writer),
		columns: make([]string, NumberOfColumns),
	}, nil
}

func (f *csvWriter) Begin(ctx context.Context, stations []*data.Station, sensors []*repositories.SensorAndModuleMeta) error {
	f.columns[0] = "time_js"
	f.columns[1] = "uptime"
	f.columns[2] = "station_id"
	f.columns[3] = "station_name"
	f.columns[4] = "latitude"
	f.columns[5] = "longitude"

	for _, sensor := range sensors {
		f.columns = append(f.columns, sensor.Sensor.FullKey)
	}

	f.writer.Write(f.columns)

	return nil
}

func (f *csvWriter) Row(ctx context.Context, station *data.Station, sensors []*repositories.ReadingValue, row *repositories.FilteredRecord) error {
	f.columns[0] = fmt.Sprintf("%v", row.Record.Time*1000)
	f.columns[1] = fmt.Sprintf("%v", 0)
	f.columns[2] = fmt.Sprintf("%v", station.ID)
	f.columns[3] = station.Name
	if row.Record.Location != nil {
		l := row.Record.Location
		f.columns[4] = fmt.Sprintf("%v", l[0])
		f.columns[5] = fmt.Sprintf("%v", l[1])
	} else {
		f.columns[4] = ""
		f.columns[5] = ""
	}

	for i, sam := range sensors {
		if sam != nil {
			f.columns[6+i] = fmt.Sprintf("%v", sam.Value)
		} else {
			f.columns[6+i] = ""
		}
	}

	f.writer.Write(f.columns)

	return nil
}

func (f *csvWriter) End(ctx context.Context) error {
	f.writer.Flush()
	return nil
}
