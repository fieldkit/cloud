package logging

import (
	"fmt"
	"net"
	"net/http"
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

			log.Infow("req:begin", "req", r.Method+" "+r.URL.String(), "from", from(r))

			cw := CaptureResponse(w)
			h.ServeHTTP(AllowWriteHeaderPrevention(cw), withCtx)

			elapsed := time.Since(started)

			log.Infow("req:done", "status", cw.StatusCode, "bytes", cw.ContentLength, "time", fmt.Sprintf("%vns", elapsed.Nanoseconds()), "time_human", elapsed.String(), "req", r.Method+" "+r.URL.String(), "from", from(r))
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
