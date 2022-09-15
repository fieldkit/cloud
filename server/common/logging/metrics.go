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

func (m *Metrics) DataErrorsParsing() {
	m.SC.Increment("api.data.errors.parsing")
}

func (m *Metrics) DataErrorsMissingMeta() {
	m.SC.Increment("api.data.errors.meta")
}

func (m *Metrics) MessagePublished() {
	m.SC.Increment("messages.published")
}

type Timing struct {
	sc          *statsd.Client
	timer       statsd.Timing
	timingKeys  []string
	counterKeys []string
}

func (t *Timing) Send() {
	for _, key := range t.timingKeys {
		t.timer.Send(key)
	}
	for _, key := range t.counterKeys {
		t.sc.Increment(key)
	}
}

func (m *Metrics) FileUpload() *Timing {
	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"files.uploading.time"},
		counterKeys: []string{"files.uploaded"},
	}
}

func (m *Metrics) TailMultiQuery(batch int) *Timing {
	m.SC.Count("api.data.tail.multi.query.batch", batch)

	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.tail.multi.query.time"},
		counterKeys: []string{"api.data.tail.multi.query"},
	}
}

func (m *Metrics) RecentlyMultiQuery(batch int) *Timing {
	m.SC.Count("api.data.recently.multi.query.time", batch)

	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.recently.multi.query.time"},
		counterKeys: []string{"api.data.recently.multi.query"},
	}
}

func (m *Metrics) LastTimesQuery(batch int) *Timing {
	m.SC.Count("api.data.lasttimes.multi.query.time", batch)

	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.lasttimes.query.time"},
		counterKeys: []string{"api.data.lasttimes.query"},
	}
}

func (m *Metrics) DailyQuery() *Timing {
	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.daily.query.time"},
		counterKeys: []string{"api.data.daily.query"},
	}
}

func (m *Metrics) TailQuery() *Timing {
	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.tail.query.time"},
		counterKeys: []string{"api.data.tail.query"},
	}
}

func (m *Metrics) DataQuery(aggregate string) *Timing {
	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"api.data.query.time", fmt.Sprintf("api.data.query.%s.time", aggregate)},
		counterKeys: []string{"api.data.query", fmt.Sprintf("api.data.query.%s", aggregate)},
	}
}

func (m *Metrics) HandleMessage(jobType string) *Timing {
	timer := m.SC.NewTiming()

	return &Timing{
		sc:          m.SC,
		timer:       timer,
		timingKeys:  []string{"messages.handling.time", fmt.Sprintf("messages.%s.handling.time", jobType)},
		counterKeys: []string{"messages.processed", fmt.Sprintf("messages.%s.processed", jobType)},
	}
}

func (m *Metrics) PileHit(pile string) {
	m.SC.Increment(fmt.Sprintf("api.pile.%s.hit", pile))
}

func (m *Metrics) PileMiss(pile string) {
	m.SC.Increment(fmt.Sprintf("api.pile.%s.miss", pile))
}

func (m *Metrics) PileBytes(pile string, bytes int64) {
	m.SC.Gauge(fmt.Sprintf("api.pile.%s.bytes", pile), bytes)
}

func (m *Metrics) UserValidated() {
	m.SC.Increment("api.users.validated")
}

func (m *Metrics) UserAdded() {
	m.SC.Increment("api.users.added")
}

func (m *Metrics) EmailVerificationSent() {
	m.SC.Increment("emails.verification")
}

func (m *Metrics) EmailRecoverPasswordSent() {
	m.SC.Increment("emails.password.recover")
}

func (m *Metrics) DataErrorsUnknown() {
	m.SC.Increment("api.data.errors.unknown")
}

func (m *Metrics) IngestionDevice(id []byte) {
	m.SC.Unique("api.ingestion.devices", hex.EncodeToString(id))
}

func (m *Metrics) UserID(id int32) {
	m.SC.Unique("api.users", fmt.Sprintf("%d", id))
}

func (m *Metrics) GatherMetrics(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := m.SC.NewTiming()

		h.ServeHTTP(w, r)

		if IsWebSocket(r) {
			t.Send("ws.req.time")
		} else {
			if IsIngestion(r) {
				t.Send("http.ingestion.time")
			} else {
				t.Send("http.req.time")
			}
		}
	})
}

func (m *Metrics) RecordsViewed(records int) {
	m.SC.Count("api.data.records.viewed", records)
}

func (m *Metrics) ReadingsViewed(readings int) {
	m.SC.Count("api.data.readings.viewed", readings)
}

func IsWebSocket(r *http.Request) bool {
	return r.Header.Get("Upgrade") == "websocket"
}

func IsIngestion(r *http.Request) bool {
	return r.URL.Path == "/ingestion"
}
