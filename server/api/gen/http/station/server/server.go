// Code generated by goa v3.2.4, DO NOT EDIT.
//
// station HTTP server
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"context"
	"net/http"
	"regexp"

	station "github.com/fieldkit/cloud/server/api/gen/station"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the station service endpoint HTTP handlers.
type Server struct {
	Mounts                []*MountPoint
	Add                   http.Handler
	Get                   http.Handler
	Transfer              http.Handler
	DefaultPhoto          http.Handler
	Update                http.Handler
	ListMine              http.Handler
	ListProject           http.Handler
	ListAssociated        http.Handler
	ListProjectAssociated http.Handler
	DownloadPhoto         http.Handler
	ListAll               http.Handler
	Delete                http.Handler
	AdminSearch           http.Handler
	Progress              http.Handler
	UpdateModule          http.Handler
	CORS                  http.Handler
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

// New instantiates HTTP handlers for all the station service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *station.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Add", "POST", "/stations"},
			{"Get", "GET", "/stations/{id}"},
			{"Transfer", "POST", "/stations/{id}/transfer/{ownerId}"},
			{"DefaultPhoto", "POST", "/stations/{id}/photo/{photoId}"},
			{"Update", "PATCH", "/stations/{id}"},
			{"ListMine", "GET", "/user/stations"},
			{"ListProject", "GET", "/projects/{id}/stations"},
			{"ListAssociated", "GET", "/stations/{id}/associated"},
			{"ListProjectAssociated", "GET", "/projects/{projectId}/associated"},
			{"DownloadPhoto", "GET", "/stations/{stationId}/photo"},
			{"ListAll", "GET", "/admin/stations"},
			{"Delete", "DELETE", "/admin/stations/{stationId}"},
			{"AdminSearch", "POST", "/admin/stations/search"},
			{"Progress", "GET", "/stations/{stationId}/progress"},
			{"UpdateModule", "PATCH", "/stations/{id}/modules/{moduleId}"},
			{"CORS", "OPTIONS", "/stations"},
			{"CORS", "OPTIONS", "/stations/{id}"},
			{"CORS", "OPTIONS", "/stations/{id}/transfer/{ownerId}"},
			{"CORS", "OPTIONS", "/stations/{id}/photo/{photoId}"},
			{"CORS", "OPTIONS", "/user/stations"},
			{"CORS", "OPTIONS", "/projects/{id}/stations"},
			{"CORS", "OPTIONS", "/stations/{id}/associated"},
			{"CORS", "OPTIONS", "/projects/{projectId}/associated"},
			{"CORS", "OPTIONS", "/stations/{stationId}/photo"},
			{"CORS", "OPTIONS", "/admin/stations"},
			{"CORS", "OPTIONS", "/admin/stations/{stationId}"},
			{"CORS", "OPTIONS", "/admin/stations/search"},
			{"CORS", "OPTIONS", "/stations/{stationId}/progress"},
			{"CORS", "OPTIONS", "/stations/{id}/modules/{moduleId}"},
		},
		Add:                   NewAddHandler(e.Add, mux, decoder, encoder, errhandler, formatter),
		Get:                   NewGetHandler(e.Get, mux, decoder, encoder, errhandler, formatter),
		Transfer:              NewTransferHandler(e.Transfer, mux, decoder, encoder, errhandler, formatter),
		DefaultPhoto:          NewDefaultPhotoHandler(e.DefaultPhoto, mux, decoder, encoder, errhandler, formatter),
		Update:                NewUpdateHandler(e.Update, mux, decoder, encoder, errhandler, formatter),
		ListMine:              NewListMineHandler(e.ListMine, mux, decoder, encoder, errhandler, formatter),
		ListProject:           NewListProjectHandler(e.ListProject, mux, decoder, encoder, errhandler, formatter),
		ListAssociated:        NewListAssociatedHandler(e.ListAssociated, mux, decoder, encoder, errhandler, formatter),
		ListProjectAssociated: NewListProjectAssociatedHandler(e.ListProjectAssociated, mux, decoder, encoder, errhandler, formatter),
		DownloadPhoto:         NewDownloadPhotoHandler(e.DownloadPhoto, mux, decoder, encoder, errhandler, formatter),
		ListAll:               NewListAllHandler(e.ListAll, mux, decoder, encoder, errhandler, formatter),
		Delete:                NewDeleteHandler(e.Delete, mux, decoder, encoder, errhandler, formatter),
		AdminSearch:           NewAdminSearchHandler(e.AdminSearch, mux, decoder, encoder, errhandler, formatter),
		Progress:              NewProgressHandler(e.Progress, mux, decoder, encoder, errhandler, formatter),
		UpdateModule:          NewUpdateModuleHandler(e.UpdateModule, mux, decoder, encoder, errhandler, formatter),
		CORS:                  NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "station" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Add = m(s.Add)
	s.Get = m(s.Get)
	s.Transfer = m(s.Transfer)
	s.DefaultPhoto = m(s.DefaultPhoto)
	s.Update = m(s.Update)
	s.ListMine = m(s.ListMine)
	s.ListProject = m(s.ListProject)
	s.ListAssociated = m(s.ListAssociated)
	s.ListProjectAssociated = m(s.ListProjectAssociated)
	s.DownloadPhoto = m(s.DownloadPhoto)
	s.ListAll = m(s.ListAll)
	s.Delete = m(s.Delete)
	s.AdminSearch = m(s.AdminSearch)
	s.Progress = m(s.Progress)
	s.UpdateModule = m(s.UpdateModule)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the station endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountAddHandler(mux, h.Add)
	MountGetHandler(mux, h.Get)
	MountTransferHandler(mux, h.Transfer)
	MountDefaultPhotoHandler(mux, h.DefaultPhoto)
	MountUpdateHandler(mux, h.Update)
	MountListMineHandler(mux, h.ListMine)
	MountListProjectHandler(mux, h.ListProject)
	MountListAssociatedHandler(mux, h.ListAssociated)
	MountListProjectAssociatedHandler(mux, h.ListProjectAssociated)
	MountDownloadPhotoHandler(mux, h.DownloadPhoto)
	MountListAllHandler(mux, h.ListAll)
	MountDeleteHandler(mux, h.Delete)
	MountAdminSearchHandler(mux, h.AdminSearch)
	MountProgressHandler(mux, h.Progress)
	MountUpdateModuleHandler(mux, h.UpdateModule)
	MountCORSHandler(mux, h.CORS)
}

