package metrics

import (
	"context"
	"net/http"

	"gopkg.in/alexcesaro/statsd.v2"
)

type MetricsSettings struct {
	Address string
}

func GatherMetrics(ctx context.Context, settings MetricsSettings, next http.Handler) http.Handler {
	c, err := statsd.New(statsd.Address(settings.Address))
	if err != nil {
		log := Logger(ctx).Sugar()
		log.Warnw("error starting statsd", "error", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := c.NewTiming()

		next.ServeHTTP(w, r)

		t.Send("fk.http.req.time")
	})
}
