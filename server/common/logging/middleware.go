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
			sanitizedUrl := sanitize(r.URL)

			log.Infow("req:begin", "req", r.Method+" "+sanitizedUrl.String(), "from", from(r))

			cw := CaptureResponse(w)
			h.ServeHTTP(AllowWriteHeaderPrevention(cw), withCtx)

			elapsed := time.Since(started)

			log.Infow("req:done", "status", cw.StatusCode, "bytes", cw.ContentLength,
				"time", fmt.Sprintf("%vns", elapsed.Nanoseconds()),
				"time_human", elapsed.String(),
				"req", r.Method+" "+sanitizedUrl.String(),
				"from", from(r))
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

func sanitize(url *url.URL) *url.URL {
	q := url.Query()
	if q.Get("auth") != "" {
		q.Set("auth", "PRIVATE")
	}
	url.RawQuery = q.Encode()
	return url
}
