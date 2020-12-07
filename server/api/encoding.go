package api

import (
	"context"
	"fmt"
	"net/http"

	goahttp "goa.design/goa/v3/http"

	projectSvr "github.com/fieldkit/cloud/server/api/gen/http/project/server"
	stationSvr "github.com/fieldkit/cloud/server/api/gen/http/station/server"
	userSvr "github.com/fieldkit/cloud/server/api/gen/http/user/server"

	"github.com/fieldkit/cloud/server/common/logging"
)

type StreamAndCacheFriendlyEncoder struct {
	ctx context.Context
	rw  http.ResponseWriter
	de  func(context.Context, http.ResponseWriter) goahttp.Encoder
}

func NewStreamAndCacheFriendlyEncoder(ctx context.Context, defaultEncoder func(context.Context, http.ResponseWriter) goahttp.Encoder, rw http.ResponseWriter) *StreamAndCacheFriendlyEncoder {
	return &StreamAndCacheFriendlyEncoder{
		ctx: ctx,
		de:  defaultEncoder,
		rw:  rw,
	}
}

type DownloadResponseBody struct {
	Length      int64
	ContentType string
	Body        []byte
	Etag        string
}

func (e *StreamAndCacheFriendlyEncoder) EncodeDownload(body *DownloadResponseBody) error {
	e.rw.Header().Set("Content-Length", fmt.Sprintf("%v", body.Length))
	if body.ContentType != "" {
		e.rw.Header().Set("Content-Type", body.ContentType)
	}
	e.rw.Header().Set("ETag", fmt.Sprintf(`"%v"`, body.Etag))
	if len(body.Body) > 0 {
		e.rw.WriteHeader(http.StatusOK)
		e.rw.Write(body.Body)
	} else {
		e.rw.WriteHeader(http.StatusNotModified)
	}
	return nil
}

func (e *StreamAndCacheFriendlyEncoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case *projectSvr.DownloadPhotoResponseBody:
		return e.EncodeDownload(&DownloadResponseBody{
			Length:      v.Length,
			ContentType: v.ContentType,
			Etag:        v.Etag,
			Body:        v.Body,
		})
	case *userSvr.DownloadPhotoResponseBody:
		return e.EncodeDownload(&DownloadResponseBody{
			Length:      v.Length,
			ContentType: v.ContentType,
			Etag:        v.Etag,
			Body:        v.Body,
		})
	case *stationSvr.DownloadPhotoResponseBody:
		return e.EncodeDownload(&DownloadResponseBody{
			Length:      v.Length,
			ContentType: v.ContentType,
			Etag:        v.Etag,
			Body:        v.Body,
		})
	default:
		psr := e.rw.(*logging.PreventableStatusResponse)
		psr.FineWriteHeader()
		de := e.de(e.ctx, e.rw)
		return de.Encode(v)
	}
}

func InterceptDownloadResponses(defaultEncoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
	return func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
		if psr, ok := w.(*logging.PreventableStatusResponse); ok {
			psr.BufferNextWriteHeader()
			return NewStreamAndCacheFriendlyEncoder(ctx, defaultEncoder, w)
		}
		return defaultEncoder(ctx, w)
	}
}
