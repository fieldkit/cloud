package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"

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

type MetaRecord struct {
	ID          int64          `db:"id"`
	IngestionID int64          `db:"ingestion_id"`
	Time        time.Time      `db:"time"`
	Number      int64          `db:"number"`
	Data        types.JSONText `db:"raw"`
}

func (d *DataRecord) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.Data = jsonData
	return nil
}

func (d *DataRecord) GetData() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

type DataRecord struct {
	ID          int64          `db:"id"`
	IngestionID int64          `db:"ingestion_id"`
	Time        time.Time      `db:"time"`
	Number      int64          `db:"number"`
	Data        types.JSONText `db:"raw"`
}

func (d *MetaRecord) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.Data = jsonData
	return nil
}

func (d *MetaRecord) GetData() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

func (ra *RecordAdder) Handle(ctx context.Context, i *data.Ingestion, pr *ParsedRecord) error {
	if pr.SignedRecord != nil {
		metaRecord := MetaRecord{
			IngestionID: i.ID,
			Time:        i.Time,
			Number:      int64(pr.SignedRecord.Record),
		}
		if err := metaRecord.SetData(pr.DataRecord); err != nil {
			return err
		}
		return ra.Database.NamedGetContext(ctx, metaRecord, `INSERT INTO fieldkit.meta_record (ingestion_id, time, number, raw) VALUES (:ingestion_id, :time, :number, :raw) RETURNING *`, metaRecord)
	} else {
		dataRecord := DataRecord{
			IngestionID: i.ID,
			Time:        i.Time,
			Number:      int64(pr.DataRecord.Readings.Reading),
		}
		if err := dataRecord.SetData(pr.DataRecord); err != nil {
			return err
		}
		return ra.Database.NamedGetContext(ctx, dataRecord, `INSERT INTO fieldkit.data_record (ingestion_id, time, number, raw) VALUES (:ingestion_id, :time, :number, :raw) RETURNING *`, dataRecord)
	}
	return nil
}

func (ra *RecordAdder) WriteRecords(ctx context.Context, i *data.Ingestion) error {
	log := Logger(ctx).Sugar()

	svc := s3.New(ra.Session)

	object, err := GetBucketAndKey(i.URL)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %v", err)
	}

	log.Infow("file", "file_url", i.URL, "file_stamp", i.Time, "stream_id", i.UploadID, "file_size", i.Size, "blocks", i.Blocks, "device_id", i.DeviceID, "user_id", i.UserID)

	goi := &s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	}

	obj, err := svc.GetObject(goi)
	if err != nil {
		return fmt.Errorf("Error reading object %v: %v", object.Key, err)
	}

	defer obj.Body.Close()

	meta := false
	data := false

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
				err = ra.Handle(ctx, i, &ParsedRecord{DataRecord: &dataRecord})
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
				err = ra.Handle(ctx, i, &ParsedRecord{
					SignedRecord: &signedRecord,
					DataRecord:   &dataRecord,
				})

			}
		}

		// Parsing either one failed, otherwise this would have been set. So return the error.
		if !meta && !data {
			return nil, unmarshalError
		}

		return nil, nil
	})

	_, _, err = ReadLengthPrefixedCollection(ctx, MaximumDataRecordLength, obj.Body, unmarshalFunc)
	if err != nil {
		newErr := fmt.Errorf("Error reading collection: %v", err)
		log.Errorw("Error", "error", newErr)
		return newErr
	}

	return nil
}
