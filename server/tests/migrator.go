package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-pg/pg/v10/orm"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

var (
	registered bool
)

func tryMigrate(url string) error {
	migrationsDir, err := findMigrationsDirectory()
	if err != nil {
		return err
	}

	log.Printf("trying to migrate...")
	log.Printf("postgres = %s", url)
	log.Printf("migrations = %s", migrationsDir)

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		log.Fatal(err)
	}

	if !registered {
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

		registered = true
	}

	options, err := pg.ParseURL(url)
	if err != nil {
		return err
	}

	db := pg.Connect(options)

	if err := migrations.Run(db, migrationsDir, []string{"", "migrate"}); err != nil {
		return err
	}

	return nil
}

func findMigrationsDirectory() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to find migrations directory: %v", err)
	}

	for {
		test := filepath.Join(path, "migrations")
		if _, err := os.Stat(test); !os.IsNotExist(err) {
			return test, nil
		}

		path = filepath.Dir(path)
	}

	return "", fmt.Errorf("unable to find migrations directory")
}
