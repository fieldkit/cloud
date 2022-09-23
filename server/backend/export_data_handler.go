package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/profile"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

const (
	SecondsBetweenProgressUpdates = 1
)

type ExportDataHandler struct {
	db        *sqlxcache.DB
	files     files.FileArchive
	metrics   *logging.Metrics
	updatedAt time.Time
}

func NewExportDataHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics) *ExportDataHandler {
	return &ExportDataHandler{
		db:        db,
		files:     files,
		metrics:   metrics,
		updatedAt: time.Now(),
	}
}

func (h *ExportDataHandler) createFormat(format string) (ExportFormat, error) {
	if format == CSVFormatter {
		return NewCsvFormatter(), nil
	}
	if format == JSONLinesFormatter {
		return NewJSONLinesFormatter(), nil
	}
	return nil, fmt.Errorf("unknown format: %v", format)
}

func (h *ExportDataHandler) updateProgress(ctx context.Context, de *data.DataExport, progress float64, message string) error {
	elapsed := time.Now().Sub(h.updatedAt)
	if elapsed.Seconds() < SecondsBetweenProgressUpdates {
		return nil
	}
	h.updatedAt = time.Now()
	de.Progress = progress
	de.Message = &message
	r, err := repositories.NewExportRepository(h.db)
	if err != nil {
		return err
	}
	if _, err := r.UpdateDataExport(ctx, de); err != nil {
		return err
	}
	return nil
}

func (h *ExportDataHandler) Handle(ctx context.Context, m *messages.ExportData) error {
	log := Logger(ctx).Sugar().Named("exporting").With("data_export_id", m.ID).With("user_id", m.UserID).With("formatter", m.Format)

	log.Infow("processing")

	if false {
		defer profile.Start().Stop()
	}

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
		return fmt.Errorf("invalid query params: %w", err)
	}

	format, err := h.createFormat(m.Format)
	if err != nil {
		return err
	}

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

	progressFunc := func(ctx context.Context, progress float64, message string) error {
		return h.updateProgress(ctx, de, progress, message)
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
