package backend

import (
	"context"
	"log"
	_ "sync"
	"time"
)

type NaiveBackgroundJobs struct {
	be *Backend
}

func NewNaiveBackgroundJobs(be *Backend) *NaiveBackgroundJobs {
	return &NaiveBackgroundJobs{
		be: be,
	}
}

func (j *NaiveBackgroundJobs) Start() error {
	delayed := make(chan SourceChange, 100)

	go func() {
		log.Printf("Started background jobs...")

		tick := time.Tick(1000 * time.Millisecond)

		buffer := make(map[int64]SourceChange)

		for {
			select {
			case <-tick:
				for _, value := range buffer {
					delayed <- value
				}
				buffer = make(map[int64]SourceChange)
			case c := <-j.be.SourceChanges:
				buffer[c.SourceID] = c
			}

		}
	}()

	go func() {
		for {
			select {
			case c := <-delayed:
				started := time.Now()
				log.Printf("Processing %v...", c)
				generator := NewPregenerator(j.be)
				ctx := context.Background()
				err := generator.Pregenerate(ctx, c.SourceID)
				if err != nil {
					log.Printf("Error: %v", err)
				}
				log.Printf("Done %v in %v", c, time.Now().Sub(started))
			}
		}
	}()

	ctx := context.Background()
	devices, err := j.be.ListAllDeviceSources(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		log.Printf("Queuing refresh %d...", device.ID)
		delayed <- SourceChange{SourceID: int64(device.ID)}
	}

	return nil
}
