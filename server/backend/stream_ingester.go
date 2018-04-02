package backend

import (
	"context"
	"crypto/sha1"
	"hash"
	"io"
	"log"
	"net/http"

	"github.com/conservify/sqlxcache"
	"github.com/google/uuid"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

const (
	FkDataBinaryContentType = "application/vnd.fk.data+binary"
	ContentTypeHeaderName   = "Content-Type"
	ContentLengthHeaderName = "Content-Length"
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

type ReaderWrapper struct {
	BytesRead int64
	Target    io.Reader
	Hash      hash.Hash
}

func (rw *ReaderWrapper) Read(p []byte) (n int, err error) {
	n, err = rw.Target.Read(p)
	rw.BytesRead += int64(n)
	sliced := p[:n]
	rw.Hash.Write(sliced)
	return n, err
}

func (si *StreamIngester) synchronous(w http.ResponseWriter, req *http.Request, id *uuid.UUID) {
	contentLength := req.Header.Get(ContentLengthHeaderName)

	status := http.StatusOK

	log.Printf("Stream [%s]: begin (%v)", id, contentLength)

	reader := &ReaderWrapper{
		BytesRead: 0,
		Target:    req.Body,
		Hash:      sha1.New(),
	}

	err := si.backend.db.WithNewTransaction(req.Context(), func(txCtx context.Context) error {
		saver := NewFormattedMessageSaver(si.backend)
		binaryReader := NewFkBinaryReader(saver)

		if err := binaryReader.Read(txCtx, reader); err != nil {
			status = http.StatusInternalServerError
			log.Printf("Stream [%s]: error: %v (%d bytes)", id, err, reader.BytesRead)
			return nil
		}

		saver.EmitChanges(si.sourceChanges)

		return nil
	})

	if err != nil {
		status = http.StatusInternalServerError
		log.Printf("Stream [%s]: error: %v (%d bytes) (%x)", id, err, reader.BytesRead, reader.Hash.Sum(nil))
	} else {
		log.Printf("Stream [%s]: done (%d bytes) (%x)", id, reader.BytesRead, reader.Hash.Sum(nil))
	}

	w.WriteHeader(status)
}

func (si *StreamIngester) asynchronous(w http.ResponseWriter, req *http.Request, id *uuid.UUID) {
	contentType := req.Header.Get(ContentTypeHeaderName)
	contentLength := req.Header.Get(ContentLengthHeaderName)

	log.Printf("Stream [%s]: begin (%v) (async)", id, contentLength)

	if err := si.streamArchiver.Archive(contentType, req.Body); err != nil {
		log.Printf("Stream [%s]: error (%v) (async)", id, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (si *StreamIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get(ContentTypeHeaderName)
	fkProcessing := req.Header.Get(FkProcessingHeaderName)

	id, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if contentType != FkDataBinaryContentType {
		log.Printf("Stream [%v]: unknown content type: %v", id, contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if fkProcessing == "" {
		si.synchronous(w, req, &id)
	} else if fkProcessing == "async" {
		si.asynchronous(w, req, &id)
	}
}
