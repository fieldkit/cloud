package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-pg/pg/v10/orm"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func main() {
	options := &pg.Options{
		Addr:     "",
		User:     "",
		Database: "",
		Password: "",
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			log.Printf("Creating schema...")

			if _, err := conn.Exec("CREATE SCHEMA IF NOT EXISTS fieldkit"); err != nil {
				return fmt.Errorf("error creating: %v", err)
			}

			if _, err := conn.Exec("GRANT USAGE ON SCHEMA fieldkit TO fieldkit"); err != nil {
				return fmt.Errorf("error granting: %v", err)
			}

			if _, err := conn.Exec("GRANT CREATE ON SCHEMA fieldkit TO fieldkit"); err != nil {
				return fmt.Errorf("error granting: %v", err)
			}

			log.Printf("Done creating schema...")

			return nil
		},
	}

	url := os.Getenv("MIGRATE_DATABASE_URL")
	if url == "" {
		log.Fatalln("MIGRATE_DATABASE_URL is requied")
	}

	o, err := pg.ParseURL(url)
	if err != nil {
		log.Fatalln(err)
	}

	options = o

	directory := os.Getenv("MIGRATE_PATH")
	if directory == "" {
		log.Fatalln("MIGRATE_PATH is requied")
	}

	log.Printf("Scanning %s...", directory)

	files, err := filepath.Glob(filepath.Join(directory, "*.up.sql"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d migration(s) in %s...", len(files), directory)

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

		_, fileOnly := filepath.Split(file)

		migrations.Register(fileOnly, up, down, opts)
	}

	db := pg.Connect(options)

	if err := migrations.Run(db, directory, os.Args); err != nil {
		log.Fatalln(err)
	}
}