// MountAddHandler configures the mux to serve the "station" service "add"
// endpoint.
func MountAddHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/stations", f)
}

// NewAddHandler creates a HTTP handler which loads the HTTP request and calls
// the "station" service "add" endpoint.
func NewAddHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeAddRequest(mux, decoder)
		encodeResponse = EncodeAddResponse(encoder)
		encodeError    = EncodeAddError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "add")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountGetHandler configures the mux to serve the "station" service "get"
// endpoint.
func MountGetHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/stations/{id}", f)
}

// NewGetHandler creates a HTTP handler which loads the HTTP request and calls
// the "station" service "get" endpoint.
func NewGetHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetRequest(mux, decoder)
		encodeResponse = EncodeGetResponse(encoder)
		encodeError    = EncodeGetError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountTransferHandler configures the mux to serve the "station" service
// "transfer" endpoint.
func MountTransferHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/stations/{id}/transfer/{ownerId}", f)
}

// NewTransferHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "transfer" endpoint.
func NewTransferHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeTransferRequest(mux, decoder)
		encodeResponse = EncodeTransferResponse(encoder)
		encodeError    = EncodeTransferError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "transfer")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountDefaultPhotoHandler configures the mux to serve the "station" service
// "default photo" endpoint.
func MountDefaultPhotoHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/stations/{id}/photo/{photoId}", f)
}

// NewDefaultPhotoHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "default photo" endpoint.
func NewDefaultPhotoHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeDefaultPhotoRequest(mux, decoder)
		encodeResponse = EncodeDefaultPhotoResponse(encoder)
		encodeError    = EncodeDefaultPhotoError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "default photo")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountUpdateHandler configures the mux to serve the "station" service
// "update" endpoint.
func MountUpdateHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("PATCH", "/stations/{id}", f)
}

// NewUpdateHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "update" endpoint.
func NewUpdateHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUpdateRequest(mux, decoder)
		encodeResponse = EncodeUpdateResponse(encoder)
		encodeError    = EncodeUpdateError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "update")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountListMineHandler configures the mux to serve the "station" service "list
// mine" endpoint.
func MountListMineHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/user/stations", f)
}

// NewListMineHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "list mine" endpoint.
func NewListMineHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListMineRequest(mux, decoder)
		encodeResponse = EncodeListMineResponse(encoder)
		encodeError    = EncodeListMineError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list mine")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountListProjectHandler configures the mux to serve the "station" service
// "list project" endpoint.
func MountListProjectHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/projects/{id}/stations", f)
}

// NewListProjectHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "list project" endpoint.
func NewListProjectHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListProjectRequest(mux, decoder)
		encodeResponse = EncodeListProjectResponse(encoder)
		encodeError    = EncodeListProjectError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list project")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountListAssociatedHandler configures the mux to serve the "station" service
// "list associated" endpoint.
func MountListAssociatedHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/stations/{id}/associated", f)
}

// NewListAssociatedHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "list associated" endpoint.
func NewListAssociatedHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListAssociatedRequest(mux, decoder)
		encodeResponse = EncodeListAssociatedResponse(encoder)
		encodeError    = EncodeListAssociatedError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list associated")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountListProjectAssociatedHandler configures the mux to serve the "station"
// service "list project associated" endpoint.
func MountListProjectAssociatedHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/projects/{projectId}/associated", f)
}

// NewListProjectAssociatedHandler creates a HTTP handler which loads the HTTP
// request and calls the "station" service "list project associated" endpoint.
func NewListProjectAssociatedHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListProjectAssociatedRequest(mux, decoder)
		encodeResponse = EncodeListProjectAssociatedResponse(encoder)
		encodeError    = EncodeListProjectAssociatedError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list project associated")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountDownloadPhotoHandler configures the mux to serve the "station" service
// "download photo" endpoint.
func MountDownloadPhotoHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/stations/{stationId}/photo", f)
}

// NewDownloadPhotoHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "download photo" endpoint.
func NewDownloadPhotoHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeDownloadPhotoRequest(mux, decoder)
		encodeResponse = EncodeDownloadPhotoResponse(encoder)
		encodeError    = EncodeDownloadPhotoError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "download photo")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountListAllHandler configures the mux to serve the "station" service "list
// all" endpoint.
func MountListAllHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/admin/stations", f)
}

// NewListAllHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "list all" endpoint.
func NewListAllHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListAllRequest(mux, decoder)
		encodeResponse = EncodeListAllResponse(encoder)
		encodeError    = EncodeListAllError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list all")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountDeleteHandler configures the mux to serve the "station" service
// "delete" endpoint.
func MountDeleteHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("DELETE", "/admin/stations/{stationId}", f)
}

// NewDeleteHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "delete" endpoint.
func NewDeleteHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeDeleteRequest(mux, decoder)
		encodeResponse = EncodeDeleteResponse(encoder)
		encodeError    = EncodeDeleteError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "delete")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountAdminSearchHandler configures the mux to serve the "station" service
// "admin search" endpoint.
func MountAdminSearchHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/admin/stations/search", f)
}

// NewAdminSearchHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "admin search" endpoint.
func NewAdminSearchHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeAdminSearchRequest(mux, decoder)
		encodeResponse = EncodeAdminSearchResponse(encoder)
		encodeError    = EncodeAdminSearchError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "admin search")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountProgressHandler configures the mux to serve the "station" service
// "progress" endpoint.
func MountProgressHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/stations/{stationId}/progress", f)
}

// NewProgressHandler creates a HTTP handler which loads the HTTP request and
// calls the "station" service "progress" endpoint.
func NewProgressHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeProgressRequest(mux, decoder)
		encodeResponse = EncodeProgressResponse(encoder)
		encodeError    = EncodeProgressError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "progress")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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

// MountUpdateModuleHandler configures the mux to serve the "station" service
// "update module" endpoint.
func MountUpdateModuleHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleStationOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("PATCH", "/stations/{id}/modules/{moduleId}", f)
}

// NewUpdateModuleHandler creates a HTTP handler which loads the HTTP request
// and calls the "station" service "update module" endpoint.
func NewUpdateModuleHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUpdateModuleRequest(mux, decoder)
		encodeResponse = EncodeUpdateModuleResponse(encoder)
		encodeError    = EncodeUpdateModuleError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "update module")
		ctx = context.WithValue(ctx, goa.ServiceKey, "station")
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
// service station.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleStationOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/stations", f)
	mux.Handle("OPTIONS", "/stations/{id}", f)
	mux.Handle("OPTIONS", "/stations/{id}/transfer/{ownerId}", f)
	mux.Handle("OPTIONS", "/stations/{id}/photo/{photoId}", f)
	mux.Handle("OPTIONS", "/user/stations", f)
	mux.Handle("OPTIONS", "/projects/{id}/stations", f)
	mux.Handle("OPTIONS", "/stations/{id}/associated", f)
	mux.Handle("OPTIONS", "/projects/{projectId}/associated", f)
	mux.Handle("OPTIONS", "/stations/{stationId}/photo", f)
	mux.Handle("OPTIONS", "/admin/stations", f)
	mux.Handle("OPTIONS", "/admin/stations/{stationId}", f)
	mux.Handle("OPTIONS", "/admin/stations/search", f)
	mux.Handle("OPTIONS", "/stations/{stationId}/progress", f)
	mux.Handle("OPTIONS", "/stations/{id}/modules/{moduleId}", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleStationOrigin applies the CORS response headers corresponding to the
// origin for the service station.
func handleStationOrigin(h http.Handler) http.Handler {
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
