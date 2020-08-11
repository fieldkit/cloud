package backend

import (
	"context"
	_ "encoding/json"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"

	_ "github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/repositories"
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
	log := Logger(ctx).Sugar().With("data_export_id", m.ID)

	log.Infow("processing")

	r, err := repositories.NewExportRepository(h.db)
	if err != nil {
		return err
	}

	de, err := r.QueryByID(ctx, m.ID)
	if err != nil {
		return err
	}

	_ = de

	/*
		rawParams := &data.RawQueryParams{}
		if err := json.Unmarshal(de.Args, rawParams); err != nil {
			return err
		}

		qp, err := rawParams.BuildQueryParams()
		if err != nil {
			return nil, csvService.MakeBadRequest(err)
		}

		log.Infow("query_parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate)

		aggregateName := "1m"
		aqp, err := NewAggregateQueryParams(qp, aggregateName, nil)
		if err != nil {
			return nil, err
		}

		dq := NewDataQuerier(c.db)

		queried, err := dq.QueryAggregate(ctx, aqp)
		if err != nil {
			return nil, err
		}

		defer queried.Close()

		for queried.Next() {
			row := &DataRow{}
			if err = queried.StructScan(row); err != nil {
				return nil, err
			}
		}
	*/
	return nil
}
