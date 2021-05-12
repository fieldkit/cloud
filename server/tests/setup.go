package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/data"
)

type TestEnv struct {
	Ctx         context.Context
	DB          *sqlxcache.DB
	PostgresURL string
	SessionKey  string
	JWTHMACKey  []byte
	Seed        int64
}

type TestConfig struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1:5432/fieldkit?sslmode=disable&search_path=public" required:"true"`
}

var (
	globalEnv *TestEnv
)

func tryMigrate(url string) error {
	migrationsDir, err := findMigrationsDirectory()
	if err != nil {
		return err
	}

	log.Printf("trying to migrate...")
	log.Printf("postgres = %s", url)
	log.Printf("migrations = %s", migrationsDir)

	migrater, err := migrate.New("file://"+migrationsDir, url)
	if err != nil {
		return err
	}

	migrater.Log = &MigratorLog{}

	if err := migrater.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

func NewTestEnv() (e *TestEnv, err error) {
	if globalEnv != nil {
		log.Printf("using existing test env")
		return globalEnv, nil
	}

	config := &TestConfig{}

	if err := envconfig.Process("FIELDKIT", config); err != nil {
		return nil, err
	}

	logging.Configure(false, "tests")

	ctx := context.Background()

	originalDb, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		return nil, err
	}

	databaseName := "fktest"
	testUrl, err := changeConnectionStringDatabase(config.PostgresURL, databaseName)
	if err != nil {
		return nil, err
	}

	if err := tryMigrate(testUrl); err != nil {
		log.Printf("creating test database")

		if _, err := originalDb.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", databaseName)); err != nil {
			return nil, err
		}

		if _, err := originalDb.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", databaseName)); err != nil {
			return nil, err
		}

		if err := tryMigrate(testUrl); err != nil {
			return nil, err
		}
	}

	testDb, err := sqlxcache.Open("postgres", testUrl)
	if err != nil {
		return nil, err
	}

	testSessionKey := "AzhG6aOqy2jbx7KSrmkOJGe+RY75nsZsr+ByneiMLBo="
	jwtHMACKey, err := base64.StdEncoding.DecodeString(testSessionKey)
	if err != nil {
		return nil, err
	}

	seed := time.Now().UnixNano()

	e = &TestEnv{
		Ctx:         ctx,
		DB:          testDb,
		PostgresURL: testUrl,
		SessionKey:  testSessionKey,
		JWTHMACKey:  jwtHMACKey,
		Seed:        seed,
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

func (e *TestEnv) NewTokenForUser(user *data.User) string {
	now := time.Now()
	refreshToken, err := data.NewRefreshToken(user.ID, 20, now.Add(time.Duration(72)*time.Hour))
	if err != nil {
		panic(err)
	}

	token := user.NewToken(now, refreshToken)
	signedToken, err := token.SignedString(e.JWTHMACKey)

	return signedToken
}

func (e *TestEnv) NewAuthorizationHeaderForUser(user *data.User) string {
	signedToken := e.NewTokenForUser(user)
	return "Bearer " + signedToken
}

func (e *TestEnv) NewAuthorizationHeader() string {
	user := &data.User{
		ID:    1,
		Admin: false,
		Email: "user@fieldkit.org",
	}

	return e.NewAuthorizationHeaderForUser(user)
}

func (e *TestEnv) NewAuthorizationHeaderForAdmin() string {
	user := &data.User{
		ID:    1,
		Admin: true,
		Email: "admin@fieldkit.org",
	}

	return e.NewAuthorizationHeaderForUser(user)
}
