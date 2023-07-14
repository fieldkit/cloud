package backend

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
)

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
	ingestion           *data.Ingestion
	provision           *data.Provision
	metaID              int64
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
		provision:           nil,
		ingestion:           nil,
		metaID:              0,
	}
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

func (ra *RecordAdder) OnSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
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

func (ra *RecordAdder) OnMeta(ctx context.Context, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
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

func (ra *RecordAdder) OnData(ctx context.Context, rawRecord *pb.DataRecord, rawMetaUnused *pb.DataRecord, bytes []byte) error {
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

func (ra *RecordAdder) OnDone(ctx context.Context) (err error) {
	return ra.handler.OnDone(ctx)
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

	fkbWalker := NewFkbWalker(ra.files, ra.metrics, ra, ra.verbose)

	walkInfo, err := fkbWalker.WalkUrl(ctx, ingestion.URL)
	if err != nil {
		return nil, err
	}

	withInfo := log.With(
		"future_ignores", ra.statistics.futureIgnores,
		"start_human", prettyTime(ra.statistics.start),
		"end_human", prettyTime(ra.statistics.end))

	withInfo.Info("adder:processed")

	info = &WriteInfo{
		Walk:          walkInfo,
		IngestionID:   ingestion.ID,
		StationID:     nil,
		FutureIgnores: int64(ra.statistics.futureIgnores),
		DataStart:     ra.statistics.start,
		DataEnd:       ra.statistics.end,
	}

	if station != nil {
		info.StationID = &station.ID
	}

	log.Infow("adder:done")

	return
}

type WriteInfo struct {
	IngestionID   int64
	StationID     *int32
	Walk          *WalkRecordsInfo
	FutureIgnores int64
	DataStart     time.Time
	DataEnd       time.Time
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
