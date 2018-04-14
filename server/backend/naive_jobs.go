package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type PgJobQueue struct {
	db       *sqlxcache.DB
	name     string
	listener *pq.Listener
}

func NewPqJobQueue(url string, name string) (*PgJobQueue, error) {
	onProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("JobQueue: %v", err)
		}
	}

	db, err := sqlxcache.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	listener := pq.NewListener(url, 10*time.Second, time.Minute, onProblem)

	return &PgJobQueue{
		name:     name,
		db:       db,
		listener: listener,
	}, nil

}

func (jq *PgJobQueue) Start() error {
	err := jq.listener.Listen(jq.name)
	if err != nil {
		return err
	}

	log.Printf("Start monitoring %s...", jq.name)

	go jq.waitForNotification()

	return nil
}

func (jq *PgJobQueue) Publish(ctx context.Context, message interface{}) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = jq.db.ExecContext(ctx, `SELECT pg_notify($1, $2)`, jq.name, bytes)
	return err
}

func (jq *PgJobQueue) waitForNotification() {
	log.Printf("Waiting!")
	for {
		select {
		case n := <-jq.listener.Notify:
			log.Printf("Received data from channel [%v]:", n.Channel)
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				log.Printf("Error processing JSON: %v", err)
				break
			}
			log.Printf(string(prettyJSON.Bytes()))
			break
		case <-time.After(90 * time.Second):
			log.Printf("Received no events for 90 seconds, checking connection")
			go func() {
				jq.listener.Ping()
			}()
			break
		}
	}
}

type NaiveBackgroundJobs struct {
	be            *Backend
	sourceChanges chan ingestion.SourceChange
}

func NewNaiveBackgroundJobs(be *Backend) *NaiveBackgroundJobs {
	return &NaiveBackgroundJobs{
		be:            be,
		sourceChanges: make(chan ingestion.SourceChange, 100),
	}
}

func (j *NaiveBackgroundJobs) SourceChanged(sourceChange ingestion.SourceChange) {
	j.sourceChanges <- sourceChange
}

func (j *NaiveBackgroundJobs) Start() error {
	delayed := make(chan ingestion.SourceChange, 100)

	go func() {
		log.Printf("Started background jobs...")

		tick := time.Tick(1000 * time.Millisecond)

		buffer := make(map[int64]ingestion.SourceChange)

		for {
			select {
			case <-tick:
				for _, value := range buffer {
					delayed <- value
				}
				buffer = make(map[int64]ingestion.SourceChange)
			case c := <-j.sourceChanges:
				buffer[c.SourceID] = c
			}

		}
	}()

	go func() {
		previous := ingestion.SourceChange{}

		for {
			select {
			case c := <-delayed:
				if previous.SourceID == c.SourceID {
					if previous.StartedAt.After(c.QueuedAt) {
						log.Printf("Skipping %v", c)
						continue
					}
				}

				c.StartedAt = time.Now()
				log.Printf("Processing %v...", c.SourceID)
				generator := NewPregenerator(j.be)
				ctx := context.Background()
				err := generator.Pregenerate(ctx, c.SourceID)
				if err != nil {
					log.Printf("Error: %v", err)
				}
				c.FinishedAt = time.Now()
				log.Printf("Done %v in %v", c.SourceID, c.FinishedAt.Sub(c.StartedAt))

				previous = c
			}
		}
	}()

	ctx := context.Background()
	devices, err := j.be.ListAllDeviceSources(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		delayed <- ingestion.NewSourceChange(int64(device.ID))
	}

	return nil
}
