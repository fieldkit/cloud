package api

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/zap"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"

	"github.com/fieldkit/cloud/server/logging"

	testsvr "github.com/fieldkit/cloud/server/api/gen/http/test/server"
	test "github.com/fieldkit/cloud/server/api/gen/test"
)

type adapter struct {
	logger *zap.Logger
}

func newGoaAdapter(logger *zap.Logger) middleware.Logger {
	return &adapter{logger: logger.Named("goa")}
}

func (a *adapter) Log(keyvals ...interface{}) error {
	fields := logging.ToZapFields(keyvals)
	a.logger.Info("goa", *fields...)
	return nil
}

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

	loggingAdapter := newGoaAdapter(logging.Logger(nil))

	var handler http.Handler = mux
	handler = httpmdlwr.Log(loggingAdapter)(handler)
	handler = httpmdlwr.RequestID()(handler)

	return handler
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		// id := ctx.Value(middleware.RequestIDKey).(string)
		// w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		// logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}
