package backend

import (
	"context"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

const (
	FkDataBinaryContentType = "application/vnd.fk.data+binary"
)

type StreamIngester struct {
	backend       *Backend
	db            *sqlxcache.DB
	sourceChanges ingestion.SourceChangesPublisher
}

func NewStreamIngester(b *Backend, sourceChanges ingestion.SourceChangesPublisher) (si *StreamIngester, err error) {
	si = &StreamIngester{
		backend:       b,
		db:            b.db,
		sourceChanges: sourceChanges,
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

	status := http.StatusOK

	err := si.backend.db.WithNewTransaction(req.Context(), func(txCtx context.Context) error {
		log.Printf("Stream [%s]: begin %v", req.RemoteAddr, contentType)

		saver := NewFormattedMessageSaver(si.backend)
		binaryReader := NewFkBinaryReader(saver)

		unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
			var record pb.DataRecord
			err := proto.Unmarshal(b, &record)
			if err != nil {
				// We keep reading, this may just be a protocol version issue.
				log.Printf("Error unmarshalling record: %v", err)
				return nil, nil
			}

			err = binaryReader.Push(txCtx, &record)
			if err, ok := err.(*ingestion.IngestError); ok {
				if err.Critical {
					return nil, err
				} else {
					log.Printf("Error: %v", err)
				}
			} else if err != nil {
				return nil, err
			}

			return &record, nil
		})

		_, err := stream.ReadLengthPrefixedCollection(req.Body, unmarshalFunc)
		if err != nil {
			status = http.StatusInternalServerError
			log.Printf("Stream [%s]: error: %v", req.RemoteAddr, err)
			return nil
		}

		saver.EmitChanges(si.sourceChanges)

		binaryReader.Done()

		return nil
	})
	if err != nil {
		log.Printf("Error: %v", err)
	}

	w.WriteHeader(status)
}
