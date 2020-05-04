package tests

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bxcodec/faker/v3"

	"github.com/fieldkit/cloud/server/data"
)

type TestEnv struct {
	Ctx context.Context
	DB  *sqlxcache.DB
}

var (
	globalEnv *TestEnv
)

func NewTestEnv(originalUrl string) (e *TestEnv, err error) {
	if globalEnv != nil {
		log.Printf("using existing test env")
		return globalEnv, nil
	}

	ctx := context.Background()

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

	migrater, err := migrate.New("file://../../../migrations", testUrl)
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

type FakeStations struct {
	OwnerID   int32
	ProjectID int32
	Stations  []*data.Station
}

func (e *TestEnv) AddStations(number int) (*FakeStations, error) {
	email := faker.Email()
	owner := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Password: []byte("IGNORE"),
		Bio:      faker.Sentence(),
	}

	if err := e.DB.NamedGetContext(e.Ctx, owner, `
		INSERT INTO fieldkit.user (name, username, email, password, bio)
		VALUES (:name, :email, :email, :password, :bio)
		RETURNING *
		`, owner); err != nil {
		return nil, err
	}

	name := faker.Name()

	project := &data.Project{
		Name: name + " Project",
		Slug: name,
	}

	if err := e.DB.NamedGetContext(e.Ctx, project, `
		INSERT INTO fieldkit.project (name, slug)
		VALUES (:name, :slug)
		RETURNING *
		`, project); err != nil {
		return nil, err
	}

	if _, err := e.DB.ExecContext(e.Ctx, `
		INSERT INTO fieldkit.project_user (project_id, user_id) VALUES ($1, $2)
		`, project.ID, owner.ID); err != nil {
		return nil, err
	}

	stations := []*data.Station{}

	for i := 0; i < number; i += 1 {
		name := fmt.Sprintf("%s #%d", owner.Name, i)

		hasher := sha1.New()
		hasher.Write([]byte(name))
		deviceID := hasher.Sum(nil)

		station := &data.Station{
			OwnerID:  owner.ID,
			Name:     name,
			DeviceID: deviceID,
		}

		if err := e.DB.NamedGetContext(e.Ctx, station, `
			INSERT INTO fieldkit.station (name, device_id, owner_id, status_json)
			VALUES (:name, :device_id, :owner_id, :status_json)
			RETURNING *
		`, station); err != nil {
			return nil, err
		}

		if _, err := e.DB.ExecContext(e.Ctx, `
			INSERT INTO fieldkit.project_station (project_id, station_id) VALUES ($1, $2)
			`, project.ID, station.ID); err != nil {
			return nil, err
		}

		stations = append(stations, station)
	}

	return &FakeStations{
		OwnerID:   owner.ID,
		ProjectID: project.ID,
		Stations:  stations,
	}, nil
}

func changeConnectionStringDatabase(original, newDatabase string) (string, error) {
	p, err := url.Parse(original)
	if err != nil {
		return "", err
	}

	p.Path = newDatabase

	return p.String(), nil
}
