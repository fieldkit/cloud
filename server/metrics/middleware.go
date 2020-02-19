package metrics

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/alexcesaro/statsd.v2"
)

type MetricsSettings struct {
	Prefix  string
	Address string
}

func GatherMetrics(ctx context.Context, settings MetricsSettings, next http.Handler) http.Handler {
	log := Logger(ctx).Sugar()

	if settings.Address == "" {
		log.Infow("statsd: skipping")
		return next
	}

	log.Infow("statsd", "address", settings.Address)

	c, err := statsd.New(statsd.Address(settings.Address))
	if err != nil {
		log.Warnw("statsd: error starting", "error", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := c.NewTiming()

		next.ServeHTTP(w, r)

		t.Send(fmt.Sprintf("%s.http.req.time", settings.Prefix))
	})
}
