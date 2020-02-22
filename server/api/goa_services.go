package api

import (
	"context"
	"net/http"
	"os"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"

	"github.com/fieldkit/cloud/server/logging"

	testsvr "github.com/fieldkit/cloud/server/api/gen/http/test/server"
	test "github.com/fieldkit/cloud/server/api/gen/test"
)

func CreateGoaV3Handler(ctx context.Context, options *ControllerOptions) http.Handler {
	debug := false

	testSvc := NewTestSevice(ctx, options)
	testEndpoints := test.NewEndpoints(testSvc)

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	mux := goahttp.NewMuxer()

	eh := errorHandler()

	testServer := testsvr.New(testEndpoints, mux, dec, enc, eh, nil)
	if debug {
		servers := goahttp.Servers{
			testServer,
		}
		servers.Use(httpmdlwr.Debug(mux, os.Stdout))
	}

	testsvr.Mount(mux, testServer)

	loggingAdapter := logging.NewGoaMiddlewareAdapter(logging.Logger(ctx))

	var handler http.Handler = mux
	handler = httpmdlwr.Log(loggingAdapter)(handler)
	handler = httpmdlwr.RequestID()(handler)

	return handler
}

func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log := Logger(ctx).Sugar()
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		log.Errorw("fatal", "id", id, "message", err.Error())
	}
}
