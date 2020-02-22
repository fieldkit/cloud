package logging

import (
	"context"
	"encoding/hex"
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
		return statsd.New(statsd.Prefix(ms.Prefix))
	}
	return statsd.New(statsd.Prefix(ms.Prefix), statsd.Address(ms.Address))
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

func (m *Metrics) AuthTry() {
	m.SC.Increment("api.auth.try")
}

func (m *Metrics) AuthSuccess() {
	m.SC.Increment("api.auth.success")
}

func (m *Metrics) AuthRefreshTry() {
	m.SC.Increment("api.auth.refresh.try")
}

func (m *Metrics) AuthRefreshSuccess() {
	m.SC.Increment("api.auth.refresh.success")
}

func (m *Metrics) Ingested(blocks, bytes int) {
	m.SC.Count("api.ingestion.blocks", blocks)
	m.SC.Count("api.ingestion.bytes", bytes)
}

func (m *Metrics) IngestionDevice(id []byte) {
	m.SC.Unique("api.ingestion.devices", hex.EncodeToString(id))
}

func (m *Metrics) UserID(id int32) {
	m.SC.Unique("api.users", fmt.Sprintf("%d", id))
}

func (m *Metrics) GatherMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := m.SC.NewTiming()

		next.ServeHTTP(w, r)

		t.Send("http.req.time")
	})
}
