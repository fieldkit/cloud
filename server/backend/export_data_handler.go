package backend

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/pkg/profile"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
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
	db                  *sqlxcache.DB
	files               files.FileArchive
	metrics             *logging.Metrics
	updatedAt           time.Time
	bytesExpectedToRead int64
	bytesRead           int64
}

func NewExportDataHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics) *ExportDataHandler {
	return &ExportDataHandler{
		db:                  db,
		files:               files,
		metrics:             metrics,
		updatedAt:           time.Now(),
		bytesExpectedToRead: 0,
		bytesRead:           0,
	}
}

func (h *ExportDataHandler) progress(ctx context.Context, de *data.DataExport, progress WalkProgress) error {
	elapsed := time.Since(h.updatedAt)
	if elapsed.Seconds() < SecondsBetweenProgressUpdates {
		return nil
	}

	h.updatedAt = time.Now()
	h.bytesRead += progress.read

	de.Progress = (float64(h.bytesRead) / float64(h.bytesExpectedToRead)) * 100.0

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

	log.Infow("parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations)

	ir := repositories.NewIngestionRepository(h.db)

	ingestions, err := ir.QueryByStationID(ctx, qp.Stations[0])
	if err != nil {
		return err
	}

	sizeOfSinglePass := int64(0)
	urls := make([]string, 0, len(ingestions))
	for _, ingestion := range ingestions {
		urls = append(urls, ingestion.URL)
		sizeOfSinglePass += ingestion.Size
	}

	h.bytesExpectedToRead = sizeOfSinglePass * 2

	readFunc := func(ctx context.Context, reader io.Reader) error {
		metadata := make(map[string]string)
		af, err := h.files.Archive(ctx, "text/csv", metadata, reader)
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

	progressFunc := func(ctx context.Context, progress WalkProgress) error {
		return h.progress(ctx, de, progress)
	}

	writeFunc := func(ctx context.Context, writer io.Writer) error {
		exporter := NewCsvExporter(h.files, h.metrics, writer, progressFunc)

		if err := exporter.Prepare(ctx, urls); err != nil {
			return fmt.Errorf("preparing: exporting failed: %w", err)
		}

		if err := exporter.Export(ctx, urls); err != nil {
			return fmt.Errorf("writing: exporting failed: %w", err)
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

type CanExport interface {
	Prepare(ctx context.Context, urls []string) error
	Export(ctx context.Context, urls []string) error
}

type JsonLinesExporter struct {
	writer io.Writer
	walker *FkbWalker
}

func NewJsonLinesExporter(files files.FileArchive, metrics *logging.Metrics, writer io.Writer, progress OnWalkProgress) (self *JsonLinesExporter) {
	self = &JsonLinesExporter{
		writer: writer,
	}
	self.walker = NewFkbWalker(files, metrics, self, progress, true)
	return
}

func (e *JsonLinesExporter) Prepare(ctx context.Context, urls []string) error {
	return nil
}

func (e *JsonLinesExporter) Export(ctx context.Context, urls []string) error {
	for _, url := range urls {
		if _, err := e.walker.WalkUrl(ctx, url); err != nil {
			return fmt.Errorf("export %v failed: %w", url, err)
		}
	}

	return nil
}

func (e *JsonLinesExporter) OnSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	log.Infow("signed-meta", "record_number", signedRecord.Record, "record", rawRecord)

	return e.write(ctx, rawRecord)
}

func (e *JsonLinesExporter) OnMeta(ctx context.Context, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	log.Infow("meta", "record_number", recordNumber, "record", rawRecord)

	return e.write(ctx, rawRecord)
}

func (e *JsonLinesExporter) OnData(ctx context.Context, rawRecord *pb.DataRecord, rawMetaUnused *pb.DataRecord, bytes []byte) error {
	for _, sensorGroup := range rawRecord.Readings.SensorGroups {
		for _, reading := range sensorGroup.Readings {
			if calibrated, ok := reading.Calibrated.(*pb.SensorAndValue_CalibratedValue); ok {
				if math.IsNaN(float64(calibrated.CalibratedValue)) {
					reading.Calibrated = &pb.SensorAndValue_CalibratedNull{}
				}
			}
		}
	}

	return e.write(ctx, rawRecord)
}

func (e *JsonLinesExporter) write(ctx context.Context, value interface{}) (err error) {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := e.writer.Write(b); err != nil {
		return err
	}

	if _, err := io.WriteString(e.writer, "\n"); err != nil {
		return err
	}

	return nil
}

func (e *JsonLinesExporter) OnDone(ctx context.Context) (err error) {
	return nil
}

type records struct {
	meta    *pb.DataRecord
	data    *pb.DataRecord
	modules map[string]bool
}

type CsvExporter struct {
	writer        *csv.Writer
	walker        *FkbWalker
	preparing     *preparingCsv
	prepared      *preparingCsv
	row           []string
	records       *records
	bytesRead     int64
	expectedBytes int64
}

type fieldFunc func(*records) string
type optionalFieldFunc func(*records) *string

type csvField struct {
	name string
	get  optionalFieldFunc
}

type fieldSet struct {
	fields []*csvField
}

func newFieldSet() *fieldSet {
	return &fieldSet{
		fields: make([]*csvField, 0),
	}
}

func (p *fieldSet) addField(name string, get optionalFieldFunc) {
	p.fields = append(p.fields, &csvField{name: name, get: get})
}

type preparingCsv struct {
	fields  *fieldSet
	metas   map[int64]*pb.DataRecord
	modules map[string]*fieldSet
	order   []string
}

func (p *preparingCsv) addField(name string, get fieldFunc) {
	p.fields.addField(name, func(r *records) *string {
		value := get(r)
		return &value
	})
}

func NewCsvExporter(files files.FileArchive, metrics *logging.Metrics, writer io.Writer, progress OnWalkProgress) (self *CsvExporter) {
	self = &CsvExporter{
		writer:        csv.NewWriter(writer),
		records:       &records{},
		row:           nil,
		bytesRead:     0,
		expectedBytes: 0,
		preparing: &preparingCsv{
			fields:  newFieldSet(),
			metas:   make(map[int64]*pb.DataRecord),
			modules: make(map[string]*fieldSet),
			order:   make([]string, 0),
		},
	}
	self.walker = NewFkbWalker(files, metrics, self, progress, true)
	return
}

func (e *CsvExporter) Prepare(ctx context.Context, urls []string) error {
	e.preparing.addField("unix_time", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Time)
	})
	e.preparing.addField("time", func(r *records) string {
		t := time.Unix(r.data.Readings.Time, 0)
		return fmt.Sprintf("%v", t)
	})
	e.preparing.addField("data_record", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Reading)
	})
	e.preparing.addField("meta_record", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Meta)
	})
	e.preparing.addField("uptime", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Uptime)
	})
	e.preparing.addField("gps", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Location.Fix)
	})
	e.preparing.addField("latitude", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Location.Latitude)
	})
	e.preparing.addField("longitude", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Location.Longitude)
	})
	e.preparing.addField("altitude", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Location.Altitude)
	})
	e.preparing.addField("gps_time", func(r *records) string {
		return fmt.Sprintf("%v", r.data.Readings.Location.Time)
	})
	e.preparing.addField("note", func(r *records) string {
		return ""
	})

	for _, url := range urls {
		if _, err := e.walker.WalkUrl(ctx, url); err != nil {
			return fmt.Errorf("prepare %v failed: %w", url, err)
		}
	}

	e.prepared = e.preparing
	e.preparing = nil

	return nil
}

