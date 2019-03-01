package backend

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	_ "io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/backend/ingestion/formatting"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/logging"
)

const (
	FkDataBinaryContentType    = "application/vnd.fk.data+binary"
	FkDataBase64ContentType    = "application/vnd.fk.data+base64"
	MultiPartFormDataMediaType = "multipart/form-data"
	ContentTypeHeaderName      = "Content-Type"
	ContentLengthHeaderName    = "Content-Length"
	FkProcessingHeaderName     = "Fk-Processing"
	FkVersionHeaderName        = "Fk-Version"
	FkBuildHeaderName          = "Fk-Build"
	FkDeviceIdHeaderName       = "Fk-DeviceId"
	FkFileIdHeaderName         = "Fk-FileId"
	FkFileOffsetHeaderName     = "Fk-FileOffset"
	FkFileNameHeaderName       = "Fk-FileName"
	FkFileVersionHeaderName    = "Fk-FileVersion"
	FkCompiledHeaderName       = "Fk-Compiled"
	FkUploadNameHeaderName     = "Fk-UploadName"
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

func (si *StreamIngester) download(ctx context.Context, headers *IncomingHeaders, target io.Reader) error {
	log := Logger(ctx).Sugar()

	startedAt := time.Now()

	log.Infow("started", headers.ToLoggingFields()...)

	reader := &ReaderWrapper{
		BytesRead: 0,
		Target:    target,
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
			decoder := json.NewDecoder(target)
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
		log.Errorw("completed", "error", err, "bytes_read", reader.BytesRead, "hash", reader.Hash.Sum(nil), "time", time.Since(startedAt).String())
	} else {
		log.Infow("completed", "bytes_read", reader.BytesRead, "hash", reader.Hash.Sum(nil), "time", time.Since(startedAt).String())
	}

	return nil
}

func (si *StreamIngester) synchronous(ctx context.Context, headers *IncomingHeaders, w http.ResponseWriter, reader io.Reader) error {
	err := si.download(ctx, headers, reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (si *StreamIngester) saveStream(ctx context.Context, headers *IncomingHeaders, saved *SavedStream) error {
	metaData, err := json.Marshal(headers)
	if err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}

	stream := data.DeviceStream{
		Time:     time.Now(),
		StreamID: saved.ID,
		Firmware: headers.FkVersion,
		DeviceID: headers.FkDeviceId,
		Size:     int64(headers.ContentLength),
		FileID:   headers.FkFileId,
		URL:      saved.URL,
		Meta:     metaData,
	}

	if _, err := si.db.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_stream (time, stream_id, firmware, device_id, size, file_id, url, meta)
		   VALUES (:time, :stream_id, :firmware, :device_id, :size, :file_id, :url, :meta)
		   `, stream); err != nil {
		return err
	}

	return nil
}

func (si *StreamIngester) asynchronous(ctx context.Context, headers *IncomingHeaders, w http.ResponseWriter, reader io.Reader) error {
	log := Logger(ctx).Sugar()

	startedAt := time.Now()

	log.Infow("started (async)", headers.ToLoggingFields()...)

	if saved, err := si.streamArchiver.Archive(ctx, headers, reader); err != nil {
		log.Errorw("completed", "error", err, "time", time.Since(startedAt).String())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if saved != nil {
			err = si.saveStream(ctx, headers, saved)
			if err != nil {
				log.Errorw("completed", "error", err, "time", time.Since(startedAt).String())
			}

			log.Infow("completed (saved)", "stream_id", saved.ID, "time", time.Since(startedAt).String())
		} else {
			log.Infow("completed", "stream_id", saved.ID, "time", time.Since(startedAt).String())
		}
		w.WriteHeader(http.StatusOK)
	}

	return nil
}

var (
	ids = logging.NewIdGenerator()
)

type IncomingHeaders struct {
	ContentType     string
	ContentLength   int32
	MediaType       string
	MediaTypeParams map[string]string
	FkProcessing    string
	FkVersion       string
	FkBuild         string
	FkDeviceId      string
	FkFileId        string
	FkFileOffset    string
	FkFileVersion   string
	FkFileName      string
	FkCompiled      string
	FkUploadName    string
}

func (h *IncomingHeaders) ToLoggingFields() []interface{} {
	keysAndValues := make([]interface{}, 0)

	keysAndValues = append(keysAndValues, "content_type", h.ContentType)
	keysAndValues = append(keysAndValues, "content_length", h.ContentLength)
	keysAndValues = append(keysAndValues, "media_type", h.MediaType)
	keysAndValues = append(keysAndValues, "firmware_version", h.FkVersion)
	keysAndValues = append(keysAndValues, "firmware_build", h.FkBuild)
	keysAndValues = append(keysAndValues, "device_id", h.FkDeviceId)
	keysAndValues = append(keysAndValues, "file_id", h.FkFileId)
	keysAndValues = append(keysAndValues, "file_name", h.FkFileName)
	keysAndValues = append(keysAndValues, "file_offset", h.FkFileOffset)
	keysAndValues = append(keysAndValues, "file_version", h.FkFileVersion)
	keysAndValues = append(keysAndValues, "upload_name", h.FkUploadName)
	keysAndValues = append(keysAndValues, "compiled", h.FkCompiled)

	return keysAndValues
}

func (h *IncomingHeaders) ToPartHeaders(contentType string, contentLength int32) *IncomingHeaders {
	mediaType, mediaTypeParams, _ := mime.ParseMediaType(contentType)

	headers := &IncomingHeaders{
		ContentType:     contentType,
		ContentLength:   contentLength,
		MediaType:       mediaType,
		MediaTypeParams: mediaTypeParams,
		FkProcessing:    h.FkProcessing,
		FkVersion:       h.FkVersion,
		FkBuild:         h.FkBuild,
		FkDeviceId:      h.FkDeviceId,
		FkFileId:        h.FkFileId,
		FkFileOffset:    h.FkFileOffset,
		FkFileVersion:   h.FkFileVersion,
		FkFileName:      h.FkFileName,
		FkCompiled:      h.FkCompiled,
		FkUploadName:    h.FkUploadName,
	}

	return headers
}

func NewIncomingHeaders(req *http.Request) (*IncomingHeaders, error) {
	contentLengthString := req.Header.Get(ContentLengthHeaderName)
	contentType := req.Header.Get(ContentTypeHeaderName)

	mediaType, mediaTypeParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return nil, err
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("Invalid Content-Length (%v)", contentLength)
	}

	headers := &IncomingHeaders{
		ContentType:     contentType,
		ContentLength:   int32(contentLength),
		MediaType:       mediaType,
		MediaTypeParams: mediaTypeParams,
		FkProcessing:    req.Header.Get(FkProcessingHeaderName),
		FkVersion:       req.Header.Get(FkVersionHeaderName),
		FkBuild:         req.Header.Get(FkBuildHeaderName),
		FkDeviceId:      req.Header.Get(FkDeviceIdHeaderName),
		FkFileId:        req.Header.Get(FkFileIdHeaderName),
		FkFileOffset:    req.Header.Get(FkFileOffsetHeaderName),
		FkFileName:      req.Header.Get(FkFileNameHeaderName),
		FkFileVersion:   req.Header.Get(FkFileVersionHeaderName),
		FkCompiled:      req.Header.Get(FkCompiledHeaderName),
		FkUploadName:    req.Header.Get(FkUploadNameHeaderName),
	}

	if headers.FkDeviceId == "" {
		return nil, fmt.Errorf("Invalid %s ('%s')", FkDeviceIdHeaderName, headers.FkDeviceId)
	}

	return headers, nil
}

func acceptableMediaType(mediaType string) bool {
	if mediaType == FkDataBinaryContentType {
		return true
	}
	if mediaType == FkDataBase64ContentType {
		return true
	}
	if mediaType == formatting.HttpProviderJsonContentType {
		return true
	}
	if mediaType == MultiPartFormDataMediaType {
		return true
	}
	return false
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

	nestedCtx := logging.WithDeviceId(ctx, headers.FkDeviceId)

	if !acceptableMediaType(headers.MediaType) {
		log.Infow("Unknown Content-Type", headers.ToLoggingFields()...)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if headers.MediaType == MultiPartFormDataMediaType {
		log.Infow("Reading form", headers.ToLoggingFields()...)

		mr := multipart.NewReader(req.Body, headers.MediaTypeParams["boundary"])
		for {
			var reader io.Reader

			buffered := false
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Errorw("failed", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			contentDisposition := p.Header.Get("Content-Disposition")
			contentType := p.Header.Get("Content-Type")
			contentLengthString := p.Header.Get("Content-Length")

			contentLength, err := strconv.Atoi(contentLengthString)
			if err != nil || contentLength <= 0 {
				buf := &bytes.Buffer{}
				bytesCopied, err := io.Copy(buf, p)
				if err != nil {
					log.Errorw("failed", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				contentLength = int(bytesCopied)
				reader = buf
				buffered = true
			} else {
				reader = p
			}

			log.Infow("Part", "content_disposition", contentDisposition, "content_type", contentType, "content_length", contentLength, "buffered", buffered)

			partHeaders := headers.ToPartHeaders(contentType, int32(contentLength))
			if headers.FkProcessing == "sync" {
				si.synchronous(nestedCtx, partHeaders, w, reader)
			} else {
				si.asynchronous(nestedCtx, partHeaders, w, reader)
			}
		}
	} else {
		if headers.FkProcessing == "sync" {
			si.synchronous(nestedCtx, headers, w, req.Body)
		} else {
			si.asynchronous(nestedCtx, headers, w, req.Body)
		}
	}
}
