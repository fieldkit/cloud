package main

import (
	"context"
	"fmt"

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

	visitor := &noopVisitor{}

	rw := NewRecordWalker(db)

	if err := rw.WalkStation(ctx, 12, visitor); err != nil {
		panic(err)
	}

	info, err := rw.Info(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ok %v data (%v meta)\n", info.DataRecords, info.MetaRecords)
}