func (e *CsvExporter) Export(ctx context.Context, urls []string) error {
	log := Logger(ctx).Sugar()

	if err := e.compactFieldSets(ctx); err != nil {
		return err
	}

	log.Infow("prepared", "conflicts", e.prepared.conflicts)
	log.Infow("prepared", "module_ids", e.prepared.order)
	log.Infow("prepared", "compacted", e.prepared.compacted)

	header := make([]string, 0, len(e.prepared.fields.fields))
	for _, field := range e.prepared.fields.fields {
		header = append(header, field.name)
	}
	for _, fs := range e.prepared.compacted {
		for _, field := range fs.fields {
			header = append(header, field.name)
		}
	}
	if err := e.writer.Write(header); err != nil {
		return err
	}

	for _, url := range urls {
		if _, err := e.walker.WalkUrl(ctx, url); err != nil {
			return fmt.Errorf("export %v failed: %w", url, err)
		}
	}

	return nil
}

func (e *CsvExporter) prepare(ctx context.Context, rawRecord *pb.DataRecord) error {
	log := Logger(ctx).Sugar()

	if rawRecord.Metadata != nil {
		for moduleIndex, module := range rawRecord.Modules {
			id := hex.EncodeToString(module.Id)
			if _, ok := e.preparing.modules[id]; ok {
				continue
			}

			log.Infow("module", "module_id", id)

			fields := newFieldSet()

			checkForModule := func(get fieldFunc) optionalFieldFunc {
				return (func(id string) optionalFieldFunc {
					return func(r *records) *string {
						if _, ok := r.modules[id]; ok {
							value := get(r)
							return &value
						} else {
							return nil
						}
					}
				})(id)
			}

			fields.addField("module_index", checkForModule(func(r *records) string {
				return fmt.Sprintf("%d", moduleIndex)
			}))
			fields.addField("module_position", checkForModule(func(r *records) string {
				return fmt.Sprintf("%d", module.Position)
			}))
			fields.addField("module_name", checkForModule(func(r *records) string {
				return module.Name
			}))
			fields.addField("module_id", checkForModule(func(r *records) string {
				return hex.EncodeToString(module.Id)
			}))

			for sensorIndex, sensor := range module.Sensors {
				fields.addField(sensor.Name, (func(moduleIndex, sensorIndex int) optionalFieldFunc {
					return checkForModule(func(r *records) string {
						if moduleIndex >= len(r.data.Readings.SensorGroups) {
							return ""
						}

						sensorGroup := r.data.Readings.SensorGroups[moduleIndex]
						if sensorIndex >= len(sensorGroup.Readings) {
							return ""
						}
						sensor := sensorGroup.Readings[sensorIndex]
						if sensor.GetCalibratedNull() {
							return ""
						}
						return fmt.Sprintf("%v", sensor.GetCalibratedValue())
					})
				})(moduleIndex, sensorIndex))

				fields.addField(fmt.Sprintf("%s_raw_v", sensor.Name), (func(moduleIndex, sensorIndex int) optionalFieldFunc {
					return checkForModule(func(r *records) string {
						if moduleIndex >= len(r.data.Readings.SensorGroups) {
							return ""
						}

						sensorGroup := r.data.Readings.SensorGroups[moduleIndex]
						if sensorIndex >= len(sensorGroup.Readings) {
							return ""
						}
						sensor := sensorGroup.Readings[sensorIndex]
						if sensor.GetUncalibratedNull() {
							return ""
						}
						return fmt.Sprintf("%v", sensor.GetUncalibratedValue())
					})
				})(moduleIndex, sensorIndex))
			}

			e.preparing.modules[id] = fields
			e.preparing.order = append(e.preparing.order, id)
		}
	}

	return nil
}

