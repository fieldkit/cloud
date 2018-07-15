package tools

import (
	"fmt"
	"io"

	"github.com/robinpowered/go-proto/collection"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"
)

type bufferedStartReader struct {
	maximum   int
	buffer    []byte
	offset    int
	available int
	target    io.Reader
}

func newBufferedStartReader(t io.Reader, maximum int) *bufferedStartReader {
	return &bufferedStartReader{
		maximum:   maximum,
		offset:    0,
		available: 0,
		target:    t,
	}
}

func (r *bufferedStartReader) Read(p []byte) (n int, err error) {
	if r.buffer == nil {
		r.buffer = make([]byte, r.maximum)

		bytes, err := r.target.Read(r.buffer)
		if err != nil {
			return 0, err
		}

		r.available = bytes
		r.offset = 0
	}

	if r.offset < r.available {
		n = copy(p, r.buffer[r.offset:len(r.buffer)])
		r.offset += n
		return n, nil
	}

	return r.target.Read(p)
}

func (r *bufferedStartReader) Seek(n int) error {
	if r.available > 0 && r.offset == r.available {
		return fmt.Errorf("Buffer empty, seeking disallowed")
	}
	r.offset = n
	return nil
}

func ReadLengthPrefixedCollectionIgnoringIncompleteBeginning(r io.Reader, maximum int, f message.UnmarshalFunc) (pbs collection.MessageCollection, junk int, err error) {
	br := newBufferedStartReader(r, maximum)

	for offset := 0; offset < maximum; offset += 1 {
		err := br.Seek(offset)
		if err != nil {
			return nil, 0, err
		}

		items, err := stream.ReadLengthPrefixedCollection(br, f)
		if err != nil {
			continue
		}

		if items != nil {
			return items, offset, nil
		}
	}

	return nil, 0, fmt.Errorf("Too much garbage at beginning of length prefixed stream")
}
