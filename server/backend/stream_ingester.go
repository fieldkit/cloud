package backend

import (
	"context"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/stream"

	"github.com/conservify/sqlxcache"
	pb "github.com/fieldkit/data-protocol"
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

	si.backend.db.WithNewTransaction(context.TODO(), func(txCtx context.Context, tx *sqlx.Tx) error {
		log.Printf("Stream [%s]: begin %v", req.RemoteAddr, contentType)

		binaryReader := NewFkBinaryReader(si.backend)

		unmarshalFunc := message.UnmarshalFunc(func(b []byte) (proto.Message, error) {
			var record pb.DataRecord
			err := proto.Unmarshal(b, &record)
			if err != nil {
				log.Printf("Error unmarshalling record: %v", err)
				return nil, err
			}

			_ = binaryReader.Push(txCtx, &record)

			return &record, nil
		})

		_, err := stream.ReadLengthPrefixedCollection(req.Body, unmarshalFunc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Stream [%s]: ingesting error: %v", req.RemoteAddr, err)
			return nil
		}

		binaryReader.Done(txCtx)

		return nil
	})

	w.WriteHeader(http.StatusOK)
}
