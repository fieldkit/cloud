package backend

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/backend/ingestion/formatting"
	"github.com/fieldkit/cloud/server/logging"
)

const (
	FkDataBinaryContentType = "application/vnd.fk.data+binary"
	FkDataBase64ContentType = "application/vnd.fk.data+base64"
	ContentTypeHeaderName   = "Content-Type"
	ContentLengthHeaderName = "Content-Length"
	FkProcessingHeaderName  = "Fk-Processing"
	FkVersionHeaderName     = "Fk-Version"
	FkBuildHeaderName       = "Fk-Build"
	FkDeviceIdHeaderName    = "Fk-DeviceId"
	FkFileIdHeaderName      = "Fk-FileId"
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

func (si *StreamIngester) synchronous(ctx context.Context, headers *IncomingHeaders, w http.ResponseWriter, req *http.Request) {
	log := Logger(ctx).Sugar()

	startedAt := time.Now()
	status := http.StatusOK

	log.Infow("started", headers.ToLoggingFields()...)

	reader := &ReaderWrapper{
		BytesRead: 0,
		Target:    req.Body,
		Hash:      sha1.New(),
	}

	err := si.backend.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		saver := NewFormattedMessageSaver(si.backend)

		if headers.MediaType == FkDataBinaryContentType {
			binaryReader := NewFkBinaryReader(saver)
			if err := binaryReader.Read(txCtx, reader); err != nil {
				return err
			}
		} else if headers.MediaType == FkDataBase64ContentType {
			decoder := base64.NewDecoder(base64.StdEncoding, reader)
			binaryReader := NewFkBinaryReader(saver)
			if err := binaryReader.Read(txCtx, decoder); err != nil {
				return err
			}
		} else {
			decoder := json.NewDecoder(req.Body)
			message := &formatting.HttpJsonMessage{}
			err := decoder.Decode(message)
			if err != nil {
				return fmt.Errorf("JSON Error: '%v'", err)
			}

			messageId, err := uuid.NewRandom()
			if err != nil {
				return err
			}

			fm, err := message.ToFormattedMessage(ingestion.MessageId(messageId.String()))
			if err != nil {
				return err
			}

			recordChange, err := saver.HandleFormattedMessage(ctx, fm)
			if err != nil {
				return err
			}

			_ = recordChange
		}

		saver.EmitChanges(txCtx, si.sourceChanges)

		return nil
	})

	if err != nil {
		status = http.StatusInternalServerError
		log.Errorw("completed", "error", err, "bytes_read", reader.BytesRead, "hash", reader.Hash.Sum(nil), "time", time.Since(startedAt).String())
	} else {
		log.Infow("completed", "bytes_read", reader.BytesRead, "hash", reader.Hash.Sum(nil), "time", time.Since(startedAt).String())
	}

	w.WriteHeader(status)
}

func (si *StreamIngester) asynchronous(ctx context.Context, headers *IncomingHeaders, w http.ResponseWriter, req *http.Request) {
	log := Logger(ctx).Sugar()

	startedAt := time.Now()

	log.Infow("started (async)", headers.ToLoggingFields()...)

	if err := si.streamArchiver.Archive(ctx, headers.ContentType, req.Body); err != nil {
		log.Errorw("completed", "error", err, "time", time.Since(startedAt).String())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Infow("completed", "time", time.Since(startedAt).String())
		w.WriteHeader(http.StatusOK)
	}
}

var (
	ids = logging.NewIdGenerator()
)

type IncomingHeaders struct {
	ContentType   string
	ContentLength int32
	MediaType     string
	FkProcessing  string
	FkVersion     string
	FkBuild       string
	FkDeviceId    string
	FkFileId      string
}

func (h *IncomingHeaders) ToLoggingFields() []interface{} {
	keysAndValues := make([]interface{}, 0)

	keysAndValues = append(keysAndValues, "content_type", h.ContentType)
	keysAndValues = append(keysAndValues, "content_length", h.ContentLength)
	keysAndValues = append(keysAndValues, "firmware_version", h.FkVersion)
	keysAndValues = append(keysAndValues, "firmware_build", h.FkBuild)
	keysAndValues = append(keysAndValues, "device_id", h.FkDeviceId)
	keysAndValues = append(keysAndValues, "file_id", h.FkFileId)

	return keysAndValues
}

func NewIncomingHeaders(req *http.Request) (*IncomingHeaders, error) {
	contentLengthString := req.Header.Get(ContentLengthHeaderName)
	contentType := req.Header.Get(ContentTypeHeaderName)

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return nil, err
	}

	headers := &IncomingHeaders{
		ContentType:   contentType,
		ContentLength: int32(contentLength),
		MediaType:     mediaType,
		FkProcessing:  req.Header.Get(FkProcessingHeaderName),
		FkVersion:     req.Header.Get(FkVersionHeaderName),
		FkBuild:       req.Header.Get(FkBuildHeaderName),
		FkDeviceId:    req.Header.Get(FkDeviceIdHeaderName),
		FkFileId:      req.Header.Get(FkFileIdHeaderName),
	}

	return headers, nil
}

func (si *StreamIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := logging.WithNewTaskId(req.Context(), ids)
	log := Logger(ctx).Sugar()

	headers, err := NewIncomingHeaders(req)
	if err != nil {
		log.Errorw("failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if headers.MediaType != FkDataBinaryContentType && headers.MediaType != FkDataBase64ContentType && headers.MediaType != formatting.HttpProviderJsonContentType {
		log.Infow("Unknown content type", headers.ToLoggingFields()...)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if headers.FkProcessing == "" {
		si.synchronous(ctx, headers, w, req)
	} else if headers.FkProcessing == "async" {
		si.asynchronous(ctx, headers, w, req)
	}
}
