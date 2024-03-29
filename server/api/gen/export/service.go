// Code generated by goa v3.2.4, DO NOT EDIT.
//
// export service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package export

import (
	"context"
	"io"

	exportviews "github.com/fieldkit/cloud/server/api/gen/export/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the export service interface.
type Service interface {
	// ListMine implements list mine.
	ListMine(context.Context, *ListMinePayload) (res *UserExports, err error)
	// Status implements status.
	Status(context.Context, *StatusPayload) (res *ExportStatus, err error)
	// Download implements download.
	Download(context.Context, *DownloadPayload) (res *DownloadResult, body io.ReadCloser, err error)
	// Csv implements csv.
	Csv(context.Context, *CsvPayload) (res *CsvResult, err error)
	// JSONLines implements json lines.
	JSONLines(context.Context, *JSONLinesPayload) (res *JSONLinesResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "export"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [5]string{"list mine", "status", "download", "csv", "json lines"}

// ListMinePayload is the payload type of the export service list mine method.
type ListMinePayload struct {
	Auth string
}

// UserExports is the result type of the export service list mine method.
type UserExports struct {
	Exports []*ExportStatus
}

// StatusPayload is the payload type of the export service status method.
type StatusPayload struct {
	Auth string
	ID   string
}

// ExportStatus is the result type of the export service status method.
type ExportStatus struct {
	ID          int64
	Token       string
	CreatedAt   int64
	CompletedAt *int64
	Format      string
	Progress    float32
	Message     *string
	StatusURL   string
	DownloadURL *string
	Size        *int32
	Args        interface{}
}

// DownloadPayload is the payload type of the export service download method.
type DownloadPayload struct {
	ID   string
	Auth string
}

// DownloadResult is the result type of the export service download method.
type DownloadResult struct {
	Length             int64
	ContentType        string
	ContentDisposition string
}

// CsvPayload is the payload type of the export service csv method.
type CsvPayload struct {
	Auth       string
	Start      *int64
	End        *int64
	Stations   *string
	Sensors    *string
	Resolution *int32
	Aggregate  *string
	Complete   *bool
	Tail       *int32
}

// CsvResult is the result type of the export service csv method.
type CsvResult struct {
	Location string
}

// JSONLinesPayload is the payload type of the export service json lines method.
type JSONLinesPayload struct {
	Auth       string
	Start      *int64
	End        *int64
	Stations   *string
	Sensors    *string
	Resolution *int32
	Aggregate  *string
	Complete   *bool
	Tail       *int32
}

// JSONLinesResult is the result type of the export service json lines method.
type JSONLinesResult struct {
	Location string
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

// NewUserExports initializes result type UserExports from viewed result type
// UserExports.
func NewUserExports(vres *exportviews.UserExports) *UserExports {
	return newUserExports(vres.Projected)
}

// NewViewedUserExports initializes viewed result type UserExports from result
// type UserExports using the given view.
func NewViewedUserExports(res *UserExports, view string) *exportviews.UserExports {
	p := newUserExportsView(res)
	return &exportviews.UserExports{Projected: p, View: "default"}
}

// NewExportStatus initializes result type ExportStatus from viewed result type
// ExportStatus.
func NewExportStatus(vres *exportviews.ExportStatus) *ExportStatus {
	return newExportStatus(vres.Projected)
}

// NewViewedExportStatus initializes viewed result type ExportStatus from
// result type ExportStatus using the given view.
func NewViewedExportStatus(res *ExportStatus, view string) *exportviews.ExportStatus {
	p := newExportStatusView(res)
	return &exportviews.ExportStatus{Projected: p, View: "default"}
}

// newUserExports converts projected type UserExports to service type
// UserExports.
func newUserExports(vres *exportviews.UserExportsView) *UserExports {
	res := &UserExports{}
	if vres.Exports != nil {
		res.Exports = make([]*ExportStatus, len(vres.Exports))
		for i, val := range vres.Exports {
			res.Exports[i] = transformExportviewsExportStatusViewToExportStatus(val)
		}
	}
	return res
}

// newUserExportsView projects result type UserExports to projected type
// UserExportsView using the "default" view.
func newUserExportsView(res *UserExports) *exportviews.UserExportsView {
	vres := &exportviews.UserExportsView{}
	if res.Exports != nil {
		vres.Exports = make([]*exportviews.ExportStatusView, len(res.Exports))
		for i, val := range res.Exports {
			vres.Exports[i] = transformExportStatusToExportviewsExportStatusView(val)
		}
	}
	return vres
}

// newExportStatus converts projected type ExportStatus to service type
// ExportStatus.
func newExportStatus(vres *exportviews.ExportStatusView) *ExportStatus {
	res := &ExportStatus{
		CompletedAt: vres.CompletedAt,
		Message:     vres.Message,
		DownloadURL: vres.DownloadURL,
		Size:        vres.Size,
		Args:        vres.Args,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Token != nil {
		res.Token = *vres.Token
	}
	if vres.CreatedAt != nil {
		res.CreatedAt = *vres.CreatedAt
	}
	if vres.Format != nil {
		res.Format = *vres.Format
	}
	if vres.Progress != nil {
		res.Progress = *vres.Progress
	}
	if vres.StatusURL != nil {
		res.StatusURL = *vres.StatusURL
	}
	return res
}

// newExportStatusView projects result type ExportStatus to projected type
// ExportStatusView using the "default" view.
func newExportStatusView(res *ExportStatus) *exportviews.ExportStatusView {
	vres := &exportviews.ExportStatusView{
		ID:          &res.ID,
		Token:       &res.Token,
		CreatedAt:   &res.CreatedAt,
		CompletedAt: res.CompletedAt,
		Format:      &res.Format,
		Progress:    &res.Progress,
		Message:     res.Message,
		StatusURL:   &res.StatusURL,
		DownloadURL: res.DownloadURL,
		Size:        res.Size,
		Args:        res.Args,
	}
	return vres
}

// transformExportviewsExportStatusViewToExportStatus builds a value of type
// *ExportStatus from a value of type *exportviews.ExportStatusView.
func transformExportviewsExportStatusViewToExportStatus(v *exportviews.ExportStatusView) *ExportStatus {
	if v == nil {
		return nil
	}
	res := &ExportStatus{
		ID:          *v.ID,
		Token:       *v.Token,
		CreatedAt:   *v.CreatedAt,
		CompletedAt: v.CompletedAt,
		Format:      *v.Format,
		Progress:    *v.Progress,
		Message:     v.Message,
		StatusURL:   *v.StatusURL,
		DownloadURL: v.DownloadURL,
		Size:        v.Size,
		Args:        v.Args,
	}

	return res
}

// transformExportStatusToExportviewsExportStatusView builds a value of type
// *exportviews.ExportStatusView from a value of type *ExportStatus.
func transformExportStatusToExportviewsExportStatusView(v *ExportStatus) *exportviews.ExportStatusView {
	res := &exportviews.ExportStatusView{
		ID:          &v.ID,
		Token:       &v.Token,
		CreatedAt:   &v.CreatedAt,
		CompletedAt: v.CompletedAt,
		Format:      &v.Format,
		Progress:    &v.Progress,
		Message:     v.Message,
		StatusURL:   &v.StatusURL,
		DownloadURL: v.DownloadURL,
		Size:        v.Size,
		Args:        v.Args,
	}

	return res
}
