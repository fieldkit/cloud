package main

import (
	"context"
	"flag"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/kelseyhightower/envconfig"

	"github.com/bxcodec/faker/v3"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"
)

type Options struct {
	PostgresURL     string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	WaitForDatabase bool
	Anonymize       bool
}

func sanitize(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	users := []*data.User{}
	if err := db.SelectContext(ctx, &users, `SELECT * FROM fieldkit.user ORDER BY id`); err != nil {
		return err
	}

	for _, user := range users {
		user.SetPassword("asdfasdfasdf")

		if !strings.Contains(user.Email, "@conservify.org") && !strings.Contains(user.Email, "@fieldkit.org") {
			user.Name = faker.Name()
			user.Bio = faker.Sentence()
			user.Email = faker.Email()
		} else {
			log.Infow("keeping", "user_id", user.ID, "email", user.Email)
		}

		if _, err := db.NamedExecContext(ctx, `
			UPDATE fieldkit.user SET password = :password, name = :name, username = :email, email = :email WHERE id = :id
			`, user); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{
		WaitForDatabase: true,
		Anonymize:       false,
	}

	flag.BoolVar(&options.Anonymize, "anonymize", false, "")
	flag.BoolVar(&options.WaitForDatabase, "waiting", false, "")

	flag.Parse()

	logging.Configure(false, "sanitizer")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	for {
		if err := sanitize(ctx, options); err != nil {
			// NOTE May be a good idea to check this error?
			log.Infow("error", "error", err)
			if !options.WaitForDatabase {
				panic(err)
			}
		} else {
			break
		}

		time.Sleep(1 * time.Second)
	}

	log.Infow("done")
}
