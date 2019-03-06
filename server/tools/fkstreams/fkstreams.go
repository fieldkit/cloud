package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	"github.com/Devatoria/go-graylog"

	"github.com/conservify/sqlxcache"

	pbtools "github.com/conservify/protobuf-tools/tools"

	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"

	fk "github.com/fieldkit/cloud/server/api/client"
	fktesting "github.com/fieldkit/cloud/server/tools"
)

func createAwsSession() (s *session.Session, err error) {
	configs := []aws.Config{
		aws.Config{
			Region:                        aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewEnvCredentials(),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}

	for _, config := range configs {
		sessionOptions := session.Options{
			Profile: "fieldkit",
			Config:  config,
		}
		session, err := session.NewSessionWithOptions(sessionOptions)
		if err == nil {
			return session, nil
		}
	}

	return nil, fmt.Errorf("Error creating AWS session: %v", err)
}

type ObjectJob struct {
	Bucket string
	Key    string
}

type ObjectDetails struct {
	Bucket       string
	Key          string
	LastModified *time.Time
	Size         int64
	DeviceId     string
	Firmware     string
	Build        string
	FileId       string
	URL          string
	Meta         map[string]*string
}

func fixMeta(m map[string]*string) map[string]*string {
	newM := make(map[string]*string)
	newM["Fk-DeviceId"] = m["Fk-Deviceid"]
	newM["Fk-FileId"] = m["Fk-Fileid"]
	newM["Fk-Build"] = m["Fk-Build"]
	newM["Fk-FileName"] = m["Fk-Filename"]
	newM["Fk-Version"] = m["Fk-Version"]

	return newM
}

func worker(ctx context.Context, id int, svc *s3.S3, db *sqlxcache.DB, jobs <-chan ObjectJob, details chan<- ObjectDetails) {
	for job := range jobs {
		files := []*data.DeviceStream{}
		if err := db.SelectContext(ctx, &files, `SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.stream_id = $1)`, job.Key); err != nil {
			panic(fmt.Errorf("Error querying for DeviceFile: %v", err))
		}

		if len(files) != 1 {
			hoi := &s3.HeadObjectInput{
				Bucket: aws.String(job.Bucket),
				Key:    aws.String(job.Key),
			}

			obj, err := svc.HeadObject(hoi)
			if err != nil {
				log.Printf("Error getting object Head: %v", err)
				continue
			}

			meta := fixMeta(obj.Metadata)
			deviceId := meta["Fk-DeviceId"]
			firmware := meta["Fk-Version"]
			build := meta["Fk-Build"]
			fileId := meta["Fk-FileId"]

			if deviceId == nil || firmware == nil || build == nil || fileId == nil {
				log.Printf("Incomplete Metadata: %v", obj)
			} else {
				od := ObjectDetails{
					Bucket:       job.Bucket,
					Key:          job.Key,
					LastModified: obj.LastModified,
					Size:         *obj.ContentLength,
					Firmware:     *firmware,
					Build:        *build,
					DeviceId:     *deviceId,
					FileId:       *fileId,
					URL:          fmt.Sprintf("https://%s.s3.amazonaws.com/%s", job.Bucket, job.Key),
					Meta:         meta,
				}

				jsonMeta, err := json.Marshal(od.Meta)
				if err != nil {
					panic(err)
				}

				stream := data.DeviceStream{
					Time:     *od.LastModified,
					StreamID: od.Key,
					Firmware: od.Firmware,
					DeviceID: od.DeviceId,
					Size:     int64(od.Size),
					FileID:   od.FileId,
					URL:      od.URL,
					Meta:     jsonMeta,
				}

				log.Printf("Missing database entry: %v %v %v", job, od.LastModified, od.Size)

				if _, err := db.NamedExecContext(ctx, `
					  INSERT INTO fieldkit.device_stream (time, stream_id, firmware, device_id, size, file_id, url, meta)
					  VALUES (:time, :stream_id, :firmware, :device_id, :size, :file_id, :url, :meta)
					`, stream); err != nil {
					panic(err)
				}
			}
		}
	}
}

func process(job ObjectJob, svc *s3.S3) error {
	goi := &s3.GetObjectInput{
		Bucket: aws.String(job.Bucket),
		Key:    aws.String(job.Key),
	}

	obj, err := svc.GetObject(goi)
	if err != nil {
		return err
	}

	fmt.Println(obj)

	unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var record pb.DataRecord

		err := proto.Unmarshal(b, &record)
		if err != nil {
			return nil, err
		}

		fmt.Println(record)

		return &record, nil
	})

	if false {
		_, junk, err := pbtools.ReadLengthPrefixedCollectionIgnoringIncompleteBeginning(obj.Body, 4096, unmarshalFunc)
		if err != nil {
			return err
		}
		if junk > 0 {
			log.Printf("Malformed stream, ignored junk")
		}
	} else {
		_, err := stream.ReadLengthPrefixedCollection(obj.Body, unmarshalFunc)
		if err != nil {
			return err
		}
	}

	return nil
}

