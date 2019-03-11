package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func (fc *FileConcatenator) WriteAllFiles(ctx context.Context) (string, error) {
	log := Logger(ctx).Sugar()

	log.Infow("Concatenating device files", "device_id", fc.DeviceID)

	files := []*data.DeviceFile{}
	if err := fc.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.file_id = ANY($1)) AND (s.device_id = $2) ORDER BY time ASC`,
		pq.StringArray(fc.TypeIDs), fc.DeviceID); err != nil {
		return "", err
	}

	temporaryFile, err := ioutil.TempFile("", fmt.Sprintf("%s.fkpb", fc.DeviceID))
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
			return "", fmt.Errorf("Error reading object: %v", err)
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

			return &record, nil
		})

		_, err = stream.ReadLengthPrefixedCollection(obj.Body, unmarshalFunc)
		if err != nil {
			return "", fmt.Errorf("Error reading collection: %v", err)
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

	fi, err := reading.Stat()
	if err != nil {
		return "", fmt.Errorf("Error getting FileInfo: %v", err)
	}

	uploader := s3manager.NewUploader(fc.Session)

	headers := IncomingHeaders{}

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

	jsonMetadata, err := json.Marshal(headers)
	if err != nil {
		return "", fmt.Errorf("JSON error: %v", err)
	}

	stream := data.DeviceStream{
		Time:     time.Now(),
		StreamID: fc.FileID,
		Firmware: "",
		DeviceID: fc.DeviceID,
		Size:     fi.Size(),
		FileID:   fc.FileTypeID,
		URL:      res.Location,
		Meta:     jsonMetadata,
	}

	if _, err := fc.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_stream (time, stream_id, firmware, device_id, size, file_id, url, meta)
		   VALUES (:time, :stream_id, :firmware, :device_id, :size, :file_id, :url, :meta)
		   `, stream); err != nil {
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
