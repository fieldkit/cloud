package tests

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestEnv struct {
	Ctx context.Context
	DB  *sqlxcache.DB
}

const PostgresURL = "postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public"

var (
	globalEnv *TestEnv
)

func NewTestEnv() (e *TestEnv, err error) {
	if globalEnv != nil {
		log.Printf("using existing test env")
		return globalEnv, nil
	}

	ctx := context.Background()

	originalUrl := PostgresURL
	originalDb, err := sqlxcache.Open("postgres", originalUrl)
	if err != nil {
		return nil, err
	}

	databaseName := "fktest"
	testUrl, err := changeConnectionStringDatabase(originalUrl, databaseName)
	if err != nil {
		return nil, err
	}

	log.Printf("creating test database")

	if _, err := originalDb.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", databaseName)); err != nil {
		return nil, err
	}

	if _, err := originalDb.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", databaseName)); err != nil {
		return nil, err
	}

	migrationsDir, err := findMigrationsDirectory()
	if err != nil {
		return nil, err
	}

	migrater, err := migrate.New("file://"+migrationsDir, testUrl)
	if err != nil {
		return nil, err
	}

	migrater.Log = &MigratorLog{}

	log.Printf("migrating test database")

	if err := migrater.Up(); err != nil {
		return nil, err
	}

	testDb, err := sqlxcache.Open("postgres", testUrl)
	if err != nil {
		return nil, err
	}

	e = &TestEnv{
		Ctx: ctx,
		DB:  testDb,
	}

	globalEnv = e

	return
}

func changeConnectionStringDatabase(original, newDatabase string) (string, error) {
	p, err := url.Parse(original)
	if err != nil {
		return "", err
	}

	p.Path = newDatabase

	return p.String(), nil
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
