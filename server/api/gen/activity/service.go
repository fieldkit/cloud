// Code generated by goa v3.2.4, DO NOT EDIT.
//
// activity service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package activity

import (
	"context"

	activityviews "github.com/fieldkit/cloud/server/api/gen/activity/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the activity service interface.
type Service interface {
	// Station implements station.
	Station(context.Context, *StationPayload) (res *StationActivityPage, err error)
	// Project implements project.
	Project(context.Context, *ProjectPayload) (res *ProjectActivityPage, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "activity"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"station", "project"}

// StationPayload is the payload type of the activity service station method.
type StationPayload struct {
	Auth *string
	ID   int64
	Page *int64
}

// StationActivityPage is the result type of the activity service station
// method.
type StationActivityPage struct {
	Activities ActivityEntryCollection
	Total      int32
	Page       int32
}

// ProjectPayload is the payload type of the activity service project method.
type ProjectPayload struct {
	Auth *string
	ID   int64
	Page *int64
}

// ProjectActivityPage is the result type of the activity service project
// method.
type ProjectActivityPage struct {
	Activities ActivityEntryCollection
	Total      int32
	Page       int32
}

type ActivityEntryCollection []*ActivityEntry

type ActivityEntry struct {
	ID        int64
	Key       string
	Project   *ProjectSummary
	Station   *StationSummary
	CreatedAt int64
	Type      string
	Meta      interface{}
}

type ProjectSummary struct {
	ID   int64
	Name string
}

type StationSummary struct {
	ID   int64
	Name string
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

// NewStationActivityPage initializes result type StationActivityPage from
// viewed result type StationActivityPage.
func NewStationActivityPage(vres *activityviews.StationActivityPage) *StationActivityPage {
	return newStationActivityPage(vres.Projected)
}

// NewViewedStationActivityPage initializes viewed result type
// StationActivityPage from result type StationActivityPage using the given
// view.
func NewViewedStationActivityPage(res *StationActivityPage, view string) *activityviews.StationActivityPage {
	p := newStationActivityPageView(res)
	return &activityviews.StationActivityPage{Projected: p, View: "default"}
}

// NewProjectActivityPage initializes result type ProjectActivityPage from
// viewed result type ProjectActivityPage.
func NewProjectActivityPage(vres *activityviews.ProjectActivityPage) *ProjectActivityPage {
	return newProjectActivityPage(vres.Projected)
}

// NewViewedProjectActivityPage initializes viewed result type
// ProjectActivityPage from result type ProjectActivityPage using the given
// view.
func NewViewedProjectActivityPage(res *ProjectActivityPage, view string) *activityviews.ProjectActivityPage {
	p := newProjectActivityPageView(res)
	return &activityviews.ProjectActivityPage{Projected: p, View: "default"}
}

// newStationActivityPage converts projected type StationActivityPage to
// service type StationActivityPage.
func newStationActivityPage(vres *activityviews.StationActivityPageView) *StationActivityPage {
	res := &StationActivityPage{}
	if vres.Total != nil {
		res.Total = *vres.Total
	}
	if vres.Page != nil {
		res.Page = *vres.Page
	}
	if vres.Activities != nil {
		res.Activities = newActivityEntryCollection(vres.Activities)
	}
	return res
}

// newStationActivityPageView projects result type StationActivityPage to
// projected type StationActivityPageView using the "default" view.
func newStationActivityPageView(res *StationActivityPage) *activityviews.StationActivityPageView {
	vres := &activityviews.StationActivityPageView{
		Total: &res.Total,
		Page:  &res.Page,
	}
	if res.Activities != nil {
		vres.Activities = newActivityEntryCollectionView(res.Activities)
	}
	return vres
}

// newActivityEntryCollection converts projected type ActivityEntryCollection
// to service type ActivityEntryCollection.
func newActivityEntryCollection(vres activityviews.ActivityEntryCollectionView) ActivityEntryCollection {
	res := make(ActivityEntryCollection, len(vres))
	for i, n := range vres {
		res[i] = newActivityEntry(n)
	}
	return res
}

// newActivityEntryCollectionView projects result type ActivityEntryCollection
// to projected type ActivityEntryCollectionView using the "default" view.
func newActivityEntryCollectionView(res ActivityEntryCollection) activityviews.ActivityEntryCollectionView {
	vres := make(activityviews.ActivityEntryCollectionView, len(res))
	for i, n := range res {
		vres[i] = newActivityEntryView(n)
	}
	return vres
}

// newActivityEntry converts projected type ActivityEntry to service type
// ActivityEntry.
func newActivityEntry(vres *activityviews.ActivityEntryView) *ActivityEntry {
	res := &ActivityEntry{
		Meta: vres.Meta,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Key != nil {
		res.Key = *vres.Key
	}
	if vres.CreatedAt != nil {
		res.CreatedAt = *vres.CreatedAt
	}
	if vres.Type != nil {
		res.Type = *vres.Type
	}
	if vres.Project != nil {
		res.Project = transformActivityviewsProjectSummaryViewToProjectSummary(vres.Project)
	}
	if vres.Station != nil {
		res.Station = transformActivityviewsStationSummaryViewToStationSummary(vres.Station)
	}
	return res
}

// newActivityEntryView projects result type ActivityEntry to projected type
// ActivityEntryView using the "default" view.
func newActivityEntryView(res *ActivityEntry) *activityviews.ActivityEntryView {
	vres := &activityviews.ActivityEntryView{
		ID:        &res.ID,
		Key:       &res.Key,
		CreatedAt: &res.CreatedAt,
		Type:      &res.Type,
		Meta:      res.Meta,
	}
	if res.Project != nil {
		vres.Project = transformProjectSummaryToActivityviewsProjectSummaryView(res.Project)
	}
	if res.Station != nil {
		vres.Station = transformStationSummaryToActivityviewsStationSummaryView(res.Station)
	}
	return vres
}

// newProjectActivityPage converts projected type ProjectActivityPage to
// service type ProjectActivityPage.
func newProjectActivityPage(vres *activityviews.ProjectActivityPageView) *ProjectActivityPage {
	res := &ProjectActivityPage{}
	if vres.Total != nil {
		res.Total = *vres.Total
	}
	if vres.Page != nil {
		res.Page = *vres.Page
	}
	if vres.Activities != nil {
		res.Activities = newActivityEntryCollection(vres.Activities)
	}
	return res
}

// newProjectActivityPageView projects result type ProjectActivityPage to
// projected type ProjectActivityPageView using the "default" view.
func newProjectActivityPageView(res *ProjectActivityPage) *activityviews.ProjectActivityPageView {
	vres := &activityviews.ProjectActivityPageView{
		Total: &res.Total,
		Page:  &res.Page,
	}
	if res.Activities != nil {
		vres.Activities = newActivityEntryCollectionView(res.Activities)
	}
	return vres
}

// transformActivityviewsProjectSummaryViewToProjectSummary builds a value of
// type *ProjectSummary from a value of type *activityviews.ProjectSummaryView.
func transformActivityviewsProjectSummaryViewToProjectSummary(v *activityviews.ProjectSummaryView) *ProjectSummary {
	if v == nil {
		return nil
	}
	res := &ProjectSummary{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// transformActivityviewsStationSummaryViewToStationSummary builds a value of
// type *StationSummary from a value of type *activityviews.StationSummaryView.
func transformActivityviewsStationSummaryViewToStationSummary(v *activityviews.StationSummaryView) *StationSummary {
	if v == nil {
		return nil
	}
	res := &StationSummary{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// transformProjectSummaryToActivityviewsProjectSummaryView builds a value of
// type *activityviews.ProjectSummaryView from a value of type *ProjectSummary.
func transformProjectSummaryToActivityviewsProjectSummaryView(v *ProjectSummary) *activityviews.ProjectSummaryView {
	res := &activityviews.ProjectSummaryView{
		ID:   &v.ID,
		Name: &v.Name,
	}

	return res
}

// transformStationSummaryToActivityviewsStationSummaryView builds a value of
// type *activityviews.StationSummaryView from a value of type *StationSummary.
func transformStationSummaryToActivityviewsStationSummaryView(v *StationSummary) *activityviews.StationSummaryView {
	res := &activityviews.StationSummaryView{
		ID:   &v.ID,
		Name: &v.Name,
	}

	return res
}
