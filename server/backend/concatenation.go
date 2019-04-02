package backend

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type FileConcatenator struct {
	Session  *session.Session
	Database *sqlxcache.DB
	Backend  *Backend

	FileID     string
	FileTypeID string
	DeviceID   string
	TypeIDs    []string
}

func NewFileConcatenator(s *session.Session, db *sqlxcache.DB, backend *Backend) (fc *FileConcatenator, err error) {
	fc = &FileConcatenator{
		Session:  s,
		Database: db,
		Backend:  backend,
	}

	return
}

func (fc *FileConcatenator) HandleLocation(ctx context.Context, file *data.DeviceFile, ts time.Time, coordinates []float64) error {
	dsl := data.DeviceStreamLocation{
		DeviceID:  file.DeviceID,
		Timestamp: ts,
		Location:  data.NewLocation(coordinates),
	}

	_, err := fc.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_stream_location (device_id, timestamp, location)
		   VALUES (:device_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326))
		   ON CONFLICT DO NOTHING
		   `, dsl)
	if err != nil {
		return err
	}

	return nil
}

var LocationRegex = regexp.MustCompile("Loc\\((\\S+), (\\S+), (\\S+)\\)")

func (fc *FileConcatenator) ProcessRecord(ctx context.Context, file *data.DeviceFile, record *pb.DataRecord) error {
	log := Logger(ctx).Sugar()

	if record.LoggedReading != nil && record.LoggedReading.Location != nil {
		location := record.LoggedReading.Location
		coordinates := []float64{float64(location.Longitude), float64(location.Latitude), float64(location.Altitude)}
		ts := time.Unix(int64(location.Time), 0)

		err := fc.HandleLocation(ctx, file, ts, coordinates)
		if err != nil {
			return err
		}
	}

	if record.Log != nil {
		m := LocationRegex.FindAllStringSubmatch(record.Log.Message, -1)
		if len(m) > 0 {
			coordinates := make([]float64, 3)
			valid := true
			for i := 0; i < 3; i += 1 {
				v, err := strconv.ParseFloat(m[0][i+1], 64)
				if err != nil {
					panic(err)
				}

				if v >= 200 {
					valid = false
				}
				coordinates[i] = v
			}

			if valid {
				ts := time.Unix(int64(record.Log.Time), 0)

				err := fc.HandleLocation(ctx, file, ts, coordinates)
				if err != nil {
					return err
				}
			}
		}
	}

	_ = log

	return nil
}

func (fc *FileConcatenator) WriteAllFiles(ctx context.Context) (string, error) {
	log := Logger(ctx).Sugar()

	log.Infow("Concatenating device files", "device_id", fc.DeviceID)

	files := []*data.DeviceFile{}
	if err := fc.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.file_id = ANY($1)) AND (s.device_id = $2) ORDER BY time ASC`,
		pq.StringArray(fc.TypeIDs), fc.DeviceID); err != nil {
		return "", err
	}

	temporaryFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.fkpb", fc.DeviceID))
	if err != nil {
		return "", fmt.Errorf("Error opening temporary file: %v", err)
	}

	defer temporaryFile.Close()

	svc := s3.New(fc.Session)

	for _, file := range files {
		object, err := GetBucketAndKey(file.URL)
		if err != nil {
			return "", fmt.Errorf("Error parsing URL: %v", err)
		}

		log.Infow("File", "file_url", file.URL, "file_stamp", file.Time, "stream_id", file.StreamID)

		goi := &s3.GetObjectInput{
			Bucket: aws.String(object.Bucket),
			Key:    aws.String(object.Key),
		}

		obj, err := svc.GetObject(goi)
		if err != nil {
			return "", fmt.Errorf("Error reading object %v: %v", object.Key, err)
		}

		unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
			var record pb.DataRecord

			err := proto.Unmarshal(b, &record)
			if err != nil {
				return nil, err
			}

			buf := proto.NewBuffer(nil)
			buf.EncodeRawBytes(b)

			_, err = temporaryFile.Write(buf.Bytes())
			if err != nil {
				return nil, err
			}

			fc.ProcessRecord(ctx, file, &record)

			return nil, nil
		})

		_, err = stream.ReadLengthPrefixedCollection(obj.Body, unmarshalFunc)
		if err != nil {
			newErr := fmt.Errorf("Error reading collection: %v", err)
			log.Errorw("Error", "error", newErr)
			if false {
				return "", newErr
			}
		}
	}

	return temporaryFile.Name(), nil
}

func (fc *FileConcatenator) Upload(ctx context.Context, path string) (string, error) {
	reading, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Error reading collection: %v", err)
	}

	defer reading.Close()

	uploader := s3manager.NewUploader(fc.Session)

	metadata := make(map[string]*string)
	metadata[FkDeviceIdHeaderName] = aws.String(fc.DeviceID)
	metadata[FkFileIdHeaderName] = aws.String(fc.FileTypeID)

	res, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		Bucket:      aws.String("fk-streams"),
		Key:         aws.String(fc.FileID),
		Body:        reading,
		ContentType: aws.String(FkDataBinaryContentType),
		Metadata:    metadata,
		Tagging:     nil,
	})
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (fc *FileConcatenator) Concatenate(ctx context.Context) {
	log := Logger(ctx).Sugar()

	name, err := fc.WriteAllFiles(ctx)
	if err != nil {
		log.Errorw("Error", "error", fmt.Errorf("Error writing files: %v", err))
		return
	}

	location, err := fc.Upload(ctx, name)
	if err != nil {
		log.Errorw("Error", "error", fmt.Errorf("Error uploading file: %v", err))
		return
	}

	log.Infow("Uploaded", "file_url", location)
}