func (e *CsvExporter) OnSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
	return e.OnMeta(ctx, int64(signedRecord.Record), rawRecord, bytes)
}

func (e *CsvExporter) OnMeta(ctx context.Context, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
	if e.preparing != nil {
		e.preparing.metas[recordNumber] = rawRecord

		if err := e.prepare(ctx, rawRecord); err != nil {
			return err
		}
	}

	return nil
}

func (e *CsvExporter) OnData(ctx context.Context, rawRecord *pb.DataRecord, rawMetaUnused *pb.DataRecord, bytes []byte) error {
	if e.prepared != nil {
		meta, ok := e.prepared.metas[int64(rawRecord.Readings.Meta)]
		if !ok {
			return fmt.Errorf("missing meta: %v", rawRecord.Readings.Meta)
		}

		e.records.data = rawRecord
		if e.records.modules == nil || e.records.meta != meta {
			e.records.meta = meta
			e.records.modules = make(map[string]bool)
			for _, module := range e.records.meta.Modules {
				e.records.modules[hex.EncodeToString(module.Id)] = true
			}
		}

		e.row = make([]string, 0, len(e.prepared.fields.fields))

		for _, field := range e.prepared.fields.fields {
			value := field.get(e.records)
			if value != nil {
				e.row = append(e.row, *value)
			} else {
				e.row = append(e.row, "")
			}
		}

		for _, id := range e.prepared.order {
			for _, field := range e.prepared.modules[id].fields {
				value := field.get(e.records)
				if value != nil {
					e.row = append(e.row, *value)
				} else {
					e.row = append(e.row, "")
				}
			}
		}

		if err := e.writer.Write(e.row); err != nil {
			return err
		}
	}

	return nil
}

func (e *CsvExporter) OnDone(ctx context.Context) (err error) {
	return nil
}

func ExportQueryParams(de *data.DataExport) (*QueryParams, error) {
	rawParams := &RawQueryParams{}
	if err := json.Unmarshal(de.Args, rawParams); err != nil {
		return nil, err
	}

	return rawParams.BuildQueryParams()
}
