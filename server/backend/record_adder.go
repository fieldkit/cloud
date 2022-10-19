package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
)

type WriteInfo struct {
	IngestionID   int64
	TotalRecords  int64
	DataRecords   int64
	MetaRecords   int64
	MetaErrors    int64
	DataErrors    int64
	FutureIgnores int64
	StationID     *int32
	DataStart     time.Time
	DataEnd       time.Time
}

type keyRecord struct {
	fileOffset   int64
	recordNumber int64
}

func (kr *keyRecord) adjustRecordNumber(fileOffset int64) int64 {
	return kr.recordNumber + (fileOffset - kr.fileOffset)
}

type RecordAdder struct {
	verbose             bool
	handler             RecordHandler
	db                  *sqlxcache.DB
	files               files.FileArchive
	metrics             *logging.Metrics
	stationRepository   *repositories.StationRepository
	provisionRepository *repositories.ProvisionRepository
	recordRepository    *repositories.RecordRepository
	statistics          *newRecordStatistics
	queue               []*ParsedRecord
	ingestion           *data.Ingestion
	provision           *data.Provision
	keyRecord           *keyRecord
	metaID              int64
}

type ParsedRecord struct {
	SignedRecord *pb.SignedRecord
	DataRecord   *pb.DataRecord
	Bytes        []byte
	FileOffset   int64
}

func NewRecordAdder(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, handler RecordHandler, verbose bool, saveData bool) (ra *RecordAdder) {
	return &RecordAdder{
		verbose:             verbose,
		db:                  db,
		files:               files,
		metrics:             metrics,
		handler:             handler,
		stationRepository:   repositories.NewStationRepository(db),
		provisionRepository: repositories.NewProvisionRepository(db),
		recordRepository:    repositories.NewRecordRepository(db, saveData),
		statistics:          &newRecordStatistics{},
		keyRecord:           nil,
		provision:           nil,
		ingestion:           nil,
		metaID:              0,
	}
}

func (ra *RecordAdder) tryParseSignedRecord(sr *pb.SignedRecord, dataRecord *pb.DataRecord) (bytes []byte, err error) {
	if err := proto.Unmarshal(sr.Data, dataRecord); err == nil {
		return sr.Data, nil
	}

	buffer := proto.NewBuffer(sr.Data)
	size, err := buffer.DecodeVarint()
	if err != nil {
		return nil, err
	}

	sizeSize := uint64(proto.SizeVarint(size))

	if size+sizeSize != uint64(len(sr.Data)) {
		return nil, fmt.Errorf("bad length prefix in signed record: %d", size)
	}

	if err := buffer.Unmarshal(dataRecord); err != nil {
		return nil, err
	}

	bytes = sr.Data[sizeSize:]

	return bytes, nil
}

func (ra *RecordAdder) tryFindStation(ctx context.Context, i *data.Ingestion) (*data.Station, error) {
	return ra.stationRepository.TryQueryStationByDeviceID(ctx, i.DeviceID)
}

func (ra *RecordAdder) findProvision(ctx context.Context, ingestion *data.Ingestion) error {
	if provision, err := ra.provisionRepository.QueryOrCreateProvision(ctx, ingestion.DeviceID, ingestion.GenerationID); err != nil {
		return err
	} else {
		ra.provision = provision
		ra.ingestion = ingestion
		return nil
	}
}

func (ra *RecordAdder) onSignedMeta(ctx context.Context, provision *data.Provision, ingestion *data.Ingestion, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	metaRecord, err := ra.recordRepository.AddSignedMetaRecord(ctx, ra.provision, ra.ingestion, signedRecord, rawRecord, bytes)
	if err != nil {
		return err
	}

	log.Infow("adder:on-meta-signed", "meta_record_id", metaRecord.ID)

	if err := ra.handler.OnMeta(ctx, ra.provision, rawRecord, metaRecord); err != nil {
		return err
	}

	ra.metaID = metaRecord.ID

	return nil
}

func (ra *RecordAdder) onMeta(ctx context.Context, provision *data.Provision, ingestion *data.Ingestion, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	metaRecord, err := ra.recordRepository.AddMetaRecord(ctx, ra.provision, ra.ingestion, recordNumber, rawRecord, bytes)
	if err != nil {
		return err
	}

	log.Infow("adder:on-meta", "meta_record_id", metaRecord.ID)

	if err := ra.handler.OnMeta(ctx, ra.provision, rawRecord, metaRecord); err != nil {
		return err
	}

	ra.metaID = metaRecord.ID

	return nil
}

