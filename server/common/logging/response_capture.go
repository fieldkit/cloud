package logging

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

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

// This is a http.ResponseWriter which allows us to circumvent the
// WriteHeader call that goa does before calling our encoder. Look,
// this sucks. I know. I just can't for the life of me figure out how
// to get Etagged response working for these images in another way.
type PreventableStatusResponse struct {
	http.ResponseWriter
	buffering bool
	buffered  int
}

func AllowWriteHeaderPrevention(w http.ResponseWriter) *PreventableStatusResponse {
	return &PreventableStatusResponse{ResponseWriter: w}
}

func (w *PreventableStatusResponse) BufferNextWriteHeader() {
	w.buffering = true
}

func (w *PreventableStatusResponse) FineWriteHeader() {
	w.ResponseWriter.WriteHeader(w.buffered)
}

func (w *PreventableStatusResponse) WriteHeader(code int) {
	if w.buffering {
		w.buffering = false
		w.buffered = code
		return
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *PreventableStatusResponse) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	return n, err
}

func (w *PreventableStatusResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("response writer does not support hijacking: %T", w.ResponseWriter)
}
