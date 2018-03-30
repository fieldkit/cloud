package backend

import (
	"context"
	"log"
	_ "sync"
	"time"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

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
