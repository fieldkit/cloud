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
	Config          *ApiConfiguration
	Session         *session.Session
	Database        *sqlxcache.DB
	Backend         *backend.Backend
	ConcatWorkers   *backend.ConcatenationWorkers
	JWTHMACKey      []byte
	Emailer         email.Emailer
	Domain          string
	Metrics         *logging.Metrics
	Publisher       jobs.MessagePublisher
	StreamProcessor backend.StreamProcessor

	// Twitter
	ConsumerKey    string
	ConsumerSecret string
}

func CreateServiceOptions(ctx context.Context, database *sqlxcache.DB, be *backend.Backend, awsSession *session.Session, ingester *backend.StreamIngester,
	publisher jobs.MessagePublisher, cw *backend.ConcatenationWorkers, config *ApiConfiguration, metrics *logging.Metrics) (controllerOptions *ControllerOptions, err error) {

	log := Logger(ctx).Sugar()

	log.Infow("config", "config", config)

	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		return nil, err
	}

	streamProcessor := backend.NewS3StreamProcessor(awsSession, ingester, config.BucketName)

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		return nil, err
	}

	controllerOptions = &ControllerOptions{
		Session:         awsSession,
		Database:        database,
		Backend:         be,
		Emailer:         emailer,
		JWTHMACKey:      jwtHMACKey,
		Domain:          config.Domain,
		Metrics:         metrics,
		StreamProcessor: streamProcessor,
		Config:          config,
		ConcatWorkers:   cw,
		Publisher:       publisher,
	}

	return
}
