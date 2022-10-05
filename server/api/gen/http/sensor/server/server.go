// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor HTTP server
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"context"
	"net/http"
	"regexp"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the sensor service endpoint HTTP handlers.
type Server struct {
	Mounts      []*MountPoint
	Meta        http.Handler
	StationMeta http.Handler
	SensorMeta  http.Handler
	Data        http.Handler
	Tail        http.Handler
	Recently    http.Handler
	Bookmark    http.Handler
	Resolve     http.Handler
	CORS        http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the sensor service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *sensor.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Meta", "GET", "/sensors"},
			{"StationMeta", "GET", "/meta/stations"},
			{"SensorMeta", "GET", "/meta/sensors"},
			{"Data", "GET", "/sensors/data"},
			{"Tail", "GET", "/sensors/data/tail"},
			{"Recently", "GET", "/sensors/data/recently"},
			{"Bookmark", "POST", "/bookmarks/save"},
			{"Resolve", "GET", "/bookmarks/resolve"},
			{"CORS", "OPTIONS", "/sensors"},
			{"CORS", "OPTIONS", "/meta/stations"},
			{"CORS", "OPTIONS", "/meta/sensors"},
			{"CORS", "OPTIONS", "/sensors/data"},
			{"CORS", "OPTIONS", "/sensors/data/tail"},
			{"CORS", "OPTIONS", "/sensors/data/recently"},
			{"CORS", "OPTIONS", "/bookmarks/save"},
			{"CORS", "OPTIONS", "/bookmarks/resolve"},
		},
		Meta:        NewMetaHandler(e.Meta, mux, decoder, encoder, errhandler, formatter),
		StationMeta: NewStationMetaHandler(e.StationMeta, mux, decoder, encoder, errhandler, formatter),
		SensorMeta:  NewSensorMetaHandler(e.SensorMeta, mux, decoder, encoder, errhandler, formatter),
		Data:        NewDataHandler(e.Data, mux, decoder, encoder, errhandler, formatter),
		Tail:        NewTailHandler(e.Tail, mux, decoder, encoder, errhandler, formatter),
		Recently:    NewRecentlyHandler(e.Recently, mux, decoder, encoder, errhandler, formatter),
		Bookmark:    NewBookmarkHandler(e.Bookmark, mux, decoder, encoder, errhandler, formatter),
		Resolve:     NewResolveHandler(e.Resolve, mux, decoder, encoder, errhandler, formatter),
		CORS:        NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "sensor" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Meta = m(s.Meta)
	s.StationMeta = m(s.StationMeta)
	s.SensorMeta = m(s.SensorMeta)
	s.Data = m(s.Data)
	s.Tail = m(s.Tail)
	s.Recently = m(s.Recently)
	s.Bookmark = m(s.Bookmark)
	s.Resolve = m(s.Resolve)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the sensor endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountMetaHandler(mux, h.Meta)
	MountStationMetaHandler(mux, h.StationMeta)
	MountSensorMetaHandler(mux, h.SensorMeta)
	MountDataHandler(mux, h.Data)
	MountTailHandler(mux, h.Tail)
	MountRecentlyHandler(mux, h.Recently)
	MountBookmarkHandler(mux, h.Bookmark)
	MountResolveHandler(mux, h.Resolve)
	MountCORSHandler(mux, h.CORS)
}

// MountMetaHandler configures the mux to serve the "sensor" service "meta"
// endpoint.
func MountMetaHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/sensors", f)
}

// NewMetaHandler creates a HTTP handler which loads the HTTP request and calls
// the "sensor" service "meta" endpoint.
func NewMetaHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeMetaResponse(encoder)
		encodeError    = EncodeMetaError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "meta")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountStationMetaHandler configures the mux to serve the "sensor" service
// "station meta" endpoint.
func MountStationMetaHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/meta/stations", f)
}

