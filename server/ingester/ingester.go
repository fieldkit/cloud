package ingester

import (
	"context"
	"net/http"
	"time"

	"github.com/goadesign/goa"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	_ "github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/logging"
)

type IngesterOptions struct {
	Database                 *sqlxcache.DB
	AwsSession               *session.Session
	AuthenticationMiddleware goa.Middleware
}

func authenticate(middleware goa.Middleware, next goa.Handler) goa.Handler {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
		ctx = goa.WithRequiredScopes(ctx, []string{"api:access"})
		return middleware(next)(ctx, res, req)
	}
}

func Ingester(ctx context.Context, o *IngesterOptions) http.Handler {
	handler := authenticate(o.AuthenticationMiddleware, func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		log := logging.Logger(ctx).Sugar()

		log.Infow("begin")

		return nil
	})

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startedAt := time.Now()

		ctx := context.Background()
		log := logging.Logger(ctx).Sugar()

		err := handler(ctx, w, req)
		if err != nil {
			log.Errorw("completed", "error", err, "time", time.Since(startedAt).String())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
