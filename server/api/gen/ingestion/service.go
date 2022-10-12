// Code generated by goa v3.2.4, DO NOT EDIT.
//
// ingestion service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package ingestion

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the ingestion service interface.
type Service interface {
	// ProcessPending implements process pending.
	ProcessPending(context.Context, *ProcessPendingPayload) (err error)
	// WalkEverything implements walk everything.
	WalkEverything(context.Context, *WalkEverythingPayload) (err error)
	// ProcessStation implements process station.
	ProcessStation(context.Context, *ProcessStationPayload) (err error)
	// ProcessStationIngestions implements process station ingestions.
	ProcessStationIngestions(context.Context, *ProcessStationIngestionsPayload) (err error)
	// ProcessIngestion implements process ingestion.
	ProcessIngestion(context.Context, *ProcessIngestionPayload) (err error)
	// RefreshViews implements refresh views.
	RefreshViews(context.Context, *RefreshViewsPayload) (err error)
	// Delete implements delete.
	Delete(context.Context, *DeletePayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "ingestion"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [7]string{"process pending", "walk everything", "process station", "process station ingestions", "process ingestion", "refresh views", "delete"}

// ProcessPendingPayload is the payload type of the ingestion service process
// pending method.
type ProcessPendingPayload struct {
	Auth string
}

// WalkEverythingPayload is the payload type of the ingestion service walk
// everything method.
type WalkEverythingPayload struct {
	Auth string
}

// ProcessStationPayload is the payload type of the ingestion service process
// station method.
type ProcessStationPayload struct {
	Auth       string
	StationID  int32
	Completely *bool
	SkipManual *bool
}

// ProcessStationIngestionsPayload is the payload type of the ingestion service
// process station ingestions method.
type ProcessStationIngestionsPayload struct {
	Auth      string
	StationID int64
}

// ProcessIngestionPayload is the payload type of the ingestion service process
// ingestion method.
type ProcessIngestionPayload struct {
	Auth        string
	IngestionID int64
}

// RefreshViewsPayload is the payload type of the ingestion service refresh
// views method.
type RefreshViewsPayload struct {
	Auth string
}

// DeletePayload is the payload type of the ingestion service delete method.
type DeletePayload struct {
	Auth        string
	IngestionID int64
}

// MakeUnauthorized builds a goa.ServiceError from an error.
func MakeUnauthorized(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "unauthorized",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeForbidden builds a goa.ServiceError from an error.
func MakeForbidden(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "forbidden",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not-found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "bad-request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
