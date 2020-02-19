package logging

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

type Metrics struct {
	SC       *statsd.Client
	Settings *MetricsSettings
}

func newClient(ms *MetricsSettings) (*statsd.Client, error) {
	if ms.Address == "" {
		return statsd.New()
	}
	return statsd.New(statsd.Address(ms.Address))
}

func NewMetrics(ctx context.Context, ms *MetricsSettings) *Metrics {
	log := Logger(ctx).Sugar()

	log.Infow("statsd", "address", ms.Address)

	sc, err := newClient(ms)
	if err != nil {
		log.Warnw("error", "error", err)
	}

	return &Metrics{
		SC:       sc,
		Settings: ms,
	}
}

func (m *Metrics) GatherMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := m.SC.NewTiming()

		next.ServeHTTP(w, r)

		t.Send(fmt.Sprintf("%s.http.req.time", m.Settings.Prefix))
	})
}
