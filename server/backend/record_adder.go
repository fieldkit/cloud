package backend

import (
	"context"
	"fmt"
	"math"
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
	database   *sqlxcache.DB
	files      files.FileArchive
	metrics    *logging.Metrics
	statistics *newRecordStatistics
}

type ParsedRecord struct {
	SignedRecord *pb.SignedRecord
	DataRecord   *pb.DataRecord
}

func NewRecordAdder(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, verbose bool) (ra *RecordAdder) {
	return &RecordAdder{
		verbose:    verbose,
		database:   db,
		files:      files,
		metrics:    metrics,
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
	if err := ra.database.SelectContext(ctx, &provisions, `SELECT p.* FROM fieldkit.provision AS p WHERE p.device_id = $1 AND p.generation = $2`, i.DeviceID, i.GenerationID); err != nil {
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
		    DO UPDATE SET updated = NOW() RETURNING id`, provision); err != nil {
		return nil, err
	}

	return provision, nil
}

func (ra *RecordAdder) findMeta(ctx context.Context, provisionId, number int64) (*data.MetaRecord, error) {
	records := []*data.MetaRecord{}
	if err := ra.database.SelectContext(ctx, &records, `SELECT r.* FROM fieldkit.meta_record AS r WHERE r.provision_id = $1 AND r.number = $2`, provisionId, number); err != nil {
		return nil, err
	}

	if len(records) != 1 {
		return nil, nil
	}

	return records[0], nil
}

func (ra *RecordAdder) findLocation(dataRecord *pb.DataRecord) (l *data.Location, err error) {
	if dataRecord.Readings == nil || dataRecord.Readings.Location == nil {
		return nil, err
	}
	location := dataRecord.Readings.Location
	lat := float64(location.Latitude)
	lon := float64(location.Longitude)
	altitude := float64(location.Altitude)

	if lat > 90 || lat < -90 || lon > 180 || lon < -180 {
		return nil, err
	}

	if lat == 0 && lon == 0 {
		return nil, err
	}

	l = data.NewLocation([]float64{lon, lat, altitude})
	return
}

func hasNaNs(dr *pb.DataRecord) bool {
	for _, sg := range dr.Readings.SensorGroups {
		for _, sr := range sg.Readings {
			if math.IsNaN(float64(sr.Value)) {
				return true
			}

		}
	}
	return false
}

func prepareForMarshalToJson(dr *pb.DataRecord) *pb.DataRecord {
	if !hasNaNs(dr) {
		return dr
	}
	for _, sg := range dr.Readings.SensorGroups {
		newReadings := make([]*pb.SensorAndValue, len(sg.Readings))
		for _, sr := range sg.Readings {
			if !math.IsNaN(float64(sr.Value)) {
				newReadings = append(newReadings, sr)
			}

		}
		sg.Readings = newReadings
	}
	return dr
}

func (ra *RecordAdder) Handle(ctx context.Context, i *data.Ingestion, pr *ParsedRecord) (warning error, fatal error) {
	log := Logger(ctx).Sugar()
	verboseLog := logging.OnlyLogIf(log, ra.verbose)

	provision, err := ra.findProvision(ctx, i)
	if err != nil {
		return nil, err
	}

	if pr.SignedRecord != nil {
		metaTime := i.Time
		if pr.SignedRecord.Time > 0 {
			metaTime = time.Unix(int64(pr.SignedRecord.Time), 0)
		}

		metaRecord := data.MetaRecord{
			ProvisionID: provision.ID,
			Time:        metaTime,
			Number:      int64(pr.SignedRecord.Record),
		}

		if err := metaRecord.SetData(pr.DataRecord); err != nil {
			return nil, fmt.Errorf("error setting meta json: %v", err)
		}

		if err := ra.database.NamedGetContext(ctx, &metaRecord, `
		    INSERT INTO fieldkit.meta_record (provision_id, time, number, raw) VALUES (:provision_id, :time, :number, :raw)
		    ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw
                    RETURNING *`, metaRecord); err != nil {
			return nil, err
		}
	} else if pr.DataRecord != nil {
		if pr.DataRecord.Readings == nil {
			verboseLog.Infow("data reading missing readings", "record", pr.DataRecord)
			ra.metrics.DataErrorsUnknown()
			return fmt.Errorf("weird record"), nil
		}

		location, err := ra.findLocation(pr.DataRecord)
		if err != nil {
			return nil, err
		}

		meta, err := ra.findMeta(ctx, provision.ID, int64(pr.DataRecord.Readings.Meta))
		if err != nil {
			return nil, err
		}
		if meta == nil {
			verboseLog.Errorw("error finding meta record", "provision_id", provision.ID, "meta_record_number", pr.DataRecord.Readings.Meta)
			ra.metrics.DataErrorsMissingMeta()
			return fmt.Errorf("error finding meta record"), nil
		}

		dataTime := i.Time
		if pr.DataRecord.Readings.Time > 0 {
			dataTime = time.Unix(int64(pr.DataRecord.Readings.Time), 0)
			ra.statistics.addTime(dataTime)
		}

		dataRecord := data.DataRecord{
			ProvisionID: provision.ID,
			Time:        dataTime,
			Number:      int64(pr.DataRecord.Readings.Reading),
			MetaID:      meta.ID,
			Location:    location,
		}

		if err := dataRecord.SetData(prepareForMarshalToJson(pr.DataRecord)); err != nil {
			return nil, fmt.Errorf("error setting data json: %v", err)
		}

		if err := ra.database.NamedGetContext(ctx, &dataRecord, `
			    INSERT INTO fieldkit.data_record (provision_id, time, number, raw, meta, location)
			    VALUES (:provision_id, :time, :number, :raw, :meta, ST_SetSRID(ST_GeomFromText(:location), 4326))
			    ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw, location = EXCLUDED.location
			    RETURNING id`, dataRecord); err != nil {
			return nil, err
		}
	}

	return
}

type WriteInfo struct {
	TotalRecords int64
	MetaErrors   int64
	DataErrors   int64
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

	if station != nil {
		withInfo = withInfo.With("station_name", station.Name)
	}

	withInfo.Infow("processed")

	if err != nil {
		newErr := fmt.Errorf("error reading collection: %v", err)
		log.Errorw("error", "error", newErr)
		return nil, newErr
	}

	info = &WriteInfo{
		TotalRecords: int64(totalRecords),
		MetaErrors:   int64(metaErrors),
		DataErrors:   int64(dataErrors),
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
