package logging

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

var (
	ids = NewIdGenerator()
)

func Monitoring(m *Metrics) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return m.GatherMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()

			newCtx := WithNewTaskID(r.Context(), ids)
			withCtx := r.WithContext(newCtx)

			log := Logger(newCtx).Named("http").Sugar()

			log.Infow("started", "req", r.Method+" "+r.URL.String(), "from", from(r))

			rw := CaptureResponse(w)
			h.ServeHTTP(rw, withCtx)

			elapsed := time.Since(started)

			log.Infow("done", "status", rw.StatusCode, "bytes", rw.ContentLength, "time", fmt.Sprintf("%vns", elapsed.Nanoseconds()), "time_human", elapsed.String())
		}))
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

// ResponseCapture is a http.ResponseWriter which captures the response status
// code and content length.
type ResponseCapture struct {
	http.ResponseWriter
	StatusCode    int
	ContentLength int
}

func CaptureResponse(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{ResponseWriter: w}
}

func (w *ResponseCapture) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseCapture) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.ContentLength += n
	return n, err
}

func (w *ResponseCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("response writer does not support hijacking: %T", w.ResponseWriter)
}