// NewStationMetaHandler creates a HTTP handler which loads the HTTP request
// and calls the "sensor" service "station meta" endpoint.
func NewStationMetaHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeStationMetaRequest(mux, decoder)
		encodeResponse = EncodeStationMetaResponse(encoder)
		encodeError    = EncodeStationMetaError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "station meta")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountSensorMetaHandler configures the mux to serve the "sensor" service
// "sensor meta" endpoint.
func MountSensorMetaHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/meta/sensors", f)
}

// NewSensorMetaHandler creates a HTTP handler which loads the HTTP request and
// calls the "sensor" service "sensor meta" endpoint.
func NewSensorMetaHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeSensorMetaResponse(encoder)
		encodeError    = EncodeSensorMetaError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "sensor meta")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountDataHandler configures the mux to serve the "sensor" service "data"
// endpoint.
func MountDataHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/sensors/data", f)
}

// NewDataHandler creates a HTTP handler which loads the HTTP request and calls
// the "sensor" service "data" endpoint.
func NewDataHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeDataRequest(mux, decoder)
		encodeResponse = EncodeDataResponse(encoder)
		encodeError    = EncodeDataError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "data")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountTailHandler configures the mux to serve the "sensor" service "tail"
// endpoint.
func MountTailHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/sensors/data/tail", f)
}

// NewTailHandler creates a HTTP handler which loads the HTTP request and calls
// the "sensor" service "tail" endpoint.
func NewTailHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeTailRequest(mux, decoder)
		encodeResponse = EncodeTailResponse(encoder)
		encodeError    = EncodeTailError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "tail")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountRecentlyHandler configures the mux to serve the "sensor" service
// "recently" endpoint.
func MountRecentlyHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/sensors/data/recently", f)
}

// NewRecentlyHandler creates a HTTP handler which loads the HTTP request and
// calls the "sensor" service "recently" endpoint.
func NewRecentlyHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeRecentlyRequest(mux, decoder)
		encodeResponse = EncodeRecentlyResponse(encoder)
		encodeError    = EncodeRecentlyError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "recently")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountBookmarkHandler configures the mux to serve the "sensor" service
// "bookmark" endpoint.
func MountBookmarkHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/bookmarks/save", f)
}

// NewBookmarkHandler creates a HTTP handler which loads the HTTP request and
// calls the "sensor" service "bookmark" endpoint.
func NewBookmarkHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeBookmarkRequest(mux, decoder)
		encodeResponse = EncodeBookmarkResponse(encoder)
		encodeError    = EncodeBookmarkError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "bookmark")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountResolveHandler configures the mux to serve the "sensor" service
// "resolve" endpoint.
func MountResolveHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleSensorOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/bookmarks/resolve", f)
}

// NewResolveHandler creates a HTTP handler which loads the HTTP request and
// calls the "sensor" service "resolve" endpoint.
func NewResolveHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeResolveRequest(mux, decoder)
		encodeResponse = EncodeResolveResponse(encoder)
		encodeError    = EncodeResolveError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "resolve")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sensor")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service sensor.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleSensorOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/sensors", f)
	mux.Handle("OPTIONS", "/meta/stations", f)
	mux.Handle("OPTIONS", "/meta/sensors", f)
	mux.Handle("OPTIONS", "/sensors/data", f)
	mux.Handle("OPTIONS", "/sensors/data/tail", f)
	mux.Handle("OPTIONS", "/sensors/data/recently", f)
	mux.Handle("OPTIONS", "/bookmarks/save", f)
	mux.Handle("OPTIONS", "/bookmarks/resolve", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleSensorOrigin applies the CORS response headers corresponding to the
// origin for the service sensor.
func handleSensorOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile("(.+[.])?fklocal.org:\\d+")
	spec1 := regexp.MustCompile("127.0.0.1:\\d+")
	spec2 := regexp.MustCompile("192.168.(\\d+).(\\d+):\\d+")
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec1) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec2) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "https://*.fieldkit.org") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "https://*.fkdev.org") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "https://dataviz.floodnet.nyc") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "https://fieldkit.org") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "https://fkdev.org") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE, PATCH, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
