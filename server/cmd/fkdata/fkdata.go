package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"
)

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

func main() {
	ctx := context.Background()

	config := &Config{}
	if err := envconfig.Process("FIELDKIT", config); err != nil {
		panic(err)
	}

	db, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	visitor := NewResolvingVisitor()

	rw := NewRecordWalker(db)

	started := time.Now()

	fmt.Printf("processing\n")

	if err := rw.WalkStation(ctx, 12, visitor); err != nil {
		panic(err)
	}

	info, err := rw.Info(ctx)
	if err != nil {
		panic(err)
	}

	finished := time.Now()

	fmt.Printf("done %v data (%v meta) %v\n", info.DataRecords, info.MetaRecords, finished.Sub(started))
}
