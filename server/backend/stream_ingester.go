package backend

import (
	"context"
	"log"
	"net/http"

	"github.com/conservify/sqlxcache"

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

		if err := binaryReader.Read(txCtx, req.Body); err != nil {
			status = http.StatusInternalServerError
			log.Printf("Stream [%s]: error: %v", req.RemoteAddr, err)
			return nil
		}

		saver.EmitChanges(si.sourceChanges)

		return nil
	})

	if err != nil {
		status = http.StatusInternalServerError
		log.Printf("Error: %v", err)
	}

	w.WriteHeader(status)
}