func listStreams(ctx context.Context, session *session.Session, db *sqlxcache.DB) error {
	svc := s3.New(session)

	loi := &s3.ListObjectsInput{
		Bucket:  aws.String("fk-streams"),
		MaxKeys: aws.Int64(100),
	}

	jobs := make(chan ObjectJob, 100)
	details := make(chan ObjectDetails, 100)

	log.Printf("Listing...")

	for w := 0; w < 10; w++ {
		go worker(ctx, w, svc, db, jobs, details)
	}

	total := 0

	err := svc.ListObjectsPages(loi, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		log.Printf("%v objects (%v)", len(page.Contents), total)

		for _, summary := range page.Contents {
			job := ObjectJob{
				Bucket: *loi.Bucket,
				Key:    *summary.Key,
			}

			jobs <- job
		}

		total += len(page.Contents)

		return true
	})
	if err != nil {
		return fmt.Errorf("Unable to list S3 bucket: %v", err)
	}

	log.Printf("%v total objects", total)

	close(jobs)

	return nil
}

func grayLogMessage() {
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	g, err := graylog.NewGraylog(graylog.Endpoint{
		Transport: graylog.TCP,
		Address:   "172.31.58.48",
		Port:      12201,
	})
	if err != nil {
		panic(err)
	}

	defer g.Close()

	err = g.Send(graylog.Message{
		Version:      "1.1",
		Host:         hn,
		ShortMessage: "test",
		FullMessage:  "test",
		Timestamp:    time.Now().Unix(),
		Level:        1,
		Extra: map[string]string{
			"tag":       "fkdev/devices",
			"stream":    "",
			"device_id": "",
		},
	})
	if err != nil {
		panic(err)
	}
}

type options struct {
	Scheme   string
	Host     string
	Username string
	Password string

	DeviceID string

	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	Sync        bool
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Scheme, "scheme", "http", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "fk instance hostname")
	flag.StringVar(&o.Username, "username", "demo-user", "username")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password")

	flag.StringVar(&o.DeviceID, "device-id", "", "device id")

	flag.BoolVar(&o.Sync, "sync", false, "sync")

	flag.Parse()

	if err := envconfig.Process("fieldkit", &o); err != nil {
		panic(err)
	}

	log.Printf("%v", o.PostgresURL)

	fkc, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		panic(err)
	}

	log.Printf("Authenticated as %s (%s)", o.Username, o.Host)

	if o.Sync {
		db, err := sqlxcache.Open("postgres", o.PostgresURL)
		if err != nil {
			panic(err)
		}

		session, err := createAwsSession()
		if err != nil {
			panic(err)
		}

		err = listStreams(ctx, session, db)
		if err != nil {
			panic(err)
		}

		return
	}

	if o.DeviceID != "" {
		page := 0
		total := 0

		for {
			res, err := fkc.ListDeviceLogFilesFiles(ctx, fk.ListDeviceDataFilesFilesPath(o.DeviceID), &page)
			if err != nil {
				panic(err)
			}

			files, err := fkc.DecodeDeviceFiles(res)
			if err != nil {
				panic(err)
			}

			if len(files.Files) == 0 {
				break
			}

			total += len(files.Files)

			for _, file := range files.Files {
				bucket, key, err := getBucketAndKey(file.URL)
				if err != nil {
					panic(err)
				}

				log.Printf("%v %v %v", o.DeviceID, bucket, key)
			}

			page += 1
		}

		log.Printf("Total Files: %v", total)

		return
	}

	if o.DeviceID == "" {
		res, err := fkc.ListDevicesFiles(ctx, fk.ListDevicesFilesPath())
		if err != nil {
			panic(err)
		}

		devices, err := fkc.DecodeDevices(res)
		if err != nil {
			panic(err)
		}

		for _, device := range devices.Devices {
			log.Printf("%v %v", device.DeviceID, device.LastFileTime)
		}

		log.Printf("Total Devices: %v", len(devices.Devices))

		return
	}
}

func getBucketAndKey(s3Url string) (bucket, key string, err error) {
	u, err := url.Parse(s3Url)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(u.Host, ".")

	return parts[0], u.Path[1:], nil
}
