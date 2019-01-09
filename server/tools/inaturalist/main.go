package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/conservify/gonaturalist"

	fk "github.com/fieldkit/cloud/server/api/client"
	fktesting "github.com/fieldkit/cloud/server/tools"

	"github.com/fieldkit/cloud/server/backend/ingestion/formatting"
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/logging"
)

var (
	GreekTheaterLocation = []float64{-118.2871352, 34.1168852}
)

type options struct {
	Scheme   string
	Host     string
	Username string
	Password string

	Project    string
	Expedition string
	DeviceName string
	DeviceId   string

	DeleteComments   bool
	NumberOfReadings int
}

func findObservation(ctx context.Context, log *zap.SugaredLogger, o *options, inc *gonaturalist.Client) (obs *gonaturalist.SimpleObservation, err error) {
	user, err := inc.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	log.Infow("CurrentUser", "user", user)

	// Now we need an observation.
	log.Infow("GetObservationsByUsername", "login", user.Login)
	myObservations, err := inc.GetObservationsByUsername(user.Login)
	if err != nil {
		return nil, err
	}

	log.Infow("My observations", "numberOfObservations", len(myObservations.Observations))

	if len(myObservations.Observations) == 0 {
		addObservation := gonaturalist.AddObservationOpt{
			SpeciesGuess:       "Duck",
			ObservedOnString:   time.Now(),
			Description:        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris quis mi quam. Ut venenatis pulvinar magna, sit amet auctor mi vehicula in.",
			Longitude:          GreekTheaterLocation[0],
			Latitude:           GreekTheaterLocation[1],
			PositionalAccuracy: 1,
		}
		obs, err := inc.AddObservation(&addObservation)
		if err != nil {
			return nil, err
		}
		log.Infow("Added", "observation", obs)

		return obs, nil
	}

	return myObservations.Observations[0], nil
}

func deleteMyOwnComments(ctx context.Context, log *zap.SugaredLogger, o *options, inc *gonaturalist.Client, obs *gonaturalist.SimpleObservation) (err error) {
	user, err := inc.GetCurrentUser()
	if err != nil {
		return err
	}

	log.Infow("GetObservationComments", "observationId", obs.Id)
	full, err := inc.GetObservationComments(obs.Id)
	if err != nil {
		return err
	}

	for _, c := range full {
		if c.UserId == user.Id {
			if o.DeleteComments {
				log.Infow("Deleting", "comment", c)
				err = inc.DeleteComment(c.Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func updateObservation(ctx context.Context, log *zap.SugaredLogger, o *options, inc *gonaturalist.Client, obs *gonaturalist.SimpleObservation) (err error) {
	now := time.Now()
	updateObservation := gonaturalist.UpdateObservationOpt{
		Id:               obs.Id,
		ObservedOnString: &now,
		Description:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris quis mi quam. Ut venenatis pulvinar magna, sit amet auctor mi vehicula in." + "\n" + time.Now().String(),
	}
	err = inc.UpdateObservation(&updateObservation)

	return err
}

func addFakeRecords(ctx context.Context, log *zap.SugaredLogger, o *options, fkc *fk.Client) error {
	for i := 0; i < o.NumberOfReadings; i += 1 {
		m := formatting.HttpJsonMessage{
			Location: GreekTheaterLocation,
			Time:     time.Now().UTC().Unix(),
			Device:   o.DeviceId,
			Fixed:    true,
			Stream:   "",
			Modules:  []string{"Fake"},
			Values: map[string]interface{}{
				"cpu":      rand.Uint32() % 100,
				"temp":     rand.Uint32() % 40,
				"humidity": rand.Uint32() % 100,
			},
		}

		log.Infow("Adding fake reading", "deviceId", o.DeviceId)

		bytes, err := json.Marshal(m)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		err = o.postRawMessage(bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func triggerRefresh(ctx context.Context, log *zap.SugaredLogger, o *options) error {
	url := o.getUrl("/tasks/five")
	log.Infow("Refreshing", "url", url)
	_, err := http.Get(url)
	return err
}

func (o *options) getUrl(path string) string {
	return fmt.Sprintf("%s://%s%s", o.Scheme, o.Host, path)
}

func (o *options) postRawMessage(body []byte) error {
	url := o.getUrl("/messages/ingestion")
	if false {
		fmt.Printf("curl -X POST -H 'Content-Type: %s' -d '%s' %s\n", formatting.HttpProviderJsonContentType, string(body), url)
	} else {
		bytes := bytes.NewBufferString(string(body))
		url := url
		url += "?token=" + "IGNORED"
		_, err := http.Post(url, formatting.HttpProviderJsonContentType, bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	o := options{}

	flag.StringVar(&o.Scheme, "scheme", "http", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "fk instance hostname")
	flag.StringVar(&o.Username, "username", "demo-user", "username")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password")

	flag.StringVar(&o.Project, "project", "www", "project")
	flag.StringVar(&o.Expedition, "expedition", "", "expedition")
	flag.StringVar(&o.DeviceName, "device-name", "fkn-test", "device name")
	flag.StringVar(&o.DeviceId, "device-id", "1c7c6a39-4208-4978-811f-5cc75c83ca13", "device name")

	flag.BoolVar(&o.DeleteComments, "delete-comments", false, "delete comments")
	flag.IntVar(&o.NumberOfReadings, "number-of-readings", 10, "number of readings")

	flag.Parse()

	log, err := logging.Configure(false)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	fkc, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		panic(err)
	}

	device, err := fktesting.CreateWebDevice(ctx, fkc, o.Project, o.DeviceName, o.DeviceId, "")
	if err != nil {
		panic(err)
	}
	log.Infow("Device", "device", device)

	config := inaturalist.ConfigConservify
	authenticator := gonaturalist.NewAuthenticatorAtCustomRoot(config.ApplicationId, config.Secret, config.RedirectUrl, config.RootUrl)
	inc := authenticator.NewClientWithAccessToken(config.AccessToken, &gonaturalist.NoopCallbacks{})

	obs, err := findObservation(ctx, log, &o, inc)
	if err != nil {
		panic(err)
	}

	err = deleteMyOwnComments(ctx, log, &o, inc, obs)
	if err != nil {
		panic(err)
	}

	err = updateObservation(ctx, log, &o, inc, obs)
	if err != nil {
		panic(err)
	}

	err = addFakeRecords(ctx, log, &o, fkc)
	if err != nil {
		panic(err)
	}

	err = triggerRefresh(ctx, log, &o)
	if err != nil {
		panic(err)
	}

	_ = fkc
	_ = log
	_ = obs
}
