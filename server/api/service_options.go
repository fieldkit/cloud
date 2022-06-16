package api

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/files"
)

type ControllerOptions struct {
	Config       *ApiConfiguration
	Session      *session.Session
	Database     *sqlxcache.DB
	Querier      *data.Querier
	Backend      *backend.Backend
	JWTHMACKey   []byte
	Emailer      email.Emailer
	Domain       string
	PortalDomain string
	Metrics      *logging.Metrics
	Publisher    jobs.MessagePublisher
	MediaFiles   files.FileArchive

	// Twitter
	ConsumerKey    string
	ConsumerSecret string

	// Services
	signer    *Signer
	locations *data.DescribeLocations
	que       *que.Client

	// Subscribed listeners
	subscriptions *Subscriptions
}

func CreateServiceOptions(ctx context.Context, config *ApiConfiguration, database *sqlxcache.DB, be *backend.Backend, publisher jobs.MessagePublisher, mediaFiles files.FileArchive,
	awsSession *session.Session, metrics *logging.Metrics, que *que.Client) (controllerOptions *ControllerOptions, err error) {
	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		return nil, err
	}

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		return nil, err
	}

	locations := data.NewDescribeLocations(config.MapboxToken)

	controllerOptions = &ControllerOptions{
		Session:       awsSession,
		Database:      database,
		Querier:       data.NewQuerier(database),
		Backend:       be,
		Emailer:       emailer,
		JWTHMACKey:    jwtHMACKey,
		Domain:        config.Domain,
		PortalDomain:  config.PortalDomain,
		Metrics:       metrics,
		Config:        config,
		Publisher:     publisher,
		MediaFiles:    mediaFiles,
		signer:        NewSigner(jwtHMACKey),
		locations:     locations,
		que:           que,
		subscriptions: NewSubscriptions(),
	}

	return
}

func (o *ControllerOptions) Close() error {
	return nil
}
