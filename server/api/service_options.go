package api

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

type ControllerOptions struct {
	Config       *ApiConfiguration
	Session      *session.Session
	Database     *sqlxcache.DB
	Backend      *backend.Backend
	JWTHMACKey   []byte
	Emailer      email.Emailer
	Domain       string
	PortalDomain string
	Metrics      *logging.Metrics
	Publisher    jobs.MessagePublisher
	// Twitter
	ConsumerKey    string
	ConsumerSecret string
}

func CreateServiceOptions(ctx context.Context, database *sqlxcache.DB, be *backend.Backend, publisher jobs.MessagePublisher, awsSession *session.Session, config *ApiConfiguration, metrics *logging.Metrics) (controllerOptions *ControllerOptions, err error) {
	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		return nil, err
	}

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		return nil, err
	}

	controllerOptions = &ControllerOptions{
		Session:      awsSession,
		Database:     database,
		Backend:      be,
		Emailer:      emailer,
		JWTHMACKey:   jwtHMACKey,
		Domain:       config.Domain,
		PortalDomain: config.PortalDomain,
		Metrics:      metrics,
		Config:       config,
		Publisher:    publisher,
	}

	return
}
