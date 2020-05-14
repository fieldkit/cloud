package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/logging"
)

type RecordAdder struct {
	verbose    bool
	handler    RecordHandler
	database   *sqlxcache.DB
	files      files.FileArchive
	metrics    *logging.Metrics
	statistics *newRecordStatistics
}

type ParsedRecord struct {
	SignedRecord *pb.SignedRecord
	DataRecord   *pb.DataRecord
}

func NewRecordAdder(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, handler RecordHandler, verbose bool) (ra *RecordAdder) {
	return &RecordAdder{
		verbose:    verbose,
		database:   db,
		files:      files,
		metrics:    metrics,
		handler:    handler,
		statistics: &newRecordStatistics{},
	}
}

func (ra *RecordAdder) tryParseSignedRecord(sr *pb.SignedRecord, dataRecord *pb.DataRecord) (err error) {
	err = proto.Unmarshal(sr.Data, dataRecord)
	if err == nil {
		return nil
	}

	buffer := proto.NewBuffer(sr.Data)
	size, err := buffer.DecodeVarint()
	if err != nil {
		return
	}

	if size > uint64(len(sr.Data)) {
		return fmt.Errorf("bad length prefix in signed record: %d", size)
	}

	err = buffer.Unmarshal(dataRecord)
	if err != nil {
		return
	}

	return
}

func (ra *RecordAdder) tryFindStation(ctx context.Context, i *data.Ingestion) (*data.Station, error) {
	r, err := repositories.NewStationRepository(ra.database)
	if err != nil {
		return nil, err
	}
	return r.TryQueryStationByDeviceID(ctx, i.DeviceID)
}

func (ra *RecordAdder) findProvision(ctx context.Context, i *data.Ingestion) (*data.Provision, error) {
	// TODO Verify we have a valid generation.

	provisions := []*data.Provision{}
	if err := ra.database.SelectContext(ctx, &provisions, `
		SELECT p.* FROM fieldkit.provision AS p WHERE p.device_id = $1 AND p.generation = $2
		`, i.DeviceID, i.GenerationID); err != nil {
		return nil, err
	}

	if len(provisions) == 1 {
		return provisions[0], nil
	}

	provision := &data.Provision{
		Created:      time.Now(),
		Updated:      time.Now(),
		DeviceID:     i.DeviceID,
		GenerationID: i.GenerationID,
	}

	if err := ra.database.NamedGetContext(ctx, &provision.ID, `
		INSERT INTO fieldkit.provision (device_id, generation, created, updated)
		VALUES (:device_id, :generation, :created, :updated) ON CONFLICT (device_id, generation)
		DO UPDATE SET updated = NOW() RETURNING id
		`, provision); err != nil {
		return nil, err
	}

	return provision, nil
}

func (ra *RecordAdder) Handle(ctx context.Context, i *data.Ingestion, pr *ParsedRecord) (warning error, fatal error) {
	log := Logger(ctx).Sugar()
	verboseLog := logging.OnlyLogIf(log, ra.verbose)

	provision, err := ra.findProvision(ctx, i)
	if err != nil {
		return nil, err
	}

	recordRepository, err := repositories.NewRecordRepository(ra.database)
	if err != nil {
		return nil, err
	}

	if pr.SignedRecord != nil {
		metaRecord, err := recordRepository.AddMetaRecord(ctx, provision, i, pr.SignedRecord, pr.DataRecord)
		if err != nil {
			return nil, err
		}

		if err := ra.handler.OnMeta(ctx, provision, pr.DataRecord, metaRecord); err != nil {
			return nil, err
		}
	} else if pr.DataRecord != nil {
		dataRecord, metaRecord, err := recordRepository.AddDataRecord(ctx, provision, i, pr.DataRecord)
		if err != nil {
			if err == repositories.ErrMalformedRecord {
				verboseLog.Infow("data reading missing readings", "record", pr.DataRecord)
				ra.metrics.DataErrorsUnknown()
				return err, nil
			}
			if err == repositories.ErrMetaMissing {
				verboseLog.Errorw("error finding meta record", "provision_id", provision.ID, "meta_record_number", pr.DataRecord.Readings.Meta)
				ra.metrics.DataErrorsMissingMeta()
				return err, nil
			}
			return nil, err
		}

		ra.statistics.addTime(dataRecord.Time)

		if err := ra.handler.OnData(ctx, provision, pr.DataRecord, dataRecord, metaRecord); err != nil {
			return nil, err
		}
	}

	return
}

