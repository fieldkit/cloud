package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type RecordAdder struct {
	Session  *session.Session
	Database *sqlxcache.DB
}

func NewRecordAdder(s *session.Session, db *sqlxcache.DB) (ra *RecordAdder) {
	return &RecordAdder{
		Session:  s,
		Database: db,
	}
}

type ParsedRecord struct {
	SignedRecord *pb.SignedRecord
	DataRecord   *pb.DataRecord
}

func (ra *RecordAdder) TryParseSignedRecord(sr *pb.SignedRecord, dataRecord *pb.DataRecord) (err error) {
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

func (ra *RecordAdder) findProvision(ctx context.Context, i *data.Ingestion) (*data.Provision, error) {
	// TODO Verify we have a valid generation.

	provisions := []*data.Provision{}
	if err := ra.Database.SelectContext(ctx, &provisions, `SELECT p.* FROM fieldkit.provision AS p WHERE p.device_id = $1 AND p.generation = $2`, i.DeviceID, i.Generation); err != nil {
		return nil, err
	}

	if len(provisions) == 1 {
		return provisions[0], nil
	}

	provision := &data.Provision{
		Created:    time.Now(),
		Updated:    time.Now(),
		DeviceID:   i.DeviceID,
		Generation: i.Generation,
	}

	if err := ra.Database.NamedGetContext(ctx, &provision.ID, `
		    INSERT INTO fieldkit.provision (device_id, generation, created, updated)
		    VALUES (:device_id, :generation, :created, :updated) ON CONFLICT (device_id, generation)
		    DO UPDATE SET updated = NOW() RETURNING id`, provision); err != nil {
		return nil, err
	}

	return provision, nil
}

func (ra *RecordAdder) findMeta(ctx context.Context, provisionId, number int64) (*data.MetaRecord, error) {
	records := []*data.MetaRecord{}
	if err := ra.Database.SelectContext(ctx, &records, `SELECT r.* FROM fieldkit.meta_record AS r WHERE r.provision_id = $1 AND r.number = $2`, provisionId, number); err != nil {
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

func (ra *RecordAdder) Handle(ctx context.Context, i *data.Ingestion, pr *ParsedRecord) (warning error, fatal error) {
	log := Logger(ctx).Sugar()

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
			return nil, err
		}

		if err := ra.Database.NamedGetContext(ctx, &metaRecord, `
		    INSERT INTO fieldkit.meta_record (provision_id, time, number, raw) VALUES (:provision_id, :time, :number, :raw)
		    ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw
                    RETURNING *`, metaRecord); err != nil {
			return nil, err
		}
	} else if pr.DataRecord != nil {
		if pr.DataRecord.Readings != nil {
			location, err := ra.findLocation(pr.DataRecord)
			if err != nil {
				return nil, err
			}

			meta, err := ra.findMeta(ctx, provision.ID, int64(pr.DataRecord.Readings.Meta))
			if err != nil {
				return nil, err
			}
			if meta == nil {
				log.Errorw("error finding meta record", "provision_id", provision.ID, "meta_record_number", pr.DataRecord.Readings.Meta, "data_record", pr.DataRecord, "error", err)
				return fmt.Errorf("error finding meta record"), nil
			}

			dataTime := i.Time
			if pr.DataRecord.Readings.Time > 0 {
				dataTime = time.Unix(int64(pr.DataRecord.Readings.Time), 0)
			}

			dataRecord := data.DataRecord{
				ProvisionID: provision.ID,
				Time:        dataTime,
				Number:      int64(pr.DataRecord.Readings.Reading),
				Meta:        meta.ID,
				Location:    location,
			}

			if err := dataRecord.SetData(pr.DataRecord); err != nil {
				return nil, err
			}

			if err := ra.Database.NamedGetContext(ctx, &dataRecord, `
			    INSERT INTO fieldkit.data_record (provision_id, time, number, raw, meta, location)
			    VALUES (:provision_id, :time, :number, :raw, :meta, ST_SetSRID(ST_GeomFromText(:location), 4326))
			    ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw, location = EXCLUDED.location
			    RETURNING id`, dataRecord); err != nil {
				return nil, err
			}
		} else {
			log.Infow("weird", "record", pr.DataRecord)
			return fmt.Errorf("weird record"), nil
		}
	}

	return nil, nil
}

func (ra *RecordAdder) WriteRecords(ctx context.Context, i *data.Ingestion) error {
	log := Logger(ctx).Sugar()

	svc := s3.New(ra.Session)

	object, err := GetBucketAndKey(i.URL)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	log.Infow("file", "file_url", i.URL, "file_stamp", i.Time, "stream_id", i.UploadID, "file_size", i.Size, "blocks", i.Blocks, "device_id", i.DeviceID, "user_id", i.UserID, "type", i.Type)

	goi := &s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	}

	obj, err := svc.GetObject(goi)
	if err != nil {
		return fmt.Errorf("error reading object %v: %v", object.Key, err)
	}

	defer obj.Body.Close()

	meta := false
	data := false

	data_processed := 0
	data_errors := 0
	meta_processed := 0
	meta_errors := 0
	records := 0

	unmarshalFunc := UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var unmarshalError error
		var dataRecord pb.DataRecord
		var signedRecord pb.SignedRecord

		if data || (!data && !meta) {
			err := proto.Unmarshal(b, &dataRecord)
			if err != nil {
				if data { // If we expected this record, return the error
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
					data_processed += 1
					records += 1
				} else {
					data_errors += 1
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
					return nil, fmt.Errorf("error parsing signed record: %v", err)
				}
				unmarshalError = err
			} else {
				meta = true
				err := ra.TryParseSignedRecord(&signedRecord, &dataRecord)
				if err != nil {
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
					meta_processed += 1
					records += 1
				} else {
					meta_errors += 1
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

	_, _, err = ReadLengthPrefixedCollection(ctx, MaximumDataRecordLength, obj.Body, unmarshalFunc)

	log.Infow("processing done", "meta_processed", meta_processed, "data_processed", data_processed, "meta_errors", meta_errors, "data_errors", data_errors, "record_run", records)

	if err != nil {
		newErr := fmt.Errorf("error reading collection: %v", err)
		log.Errorw("error", "error", newErr)
		return newErr
	}

	return nil
}