func (ra *RecordAdder) onData(ctx context.Context, provision *data.Provision, ingestion *data.Ingestion, rawRecord *pb.DataRecord, bytes []byte) error {
	log := Logger(ctx).Sugar()

	dataRecord, metaRecord, err := ra.recordRepository.AddDataRecord(ctx, ra.provision, ra.ingestion, rawRecord, bytes)
	if err != nil {
		return err
	}

	// Something to remember, is that AddDataRecord is returning the persisted
	// meta record for this data record, found by looking for the record by
	// number. This is especially important when ingesting partial files.
	if ra.metaID != metaRecord.ID {
		log.Infow("adder:on-meta-unvisisted", "meta_record_id", metaRecord.ID)

		if err := ra.handler.OnMeta(ctx, ra.provision, rawRecord, metaRecord); err != nil {
			return err
		}

		ra.metaID = metaRecord.ID
	}

	var rawMeta pb.DataRecord
	if err := metaRecord.Unmarshal(&rawMeta); err != nil {
		return err
	}

	if err := ra.handler.OnData(ctx, ra.provision, rawRecord, &rawMeta, dataRecord, metaRecord); err != nil {
		return err
	}

	ra.statistics.addTime(dataRecord.Time)

	return nil
}

func (ra *RecordAdder) flushQueue(ctx context.Context, required bool) (fatal error) {
	log := Logger(ctx).Sugar()

	// Do we have queued meta records in need of a record number?
	if len(ra.queue) == 0 {
		return nil
	}

	// If we're missing a key record, we can't import this file.
	if ra.keyRecord == nil {
		if !required {
			// Right now the only way this can happen is if we have a partial
			// file with no data and just meta records.
			log.Warnw("adder:unrequired-flush no-key-record")
			return nil
		}

		return fmt.Errorf("adder:flush error: no key record")
	}

	// So we can safetly process the queued meta records using their
	// file offset for the record number.
	for _, queued := range ra.queue {
		if err := ra.onMeta(ctx, ra.provision, ra.ingestion, ra.keyRecord.adjustRecordNumber(queued.FileOffset), queued.DataRecord, queued.Bytes); err != nil {
			return err
		}
	}

	log.Infow("adder:flushed", "nrecords", len(ra.queue))

	ra.queue = make([]*ParsedRecord, 0)

	return nil
}

