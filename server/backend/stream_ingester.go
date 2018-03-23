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
	ContentTypeHeaderName   = "Content-Type"
	FkProcessingHeaderName  = "Fk-Processing"
)

type StreamIngester struct {
	backend        *Backend
	db             *sqlxcache.DB
	streamArchiver StreamArchiver
	sourceChanges  ingestion.SourceChangesPublisher
}

func NewStreamIngester(b *Backend, streamArchiver StreamArchiver, sourceChanges ingestion.SourceChangesPublisher) (si *StreamIngester, err error) {
	si = &StreamIngester{
		backend:        b,
		db:             b.db,
		streamArchiver: streamArchiver,
		sourceChanges:  sourceChanges,
	}

	return
}

func (si *StreamIngester) synchronous(w http.ResponseWriter, req *http.Request) {
	status := http.StatusOK

	log.Printf("Stream [%s]: begin", req.RemoteAddr)

	err := si.backend.db.WithNewTransaction(req.Context(), func(txCtx context.Context) error {
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

func (si *StreamIngester) asynchronous(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get(ContentTypeHeaderName)

	log.Printf("Stream [%s]: begin (async)", req.RemoteAddr)

	if err := si.streamArchiver.Archive(contentType, req.Body); err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (si *StreamIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get(ContentTypeHeaderName)
	fkProcessing := req.Header.Get(FkProcessingHeaderName)

	if contentType != FkDataBinaryContentType {
		log.Printf("Stream [%s]: unknown content type: %v", req.RemoteAddr, contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if fkProcessing == "" {
		si.synchronous(w, req)
	} else if fkProcessing == "async" {
		si.asynchronous(w, req)
	}
}
