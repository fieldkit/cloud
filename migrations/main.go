package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-pg/pg/v10/orm"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "migrations"

type options struct {
	Path       string
	SearchPath string
}

func main() {
	o := &options{}

	flag.StringVar(&o.Path, "path", "./", "path to migrations")
	flag.StringVar(&o.SearchPath, "search-path", "fieldkit,public", "search path to apply")

	flag.Parse()

	options := &pg.Options{
		Addr:     "",
		User:     "",
		Database: "",
		Password: "",
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			log.Printf("creating schema...")

			if _, err := conn.Exec("CREATE SCHEMA IF NOT EXISTS fieldkit"); err != nil {
				return fmt.Errorf("error creating: %v", err)
			}

			if _, err := conn.Exec("GRANT USAGE ON SCHEMA fieldkit TO fieldkit"); err != nil {
				return fmt.Errorf("error granting: %v", err)
			}

			if _, err := conn.Exec("GRANT CREATE ON SCHEMA fieldkit TO fieldkit"); err != nil {
				return fmt.Errorf("error granting: %v", err)
			}

			log.Printf("done creating schema...")

			return nil
		},
	}

	url := os.Getenv("PGURL")
	if url != "" {
		o, err := pg.ParseURL(url)
		if err != nil {
			log.Fatalln(err)
		}

		options = o
	}

	files, err := filepath.Glob(filepath.Join(o.Path, "*.up.sql"))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		text := string(data)

		up := func(db orm.DB) error {
			_, err := db.Exec(text)
			return err
		}

		down := func(db orm.DB) error {
			return err
		}

		opts := migrations.MigrationOptions{}

		migrations.Register(file, up, down, opts)
	}

	db := pg.Connect(options)

	if err := migrations.Run(db, directory, os.Args); err != nil {
		log.Fatalln(err)
	}
}