func (ra *RecordAdder) Handle(ctx context.Context, pr *ParsedRecord) (warning error, fatal error) {
	log := Logger(ctx).Sugar()

	if pr.SignedRecord != nil {
		if err := ra.onSignedMeta(ctx, ra.provision, ra.ingestion, pr.SignedRecord, pr.DataRecord, pr.Bytes); err != nil {
			return nil, err
		}
	} else if pr.DataRecord != nil {
		// Older versions of firmware were using a nanopb version that wouldn't
		// allow us to elide empty child messages, so we need to check for a
		// field.
		if pr.DataRecord.Metadata != nil && pr.DataRecord.Metadata.DeviceId != nil {
			// Check to see if this record has a record number. Some older
			// firmware was never setting this which causes some headaches. What
			// we do is queue this up, and hope that we come across a data
			// record with a record number to anchor things. Notice that
			// non-single file meta records will come in via SignedMetaRecords.
			record := int64(pr.DataRecord.Metadata.Record)
			if record == 0 || len(ra.queue) > 0 {
				log.Infow("adder:queue-meta")

				ra.queue = append(ra.queue, pr)
			} else {
				// I'm not sure if this is entirely necessary, just better safe than sorry.
				if ra.keyRecord != nil {
					if err := ra.onMeta(ctx, ra.provision, ra.ingestion, ra.keyRecord.adjustRecordNumber(record), pr.DataRecord, pr.Bytes); err != nil {
						return nil, err
					}
				} else {
					if err := ra.onMeta(ctx, ra.provision, ra.ingestion, record, pr.DataRecord, pr.Bytes); err != nil {
						return nil, err
					}
				}
			}
		} else {
			// Are we need of a key record and does this record have one? Notice
			// that data readings can't be the first record, so looking for
			// non-zero readings here is safe.
			if ra.keyRecord == nil && pr.DataRecord.Readings.Reading > 0 {
				ra.keyRecord = &keyRecord{
					fileOffset:   pr.FileOffset,
					recordNumber: int64(pr.DataRecord.Readings.Reading),
				}
				log.Infow("key-record", "file_offset", ra.keyRecord.fileOffset, "record_number", ra.keyRecord.recordNumber)
			}

			if err := ra.flushQueue(ctx, true); err != nil {
				return nil, err
			}

			if err := ra.onData(ctx, ra.provision, ra.ingestion, pr.DataRecord, pr.Bytes); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (ra *RecordAdder) WriteRecords(ctx context.Context, ingestion *data.Ingestion) (info *WriteInfo, err error) {
	log := Logger(ctx).Sugar().With("ingestion_id", ingestion.ID, "device_id", ingestion.DeviceID, "user_id", ingestion.UserID, "blocks", ingestion.Blocks)

	log.Infow("file", "file_url", ingestion.URL, "file_stamp", ingestion.Time, "file_id", ingestion.UploadID, "file_size", ingestion.Size, "type", ingestion.Type)

	station, err := ra.tryFindStation(ctx, ingestion)
	if err != nil {
		return nil, err
	}

	if err := ra.findProvision(ctx, ingestion); err != nil {
		return nil, err
	}

	opened, err := ra.files.OpenByURL(ctx, ingestion.URL)
	if err != nil {
		return nil, err
	}

	reader := opened.Body

	defer reader.Close()

	meta := false
	data := false

	totalRecords := 0
	dataProcessed := 0
	dataErrors := 0
	metaProcessed := 0
	metaErrors := 0
	records := 0
	recordNumber := int64(0)

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
					return nil, fmt.Errorf("adder:error parsing data record: %w", err)
				}
				unmarshalError = err
			} else {
				data = true
				warning, fatal := ra.Handle(ctx, &ParsedRecord{
					FileOffset: recordNumber,
					DataRecord: &dataRecord,
					Bytes:      b,
				})
				if fatal != nil {
					return nil, fatal
				}
				if warning == nil {
					dataProcessed += 1
					records += 1
				} else {
					dataErrors += 1
					if records > 0 {
						log.Infow("adder:processed", "record_run", records)
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
					return nil, fmt.Errorf("adder:error parsing signed record: %w", err)
				}
				unmarshalError = err
			} else {
				meta = true
				bytes, err := ra.tryParseSignedRecord(&signedRecord, &dataRecord)
				if err != nil {
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("adder:error parsing signed record: %w", err)
				}

				warning, fatal := ra.Handle(ctx, &ParsedRecord{
					FileOffset:   recordNumber,
					SignedRecord: &signedRecord,
					DataRecord:   &dataRecord,
					Bytes:        bytes,
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
						log.Infow("adder:processed", "record_run", records)
						records = 0
					}
				}
			}
		}

		// Only used with the new single file uploads.
		recordNumber += 1

		// Parsing either one failed, otherwise this would have been set. So return the error.
		if !meta && !data {
			return nil, unmarshalError
		}

		return nil, nil
	})

	if _, _, err = ReadLengthPrefixedCollection(ctx, MaximumDataRecordLength, reader, unmarshalFunc); err != nil {
		newErr := fmt.Errorf("adder:error: %w", err)
		log.Errorw("adder:error", "error", newErr)
		return nil, newErr
	}

	if err := ra.flushQueue(ctx, false); err != nil {
		return nil, err
	}

	withInfo := log.With("meta_processed", metaProcessed,
		"data_processed", dataProcessed,
		"meta_errors", metaErrors,
		"data_errors", dataErrors,
		"record_run", records,
		"future_ignores", ra.statistics.futureIgnores,
		"start_human", prettyTime(ra.statistics.start),
		"end_human", prettyTime(ra.statistics.end))

	var stationID *int32
	if station != nil {
		withInfo = withInfo.With("station_name", station.Name)
		stationID = &station.ID
	}

	withInfo.Infow("adder:processed")

	if err := ra.handler.OnDone(ctx); err != nil {
		return nil, fmt.Errorf("adder:done handler error: %w", err)
	}

	info = &WriteInfo{
		IngestionID:   ingestion.ID,
		TotalRecords:  int64(totalRecords),
		MetaRecords:   int64(metaProcessed),
		DataRecords:   int64(dataProcessed),
		MetaErrors:    int64(metaErrors),
		DataErrors:    int64(dataErrors),
		FutureIgnores: int64(ra.statistics.futureIgnores),
		StationID:     stationID,
		DataStart:     ra.statistics.start,
		DataEnd:       ra.statistics.end,
	}

	withInfo.Infow("adder:done")

	return
}

func prettyTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.String()
}

type newRecordStatistics struct {
	start         time.Time
	end           time.Time
	futureIgnores int
}

func (s *newRecordStatistics) addTime(t time.Time) {
	if t.Before(time.Now()) {
		if s.start.IsZero() || t.Before(s.start) {
			s.start = t
		}
		if s.end.IsZero() || t.After(s.end) {
			s.end = t
		}
	} else {
		s.futureIgnores++
	}
}
