package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	"github.com/Devatoria/go-graylog"

	"github.com/Conservify/sqlxcache"

	pbtools "github.com/Conservify/protobuf-tools/tools"

	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"
)

type options struct {
}

func createAwsSession() (s *session.Session, err error) {
	configs := []aws.Config{
		aws.Config{
			Region: aws.String("us-east-1"),
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

func worker(id int, svc *s3.S3, jobs <-chan ObjectJob, details chan<- ObjectDetails) {
	for job := range jobs {
		hoi := &s3.HeadObjectInput{
			Bucket: aws.String(job.Bucket),
			Key:    aws.String(job.Key),
		}

		obj, err := svc.HeadObject(hoi)
		if err != nil {
			log.Printf("%v", err)
		}

		meta := fixMeta(obj.Metadata)

		deviceId := meta["Fk-DeviceId"]
		firmware := meta["Fk-Version"]
		build := meta["Fk-Build"]
		fileId := meta["Fk-FileId"]

		if deviceId == nil || firmware == nil || build == nil || fileId == nil {
			log.Printf("Incomplete metadata: %v", obj)
		} else {
			details <- ObjectDetails{
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
		log.Printf("%v", err)
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

func writeFile(details <-chan ObjectDetails) {
	csv, err := os.OpenFile("objects.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer csv.Close()

	db, err := sqlxcache.Open("postgres", os.Getenv("FIELDKIT_POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	_ = db

	for od := range details {
		fields := []string{
			od.Bucket,
			od.Key,
			od.LastModified.String(),
			fmt.Sprintf("%v", od.Size),
			od.Firmware,
			od.Build,
			od.DeviceId,
			od.FileId,
		}

		metaData, err := json.Marshal(od.Meta)
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
			Meta:     metaData,
		}

		log.Printf("%v\n", stream)

		_ = stream

		if false {
			if _, err := db.NamedExecContext(context.TODO(), `
				   INSERT INTO fieldkit.device_stream (time, stream_id, firmware, device_id, size, file_id, url, meta)
				   VALUES (:time, :stream_id, :firmware, :device_id, :size, :file_id, :url, :meta)
				   `, stream); err != nil {
				log.Printf("%v", err)
			}
		}

		csv.WriteString(strings.Join(fields, ",") + "\n")
	}
}

func listStreams(ctx context.Context) error {
	log.Printf("Creating session...")

	session, err := createAwsSession()
	if err != nil {
		return err
	}

	svc := s3.New(session)

	loi := &s3.ListObjectsInput{
		Bucket:  aws.String("fk-streams"),
		MaxKeys: aws.Int64(100),
	}

	jobs := make(chan ObjectJob, 100)
	details := make(chan ObjectDetails, 100)

	go writeFile(details)

	if true {
		log.Printf("Listing...")

		for w := 0; w < 10; w++ {
			go worker(w, svc, jobs, details)
		}

		total := 0

		err = svc.ListObjectsPages(loi, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			log.Printf("%v objects", len(page.Contents))

			// fmt.Println(page.Contents)

			for _, summary := range page.Contents {
				job := ObjectJob{
					Bucket: *loi.Bucket,
					Key:    *summary.Key,
				}

				// process(job, svc)

				jobs <- job
			}

			total += len(page.Contents)

			return true
		})
		if err != nil {
			return err
		}

		log.Printf("%v total objects", total)
	}

	close(jobs)

	return nil
}

func main() {
	ctx := context.TODO()

	flag.Parse()

	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	if true {
		err := listStreams(ctx)
		if err != nil {
			panic(err)
		}
	} else {
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

}