type WriteInfo struct {
	TotalRecords int64
	DataRecords  int64
	MetaRecords  int64
	MetaErrors   int64
	DataErrors   int64
	StationID    *int32
}

func (ra *RecordAdder) WriteRecords(ctx context.Context, i *data.Ingestion) (info *WriteInfo, err error) {
	log := Logger(ctx).Sugar().With("ingestion_id", i.ID, "device_id", i.DeviceID, "user_id", i.UserID, "blocks", i.Blocks)

	log.Infow("file", "file_url", i.URL, "file_stamp", i.Time, "file_id", i.UploadID, "file_size", i.Size, "type", i.Type)

	station, err := ra.tryFindStation(ctx, i)
	if err != nil {
		return nil, err
	}

	reader, err := ra.files.OpenByURL(ctx, i.URL)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	meta := false
	data := false

	totalRecords := 0
	dataProcessed := 0
	dataErrors := 0
	metaProcessed := 0
	metaErrors := 0
	records := 0

	unmarshalFunc := UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var unmarshalError error
		var dataRecord pb.DataRecord
		var signedRecord pb.SignedRecord

		totalRecords += 1

		if data || (!data && !meta) {
			err := proto.Unmarshal(b, &dataRecord)
			if err != nil {
				if data { // If we expected this record, return the error
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("error parsing data record: %v", err)
				}
				unmarshalError = err
			} else {
				data = true
				warning, fatal := ra.Handle(ctx, i, &ParsedRecord{DataRecord: &dataRecord})
				if fatal != nil {
					return nil, fatal
				}
				if warning == nil {
					dataProcessed += 1
					records += 1
				} else {
					dataErrors += 1
					if records > 0 {
						log.Infow("processed", "record_run", records)
						records = 0
					}
				}
			}
		}

		if meta || (!data && !meta) {
			err := proto.Unmarshal(b, &signedRecord)
			if err != nil {
				if meta { // If we expected this record, return the error
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("error parsing signed record: %v", err)
				}
				unmarshalError = err
			} else {
				meta = true
				err := ra.tryParseSignedRecord(&signedRecord, &dataRecord)
				if err != nil {
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("error parsing signed record: %v", err)
				}

				warning, fatal := ra.Handle(ctx, i, &ParsedRecord{
					SignedRecord: &signedRecord,
					DataRecord:   &dataRecord,
				})
				if fatal != nil {
					return nil, fatal
				}
				if warning == nil {
					metaProcessed += 1
					records += 1
				} else {
					metaErrors += 1
					if records > 0 {
						log.Infow("processed", "record_run", records)
						records = 0
					}
				}
			}
		}

		// Parsing either one failed, otherwise this would have been set. So return the error.
		if !meta && !data {
			return nil, unmarshalError
		}

		return nil, nil
	})

	_, _, err = ReadLengthPrefixedCollection(ctx, MaximumDataRecordLength, reader, unmarshalFunc)

	withInfo := log.With("meta_processed", metaProcessed,
		"data_processed", dataProcessed,
		"meta_errors", metaErrors,
		"data_errors", dataErrors,
		"record_run", records,
		"start_human", prettyTime(ra.statistics.start),
		"end_human", prettyTime(ra.statistics.end))

	var stationID *int32
	if station != nil {
		withInfo = withInfo.With("station_name", station.Name)
		stationID = &station.ID
	}

	if err := ra.handler.OnDone(ctx); err != nil {
		return nil, err
	}

	withInfo.Infow("processed")

	if err != nil {
		newErr := fmt.Errorf("error reading collection: %v", err)
		log.Errorw("error", "error", newErr)
		return nil, newErr
	}

	info = &WriteInfo{
		TotalRecords: int64(totalRecords),
		MetaRecords:  int64(metaProcessed),
		DataRecords:  int64(dataProcessed),
		MetaErrors:   int64(metaErrors),
		DataErrors:   int64(dataErrors),
		StationID:    stationID,
	}

	return
}

func prettyTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.String()
}

type newRecordStatistics struct {
	start time.Time
	end   time.Time
}

func (s *newRecordStatistics) addTime(t time.Time) {
	if s.start.IsZero() || t.Before(s.start) {
		s.start = t
	}
	if s.end.IsZero() || t.After(s.end) {
		s.end = t
	}
}
