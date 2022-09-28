package logging

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
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
		return statsd.New(ms.Address)
	}
	return statsd.New(ms.Address)
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

func (m *Metrics) name(path string) string {
	return m.Settings.Prefix + "." + path
}

func (m *Metrics) AuthTry() {
	m.SC.Incr(m.name("api.auth.try"), []string{}, 1)
}

func (m *Metrics) AuthSuccess() {
	m.SC.Incr(m.name("api.auth.success"), []string{}, 1)
}

func (m *Metrics) AuthRefreshTry() {
	m.SC.Incr(m.name("api.auth.refresh.try"), []string{}, 1)
}

func (m *Metrics) AuthRefreshSuccess() {
	m.SC.Incr(m.name("api.auth.refresh.success"), []string{}, 1)
}

func (m *Metrics) Ingested(blocks, bytes int) {
	m.SC.Count(m.name("api.ingestion.blocks"), int64(blocks), []string{}, 1)
	m.SC.Count(m.name("api.ingestion.bytes"), int64(bytes), []string{}, 1)
}

func (m *Metrics) DataErrorsParsing() {
	m.SC.Incr(m.name("api.data.errors.parsing"), []string{}, 1)
}

func (m *Metrics) DataErrorsMissingMeta() {
	m.SC.Incr(m.name("api.data.errors.meta"), []string{}, 1)
}

func (m *Metrics) MessagePublished() {
	m.SC.Incr(m.name("messages.published"), []string{}, 1)
}

type KeyMeta struct {
	key  string
	tags []string
}

type Timing struct {
	m           *Metrics
	start       time.Time
	timingKeys  []KeyMeta
	counterKeys []KeyMeta
}

func (t *Timing) Send() {
	elapsed := time.Now().UTC().Sub(t.start)

	for _, keyMeta := range t.timingKeys {
		t.m.SC.Timing(t.m.name(keyMeta.key), elapsed, keyMeta.tags, 1)
	}
	for _, keyMeta := range t.counterKeys {
		t.m.SC.Incr(t.m.name(keyMeta.key), keyMeta.tags, 1)
	}
}

func (m *Metrics) FileUpload() *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "files.uploading.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "files.uploaded", tags: []string{}}},
	}
}

func (m *Metrics) TailMultiQuery(batch int) *Timing {
	m.SC.Count(m.name("api.data.tail.multi.query.batch"), int64(batch), []string{}, 1)

	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.tail.multi.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.tail.multi.query", tags: []string{}}},
	}
}

func (m *Metrics) RecentlyMultiQuery(batch int) *Timing {
	m.SC.Count(m.name("api.data.recently.multi.query.time"), int64(batch), []string{}, 1)

	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.recently.multi.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.recently.multi.query", tags: []string{}}},
	}
}

func (m *Metrics) LastTimesQuery(batch int) *Timing {
	m.SC.Count(m.name("api.data.lasttimes.multi.query.time"), int64(batch), []string{}, 1)

	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.lasttimes.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.lasttimes.query", tags: []string{}}},
	}
}

func (m *Metrics) DataRangesQuery() *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.ranges.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.ranges.query", tags: []string{}}},
	}
}

func (m *Metrics) DailyQuery() *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.daily.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.daily.query", tags: []string{}}},
	}
}

func (m *Metrics) TailQuery() *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.tail.query.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.tail.query", tags: []string{}}},
	}
}

func (m *Metrics) DataQuery(aggregate string) *Timing {
	aggregateTag := fmt.Sprintf("aggregate:%s", aggregate)
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.data.query.time", tags: []string{aggregateTag}}, KeyMeta{key: fmt.Sprintf("api.data.query.%s.time", aggregate), tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.data.query", tags: []string{aggregateTag}}, KeyMeta{key: fmt.Sprintf("api.data.query.%s", aggregate), tags: []string{}}},
	}
}

func (m *Metrics) HandleMessage(jobType string) *Timing {
	jobTypeTag := fmt.Sprintf("job_type:%s", jobType)
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "messages.handling.time", tags: []string{jobTypeTag}}, KeyMeta{key: fmt.Sprintf("messages.%s.handling.time", jobType), tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "messages.processed", tags: []string{jobTypeTag}}, KeyMeta{key: fmt.Sprintf("messages.%s.processed", jobType), tags: []string{}}},
	}
}

func (m *Metrics) ThirdPartyLocationDescribe() *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: "api.thirdparty.location.describe.time", tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: "api.thirdparty.location.describe.queries", tags: []string{}}},
	}
}

func (m *Metrics) ThirdPartyLocation(provider string) *Timing {
	return &Timing{
		m:           m,
		start:       time.Now().UTC(),
		timingKeys:  []KeyMeta{KeyMeta{key: fmt.Sprintf("api.thirdparty.location.%s.time", provider), tags: []string{}}},
		counterKeys: []KeyMeta{KeyMeta{key: fmt.Sprintf("api.thirdparty.location.%s.queries", provider), tags: []string{}}},
	}
}

func (m *Metrics) PileHit(pile string) {
	m.SC.Incr(m.name(fmt.Sprintf("api.pile.%s.hit", pile)), []string{}, 1)
}

func (m *Metrics) PileMiss(pile string) {
	m.SC.Incr(m.name(fmt.Sprintf("api.pile.%s.miss", pile)), []string{}, 1)
}

func (m *Metrics) PileBytes(pile string, bytes int64) {
	m.SC.Gauge(m.name(fmt.Sprintf("api.pile.%s.bytes", pile)), float64(bytes), []string{}, 1)
}

func (m *Metrics) UserValidated() {
	m.SC.Incr(m.name("api.users.validated"), []string{}, 1)
}

func (m *Metrics) UserAdded() {
	m.SC.Incr(m.name("api.users.added"), []string{}, 1)
}

func (m *Metrics) EmailVerificationSent() {
	m.SC.Incr(m.name("emails.verification"), []string{}, 1)
}

func (m *Metrics) EmailRecoverPasswordSent() {
	m.SC.Incr(m.name("emails.password.recover"), []string{}, 1)
}

func (m *Metrics) DataErrorsUnknown() {
	m.SC.Incr(m.name("api.data.errors.unknown"), []string{}, 1)
}

func (m *Metrics) IngestionDevice(id []byte) {
	m.SC.Set(m.name("api.ingestion.devices"), hex.EncodeToString(id), []string{}, 1)
}

func (m *Metrics) UserID(id int32) {
	m.SC.Set(m.name("api.users"), fmt.Sprintf("%d", id), []string{}, 1)
}

func (m *Metrics) GatherMetrics(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		h.ServeHTTP(w, r)

		elapsed := time.Now().UTC().Sub(start)

		if IsWebSocket(r) {
			m.SC.Timing(m.name("ws.req.time"), elapsed, []string{}, 1)
		} else {
			if IsIngestion(r) {
				m.SC.Timing(m.name("http.ingestion.time"), elapsed, []string{}, 1)
			} else {
				m.SC.Timing(m.name("http.req.time"), elapsed, []string{}, 1)
			}
		}
	})
}

func (m *Metrics) RecordsViewed(records int) {
	m.SC.Count(m.name("api.data.records.viewed"), int64(records), []string{}, 1)
}

func (m *Metrics) ReadingsViewed(readings int) {
	m.SC.Count(m.name("api.data.readings.viewed"), int64(readings), []string{}, 1)
}

func IsWebSocket(r *http.Request) bool {
	return r.Header.Get("Upgrade") == "websocket"
}

func IsIngestion(r *http.Request) bool {
	return r.URL.Path == "/ingestion"
}
