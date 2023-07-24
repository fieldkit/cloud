package logging

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	ids = NewIdGenerator()
)

func LoggingAndInfrastructure(name string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()

			newCtx := WithNewTaskID(r.Context(), ids)
			withCtx := r.WithContext(newCtx)

			log := Logger(newCtx).Named(name).Sugar()
			sanitizedUrl, err := sanitize(r.URL)
			if err != nil {
				log.Warnw("req:sanitize:failed", "error", err)
			}

			log.Infow("req:begin", "req", r.Method+" "+sanitizedUrl, "from", from(r), "ws", IsWebSocket(r))

			cw := CaptureResponse(w)
			h.ServeHTTP(AllowWriteHeaderPrevention(cw), withCtx)

			elapsed := time.Since(started)

			log.Infow("req:done", "status", cw.StatusCode, "bytes", cw.ContentLength,
				"time", fmt.Sprintf("%vns", elapsed.Nanoseconds()),
				"time_human", elapsed.String(),
				"req", r.Method+" "+sanitizedUrl,
				"from", from(r), "ws", IsWebSocket(r))
		})
	}
}

func Monitoring(name string, m *Metrics) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return m.GatherMetrics(LoggingAndInfrastructure(name)(h))
	}
}

// from makes a best effort to compute the request client IP.
func from(req *http.Request) string {
	if f := req.Header.Get("X-Forwarded-For"); f != "" {
		return f
	}
	f := req.RemoteAddr
	ip, _, err := net.SplitHostPort(f)
	if err != nil {
		return f
	}
	return ip
}

func sanitize(orig *url.URL) (string, error) {
	copy, err := url.Parse(orig.String())
	if err != nil {
		return "", err
	}
	q := copy.Query()
	if q.Get("auth") != "" {
		q.Set("auth", "PRIVATE")
	}
	copy.RawQuery = q.Encode()
	return copy.String(), nil
}
