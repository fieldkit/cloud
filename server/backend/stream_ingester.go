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

func NewStreamIngester(b *Backend) (rmi *StreamIngester, err error) {
	rmi = &StreamIngester{
		backend: b,
		db:      b.db,
	}

	return
}

func (rmi *StreamIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType != FkDataBinaryContentType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var record pb.DataRecord
		err := proto.Unmarshal(b, &record)
		if err != nil {
			log.Printf("Error unmarshalling record: %v", err)
			return nil, err
		}

		log.Printf("%v", record)

		return &record, nil
	})

	_, err := stream.ReadLengthPrefixedCollection(req.Body, unmarshalFunc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
