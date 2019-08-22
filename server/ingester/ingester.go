package ingester

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/logging"
)

type IngesterOptions struct {
	Database   *sqlxcache.DB
	AwsSession *session.Session
}

func Ingester(ctx context.Context, o *IngesterOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log := logging.Logger(ctx).Sugar()

		log.Infow("begin")

		http.NotFoundHandler().ServeHTTP(w, req)
	})
}
