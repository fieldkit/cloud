package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

type ExportCsvHandler struct {
	db      *sqlxcache.DB
	files   files.FileArchive
	metrics *logging.Metrics
}

func NewExportCsvHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics) *ExportCsvHandler {
	return &ExportCsvHandler{
		db:      db,
		files:   files,
		metrics: metrics,
	}
}

func (h *ExportCsvHandler) Handle(ctx context.Context, m *messages.ExportCsv) error {
	log := Logger(ctx).Sugar().Named("exporting").With("data_export_id", m.ID).With("user_id", m.UserID)

	log.Infow("processing")

	r, err := repositories.NewExportRepository(h.db)
	if err != nil {
		return err
	}

	de, err := r.QueryByID(ctx, m.ID)
	if err != nil {
		return err
	}

	rawParams := &RawQueryParams{}
	if err := json.Unmarshal(de.Args, rawParams); err != nil {
		return err
	}

	qp, err := rawParams.BuildQueryParams()
	if err != nil {
		return fmt.Errorf("invalid query params: %v", err)
	}

	aggregateName := "1m"
	qp.Aggregate = aggregateName

	log.Infow("query_parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate, "complete", qp.Complete)

	aqp, err := NewAggregateQueryParams(qp, aggregateName, nil)
	if err != nil {
		return err
	}

	dq := NewDataQuerier(h.db)

	queried, err := dq.QueryAggregate(ctx, aqp)
	if err != nil {
		return err
	}

	defer queried.Close()

	reader, writer := io.Pipe()

	var wg sync.WaitGroup
	var archiveError error

	wg.Add(1)

	go func() {
		log.Infow("archiver:ready")

		metadata := make(map[string]string)
		contentType := "text/csv"
		af, err := h.files.Archive(ctx, contentType, metadata, reader)
		if err != nil {
			log.Errorw("archiver:error", "error", err)
			archiveError = err
		} else {
			log.Infow("archiver:done", "key", af.Key)
		}

		wg.Done()
	}()

	for queried.Next() {
		row := &DataRow{}
		if err = queried.StructScan(row); err != nil {
			return err
		}
		fmt.Fprintf(writer, "%v\n", row)
	}

	writer.Close()

	log.Infow("waiting")

	wg.Wait()

	log.Infow("done")

	if archiveError != nil {
		return archiveError
	}

	return nil
}

type DataRow struct {
	Time      data.NumericWireTime `db:"time" json:"time"`
	ID        *int64               `db:"id" json:"-"`
	StationID *int32               `db:"station_id" json:"stationId,omitempty"`
	SensorID  *int64               `db:"sensor_id" json:"sensorId,omitempty"`
	Location  *data.Location       `db:"location" json:"location,omitempty"`
	Value     *float64             `db:"value" json:"value,omitempty"`
	TimeGroup *int32               `db:"time_group" json:"tg,omitempty"`
}
