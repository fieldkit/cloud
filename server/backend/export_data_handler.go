package backend

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/profile"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
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

	log.Infow("parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations)

	ir := repositories.NewIngestionRepository(h.db)

	ingestions, err := ir.QueryByStationID(ctx, qp.Stations[0])
	if err != nil {
		return err
	}

	exporter := NewFkbExporter(h.files, h.metrics, nil)

	for _, ingestion := range ingestions {
		if err := exporter.Prepare(ctx, ingestion.URL); err != nil {
			return fmt.Errorf("pass 1: exporting %v failed: %w", ingestion.URL, err)
		}
	}

	for _, ingestion := range ingestions {
		if err := exporter.Export(ctx, ingestion.URL); err != nil {
			return fmt.Errorf("pass 2: exporting %v failed: %w", ingestion.URL, err)
		}
	}

	return nil
}

type FkbExporter struct {
	writer *csv.Writer
	walker *FkbWalker
}

func NewFkbExporter(files files.FileArchive, metrics *logging.Metrics, writer io.Writer) (self *FkbExporter) {
	self = &FkbExporter{
		writer: csv.NewWriter(writer),
	}
	self.walker = NewFkbWalker(files, metrics, self, true)
	return
}

func (e *FkbExporter) Prepare(ctx context.Context, url string) error {
	if _, err := e.walker.WalkUrl(ctx, url); err != nil {
		return err
	}

	return nil
}

func (e *FkbExporter) Export(ctx context.Context, url string) error {
	if _, err := e.walker.WalkUrl(ctx, url); err != nil {
		return err
	}

	return nil
}

func (e *FkbExporter) OnSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	log.Infow("signed-meta", "record_number", signedRecord.Record, "record", rawRecord)

	return nil
}

func (e *FkbExporter) OnMeta(ctx context.Context, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	log.Infow("meta", "record_number", recordNumber, "record", rawRecord)

	return nil
}

func (e *FkbExporter) OnData(ctx context.Context, rawRecord *pb.DataRecord, rawMetaUnused *pb.DataRecord, bytes []byte) error {
	return nil
}

func (e *FkbExporter) OnDone(ctx context.Context) (err error) {
	return nil
}
