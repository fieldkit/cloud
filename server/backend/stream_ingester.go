package backend

import (
	"github.com/conservify/sqlxcache"
	pb "github.com/fieldkit/data-protocol"
	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"
	"log"
	"net/http"
)

const (
	FkDataBinaryContentType = "application/vnd.fk.data+binary"
)

type StreamIngester struct {
	backend *Backend
	db      *sqlxcache.DB
}

func NewStreamIngester(b *Backend) (si *StreamIngester, err error) {
	si = &StreamIngester{
		backend: b,
		db:      b.db,
	}

	return
}

func (si *StreamIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if contentType != FkDataBinaryContentType {
		log.Printf("Stream [%s]: unknown content type: %v", req.RemoteAddr, contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Stream [%s]: begin %v", req.RemoteAddr, contentType)

	binaryReader := NewFkBinaryReader(si.backend)

	unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var record pb.DataRecord
		err := proto.Unmarshal(b, &record)
		if err != nil {
			log.Printf("Error unmarshalling record: %v", err)
			return nil, err
		}

		_ = binaryReader.Push(&record)

		return &record, nil
	})

	_, err := stream.ReadLengthPrefixedCollection(req.Body, unmarshalFunc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Stream [%s]: ingesting error: %v", req.RemoteAddr, err)
		return
	}

	binaryReader.Done()

	w.WriteHeader(http.StatusOK)
}
