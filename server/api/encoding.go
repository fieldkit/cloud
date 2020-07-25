package api

import (
	"context"
	"fmt"
	"net/http"

	goahttp "goa.design/goa/v3/http"

	projectSvr "github.com/fieldkit/cloud/server/api/gen/http/project/server"
	userSvr "github.com/fieldkit/cloud/server/api/gen/http/user/server"

	"github.com/fieldkit/cloud/server/common/logging"
)

type StreamAndCacheFriendlyEncoder struct {
	ctx context.Context
	rw  http.ResponseWriter
}

func NewStreamAndCacheFriendlyEncoder(ctx context.Context, rw http.ResponseWriter) *StreamAndCacheFriendlyEncoder {
	return &StreamAndCacheFriendlyEncoder{
		ctx: ctx,
		rw:  rw,
	}
}

func (e *StreamAndCacheFriendlyEncoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case *projectSvr.DownloadPhotoResponseBody:
		e.rw.Header().Set("Content-Length", fmt.Sprintf("%v", v.Length))
		if v.ContentType != "" {
			e.rw.Header().Set("Content-Type", v.ContentType)
		}
		e.rw.Header().Set("ETag", fmt.Sprintf(`"%v"`, v.Etag))
		if len(v.Body) > 0 {
			e.rw.WriteHeader(http.StatusOK)
			e.rw.Write(v.Body)
		} else {
			e.rw.WriteHeader(http.StatusNotModified)
		}
		return nil
	case *userSvr.DownloadPhotoResponseBody:
		e.rw.Header().Set("Content-Length", fmt.Sprintf("%v", v.Length))
		if v.ContentType != "" {
			e.rw.Header().Set("Content-Type", v.ContentType)
		}
		e.rw.Header().Set("ETag", fmt.Sprintf(`"%v"`, v.Etag))
		if len(v.Body) > 0 {
			e.rw.WriteHeader(http.StatusOK)
			e.rw.Write(v.Body)
		} else {
			e.rw.WriteHeader(http.StatusNotModified)
		}
		return nil
	default:
		psr := e.rw.(*logging.PreventableStatusResponse)
		psr.FineWriteHeader()
		de := goahttp.ResponseEncoder(e.ctx, e.rw)
		return de.Encode(v)
	}
}

func InterceptDownloadResponses() func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
	return func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
		psr := w.(*logging.PreventableStatusResponse)
		psr.BufferNextWriteHeader()
		return NewStreamAndCacheFriendlyEncoder(ctx, w)
	}
}
