package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/profile"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

type ExportDataHandler struct {
	db      *sqlxcache.DB
	files   files.FileArchive
	metrics *logging.Metrics
}

func NewExportDataHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics) *ExportDataHandler {
	return &ExportDataHandler{
		db:      db,
		files:   files,
		metrics: metrics,
	}
}

func (h *ExportDataHandler) Handle(ctx context.Context, m *messages.ExportData) error {
	log := Logger(ctx).Sugar().Named("exporting").With("data_export_id", m.ID).With("user_id", m.UserID).With("formatter", m.Formatter)

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

	format := NewCsvFormatter()

	log.Infow("parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations)

	exporter, err := NewExporter(h.db)
	if err != nil {
		return err
	}

	readFunc := func(ctx context.Context, reader io.Reader) error {
		log.Infow("archiver:ready")

		metadata := make(map[string]string)
		contentType := format.MimeType()
		af, err := h.files.Archive(ctx, contentType, metadata, reader)
		if err != nil {
			log.Errorw("archiver:error", "error", err)
			return err
		} else {
			log.Infow("archiver:done", "key", af.Key, "bytes", af.BytesRead)
		}

		now := time.Now()
		size := int32(af.BytesRead)

		de.DownloadURL = &af.URL
		de.CompletedAt = &now
		de.Progress = 100
		de.Size = &size
		if _, err := r.UpdateDataExport(ctx, de); err != nil {
			return err
		}
		return nil
	}

	progressFunc := func(ctx context.Context) error {
		log.Infow("progress")
		return nil
	}

	writeFunc := func(ctx context.Context, writer io.Writer) error {
		criteria := &ExportCriteria{
			Start:    qp.Start,
			End:      qp.End,
			Stations: qp.Stations,
			Sensors:  qp.Sensors,
		}

		err = exporter.Export(ctx, criteria, format, progressFunc, writer)
		if err != nil {
			return err
		}

		return nil
	}

	async := NewAsyncFileWriter(readFunc, writeFunc)
	if err := async.Start(ctx); err != nil {
		return err
	}

	if err := async.Wait(ctx); err != nil {
		return err
	}

	return nil
}
