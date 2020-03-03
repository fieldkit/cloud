package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

type ObjectJob struct {
	common.BucketAndKey
}

func worker(ctx context.Context, id int, svc *s3.S3, db *sqlxcache.DB, jobs <-chan ObjectJob) {
	for job := range jobs {
		hoi := &s3.HeadObjectInput{
			Bucket: aws.String(job.Bucket),
			Key:    aws.String(job.Key),
		}

		obj, err := svc.HeadObject(hoi)
		if err != nil {
			log.Printf("error getting object: %v", err)
			continue
		}

		meta := common.SanitizeMeta(obj.Metadata)

		ingestions := []*data.Ingestion{}
		if err := db.SelectContext(ctx, &ingestions, `SELECT * FROM fieldkit.ingestion WHERE (file_id = $1)`, job.Key); err != nil {
			panic(fmt.Errorf("error: %v", err))
		}

		if len(ingestions) == 0 {
			log.Printf("deleting %v %v", job.Key, meta)
		}
	}
}

func listStreams(ctx context.Context, svc *s3.S3, bucketName string, jobs chan ObjectJob) error {
	loi := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(100),
	}

	log.Printf("listing")

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
		return fmt.Errorf("unable to list bucket: %v", err)
	}

	return nil
}

type options struct {
	PostgresURL string `split_words:"true"`
	BucketName  string
	Sync        bool
	Force       bool
}

func main() {
	ctx := context.Background()

	o := options{}

	flag.StringVar(&o.BucketName, "bucket-name", "", "bucket name")
	flag.BoolVar(&o.Sync, "sync", false, "sync")
	flag.BoolVar(&o.Force, "force", false, "force")

	flag.Parse()

	if err := envconfig.Process("FIELDKIT", &o); err != nil {
		panic(err)
	}

	if o.Sync {
		db, err := sqlxcache.Open("postgres", o.PostgresURL)
		if err != nil {
			panic(err)
		}

		session, err := createAwsSession()
		if err != nil {
			panic(err)
		}

		svc := s3.New(session)

		jobs := make(chan ObjectJob, 100)

		for w := 0; w < 10; w++ {
			go worker(ctx, w, svc, db, jobs)
		}

		err = listStreams(ctx, svc, o.BucketName, jobs)
		if err != nil {
			panic(err)
		}

		close(jobs)

		return
	}
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

	return nil, fmt.Errorf("aws error: %v", err)
}
