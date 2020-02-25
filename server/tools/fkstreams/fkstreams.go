package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common"
	data "github.com/fieldkit/cloud/server/data"
	fktesting "github.com/fieldkit/cloud/server/tools"
)

type options struct {
	Scheme   string
	Host     string
	Username string
	Password string

	DeviceID string

	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	Sync        bool
}

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
	common.BucketAndKey
}

type ObjectDetails struct {
	Bucket       string
	Key          string
	LastModified *time.Time
	Size         int64
	DeviceId     string
	URL          string
	Meta         map[string]*string
}

func worker(ctx context.Context, id int, svc *s3.S3, db *sqlxcache.DB, jobs <-chan ObjectJob, details chan<- ObjectDetails) {
	for job := range jobs {
		files := []*data.DeviceFile{}
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

			meta := common.SanitizeMeta(obj.Metadata)
			deviceId := meta[common.FkDeviceIdHeaderName]

			if deviceId == nil {
				log.Printf("Incomplete Metadata: %v", obj)
			} else {
				od := ObjectDetails{
					Bucket:       job.Bucket,
					Key:          job.Key,
					LastModified: obj.LastModified,
					Size:         *obj.ContentLength,
					DeviceId:     *deviceId,
					URL:          fmt.Sprintf("https://%s.s3.amazonaws.com/%s", job.Bucket, job.Key),
					Meta:         meta,
				}

				jsonMeta, err := json.Marshal(od.Meta)
				if err != nil {
					panic(err)
				}

				stream := data.DeviceFile{
					Time:     *od.LastModified,
					StreamID: od.Key,
					DeviceID: od.DeviceId,
					Size:     int64(od.Size),
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

func listStreams(ctx context.Context, o *options, session *session.Session, db *sqlxcache.DB) error {
	svc := s3.New(session)

	loi := &s3.ListObjectsV2Input{
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

	err := svc.ListObjectsV2Pages(loi, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		log.Printf("%v objects (%v)", len(page.Contents), total)

		for _, summary := range page.Contents {
			job := ObjectJob{
				BucketAndKey: common.BucketAndKey{
					Bucket: *loi.Bucket,
					Key:    *summary.Key,
				},
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

	_, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		panic(err)
	}

	log.Printf("Authenticated as %s (%s)", o.Username, o.Host)

	session, err := createAwsSession()
	if err != nil {
		panic(err)
	}

	if o.Sync {
		db, err := sqlxcache.Open("postgres", o.PostgresURL)
		if err != nil {
			panic(err)
		}

		err = listStreams(ctx, &o, session, db)
		if err != nil {
			panic(err)
		}

		return
	}
}
