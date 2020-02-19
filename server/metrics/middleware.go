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
	c, err := statsd.New(statsd.Address(settings.Address))
	if err != nil {
		log := Logger(ctx).Sugar()
		log.Warnw("error starting statsd", "error", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := c.NewTiming()

		next.ServeHTTP(w, r)

		t.Send(fmt.Sprintf("%s.http.req.time", settings.Prefix))
	})
}
