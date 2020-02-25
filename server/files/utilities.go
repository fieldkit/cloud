package files

import (
	"io"
)

type countingReader struct {
	target    io.Reader
	bytesRead int
}

func newCountingReader(target io.Reader) *countingReader {
	return &countingReader{
		target: target,
	}
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.target.Read(p)
	r.bytesRead += n
	return n, err
}
