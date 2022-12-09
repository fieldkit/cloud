// Code generated by goa v3.2.4, DO NOT EDIT.
//
// project service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package project

import (
	"context"
	"io"

	projectviews "github.com/fieldkit/cloud/server/api/gen/project/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the project service interface.
type Service interface {
	// AddUpdate implements add update.
	AddUpdate(context.Context, *AddUpdatePayload) (res *ProjectUpdate, err error)
	// DeleteUpdate implements delete update.
	DeleteUpdate(context.Context, *DeleteUpdatePayload) (err error)
	// ModifyUpdate implements modify update.
	ModifyUpdate(context.Context, *ModifyUpdatePayload) (res *ProjectUpdate, err error)
	// Invites implements invites.
	Invites(context.Context, *InvitesPayload) (res *PendingInvites, err error)
	// LookupInvite implements lookup invite.
	LookupInvite(context.Context, *LookupInvitePayload) (res *PendingInvites, err error)
	// AcceptProjectInvite implements accept project invite.
	AcceptProjectInvite(context.Context, *AcceptProjectInvitePayload) (err error)
	// RejectProjectInvite implements reject project invite.
	RejectProjectInvite(context.Context, *RejectProjectInvitePayload) (err error)
	// AcceptInvite implements accept invite.
	AcceptInvite(context.Context, *AcceptInvitePayload) (err error)
	// RejectInvite implements reject invite.
	RejectInvite(context.Context, *RejectInvitePayload) (err error)
	// Add implements add.
	Add(context.Context, *AddPayload) (res *Project, err error)
	// Update implements update.
	Update(context.Context, *UpdatePayload) (res *Project, err error)
	// Get implements get.
	Get(context.Context, *GetPayload) (res *Project, err error)
	// ListCommunity implements list community.
	ListCommunity(context.Context, *ListCommunityPayload) (res *Projects, err error)
	// ListMine implements list mine.
	ListMine(context.Context, *ListMinePayload) (res *Projects, err error)
	// Invite implements invite.
	Invite(context.Context, *InvitePayload) (err error)
	// EditUser implements edit user.
	EditUser(context.Context, *EditUserPayload) (err error)
	// RemoveUser implements remove user.
	RemoveUser(context.Context, *RemoveUserPayload) (err error)
	// AddStation implements add station.
	AddStation(context.Context, *AddStationPayload) (err error)
	// RemoveStation implements remove station.
	RemoveStation(context.Context, *RemoveStationPayload) (err error)
	// Delete implements delete.
	Delete(context.Context, *DeletePayload) (err error)
	// UploadPhoto implements upload photo.
	UploadPhoto(context.Context, *UploadPhotoPayload, io.ReadCloser) (err error)
	// DownloadPhoto implements download photo.
	DownloadPhoto(context.Context, *DownloadPhotoPayload) (res *DownloadedPhoto, err error)
	// GetProjectsForStation implements get projects for station.
	GetProjectsForStation(context.Context, *GetProjectsForStationPayload) (res *Projects, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "project"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [23]string{"add update", "delete update", "modify update", "invites", "lookup invite", "accept project invite", "reject project invite", "accept invite", "reject invite", "add", "update", "get", "list community", "list mine", "invite", "edit user", "remove user", "add station", "remove station", "delete", "upload photo", "download photo", "get projects for station"}

// AddUpdatePayload is the payload type of the project service add update
// method.
type AddUpdatePayload struct {
	Auth      string
	ProjectID int32
	Body      string
}

// ProjectUpdate is the result type of the project service add update method.
type ProjectUpdate struct {
	ID        int64
	Body      string
	CreatedAt int64
}

// DeleteUpdatePayload is the payload type of the project service delete update
// method.
type DeleteUpdatePayload struct {
	Auth      string
	ProjectID int32
	UpdateID  int64
}

// ModifyUpdatePayload is the payload type of the project service modify update
// method.
type ModifyUpdatePayload struct {
	Auth      string
	ProjectID int32
	UpdateID  int64
	Body      string
}

// InvitesPayload is the payload type of the project service invites method.
type InvitesPayload struct {
	Auth string
}

// PendingInvites is the result type of the project service invites method.
type PendingInvites struct {
	Pending  []*PendingInvite
	Projects ProjectCollection
}

// LookupInvitePayload is the payload type of the project service lookup invite
// method.
type LookupInvitePayload struct {
	Auth  string
	Token string
}

// AcceptProjectInvitePayload is the payload type of the project service accept
// project invite method.
type AcceptProjectInvitePayload struct {
	Auth      string
	ProjectID int32
}

// RejectProjectInvitePayload is the payload type of the project service reject
// project invite method.
type RejectProjectInvitePayload struct {
	Auth      string
	ProjectID int32
}

// AcceptInvitePayload is the payload type of the project service accept invite
// method.
type AcceptInvitePayload struct {
	Auth  string
	ID    int64
	Token *string
}

// RejectInvitePayload is the payload type of the project service reject invite
// method.
type RejectInvitePayload struct {
	Auth  string
	ID    int64
	Token *string
}

// AddPayload is the payload type of the project service add method.
type AddPayload struct {
	Auth    string
	Project *AddProjectFields
}

// Project is the result type of the project service add method.
type Project struct {
	ID           int32
	Name         string
	Description  string
	Goal         string
	Location     string
	Tags         string
	Privacy      int32
	StartTime    *string
	EndTime      *string
	Photo        *string
	ReadOnly     bool
	ShowStations bool
	Bounds       *ProjectBounds
	Following    *ProjectFollowing
}

// UpdatePayload is the payload type of the project service update method.
type UpdatePayload struct {
	Auth      string
	ProjectID int32
	Project   *AddProjectFields
}

// GetPayload is the payload type of the project service get method.
type GetPayload struct {
	Auth      *string
	ProjectID int32
}

// ListCommunityPayload is the payload type of the project service list
// community method.
type ListCommunityPayload struct {
	Auth *string
}

// Projects is the result type of the project service list community method.
type Projects struct {
	Projects ProjectCollection
}

// ListMinePayload is the payload type of the project service list mine method.
type ListMinePayload struct {
	Auth string
}

// InvitePayload is the payload type of the project service invite method.
type InvitePayload struct {
	Auth      string
	ProjectID int32
	Invite    *InviteUserFields
}

// EditUserPayload is the payload type of the project service edit user method.
type EditUserPayload struct {
	Auth      string
	ProjectID int32
	Edit      *EditUserFields
}

// RemoveUserPayload is the payload type of the project service remove user
// method.
type RemoveUserPayload struct {
	Auth      string
	ProjectID int32
	Remove    *RemoveUserFields
}

// AddStationPayload is the payload type of the project service add station
// method.
type AddStationPayload struct {
	Auth      string
	ProjectID int32
	StationID int32
}

// RemoveStationPayload is the payload type of the project service remove
// station method.
type RemoveStationPayload struct {
	Auth      string
	ProjectID int32
	StationID int32
}

// DeletePayload is the payload type of the project service delete method.
type DeletePayload struct {
	Auth      string
	ProjectID int32
}

// UploadPhotoPayload is the payload type of the project service upload photo
// method.
type UploadPhotoPayload struct {
	Auth          string
	ProjectID     int32
	ContentLength int64
	ContentType   string
}

// DownloadPhotoPayload is the payload type of the project service download
// photo method.
type DownloadPhotoPayload struct {
	Auth        *string
	ProjectID   int32
	Size        *int32
	IfNoneMatch *string
}

// DownloadedPhoto is the result type of the project service download photo
// method.
type DownloadedPhoto struct {
	Length      int64
	ContentType string
	Etag        string
	Body        []byte
}

// GetProjectsForStationPayload is the payload type of the project service get
// projects for station method.
type GetProjectsForStationPayload struct {
	Auth *string
	ID   int32
}

type PendingInvite struct {
	ID      int64
	Project *ProjectSummary
	Time    int64
	Role    int32
}

type ProjectSummary struct {
	ID   int64
	Name string
}

type ProjectCollection []*Project

type ProjectBounds struct {
	Min []float64
	Max []float64
}

type ProjectFollowing struct {
	Total     int32
	Following bool
}

type AddProjectFields struct {
	Name         string
	Description  string
	Goal         *string
	Location     *string
	Tags         *string
	Privacy      *int32
	StartTime    *string
	EndTime      *string
	Bounds       *ProjectBounds
	ShowStations *bool
}

type InviteUserFields struct {
	Email string
	Role  int32
}

type EditUserFields struct {
	Email string
	Role  int32
}

type RemoveUserFields struct {
	Email string
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

// NewProjectUpdate initializes result type ProjectUpdate from viewed result
// type ProjectUpdate.
func NewProjectUpdate(vres *projectviews.ProjectUpdate) *ProjectUpdate {
	return newProjectUpdate(vres.Projected)
}

// NewViewedProjectUpdate initializes viewed result type ProjectUpdate from
// result type ProjectUpdate using the given view.
func NewViewedProjectUpdate(res *ProjectUpdate, view string) *projectviews.ProjectUpdate {
	p := newProjectUpdateView(res)
	return &projectviews.ProjectUpdate{Projected: p, View: "default"}
}

// NewPendingInvites initializes result type PendingInvites from viewed result
// type PendingInvites.
func NewPendingInvites(vres *projectviews.PendingInvites) *PendingInvites {
	return newPendingInvites(vres.Projected)
}

// NewViewedPendingInvites initializes viewed result type PendingInvites from
// result type PendingInvites using the given view.
func NewViewedPendingInvites(res *PendingInvites, view string) *projectviews.PendingInvites {
	p := newPendingInvitesView(res)
	return &projectviews.PendingInvites{Projected: p, View: "default"}
}

// NewProject initializes result type Project from viewed result type Project.
func NewProject(vres *projectviews.Project) *Project {
	return newProject(vres.Projected)
}

// NewViewedProject initializes viewed result type Project from result type
// Project using the given view.
func NewViewedProject(res *Project, view string) *projectviews.Project {
	p := newProjectView(res)
	return &projectviews.Project{Projected: p, View: "default"}
}

// NewProjects initializes result type Projects from viewed result type
// Projects.
func NewProjects(vres *projectviews.Projects) *Projects {
	return newProjects(vres.Projected)
}

// NewViewedProjects initializes viewed result type Projects from result type
// Projects using the given view.
func NewViewedProjects(res *Projects, view string) *projectviews.Projects {
	p := newProjectsView(res)
	return &projectviews.Projects{Projected: p, View: "default"}
}

// NewDownloadedPhoto initializes result type DownloadedPhoto from viewed
// result type DownloadedPhoto.
func NewDownloadedPhoto(vres *projectviews.DownloadedPhoto) *DownloadedPhoto {
	return newDownloadedPhoto(vres.Projected)
}

// NewViewedDownloadedPhoto initializes viewed result type DownloadedPhoto from
// result type DownloadedPhoto using the given view.
func NewViewedDownloadedPhoto(res *DownloadedPhoto, view string) *projectviews.DownloadedPhoto {
	p := newDownloadedPhotoView(res)
	return &projectviews.DownloadedPhoto{Projected: p, View: "default"}
}

// newProjectUpdate converts projected type ProjectUpdate to service type
// ProjectUpdate.
func newProjectUpdate(vres *projectviews.ProjectUpdateView) *ProjectUpdate {
	res := &ProjectUpdate{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Body != nil {
		res.Body = *vres.Body
	}
	return res
}

// newProjectUpdateView projects result type ProjectUpdate to projected type
// ProjectUpdateView using the "default" view.
func newProjectUpdateView(res *ProjectUpdate) *projectviews.ProjectUpdateView {
	vres := &projectviews.ProjectUpdateView{
		ID:   &res.ID,
		Body: &res.Body,
	}
	return vres
}

// newPendingInvites converts projected type PendingInvites to service type
// PendingInvites.
func newPendingInvites(vres *projectviews.PendingInvitesView) *PendingInvites {
	res := &PendingInvites{}
	if vres.Pending != nil {
		res.Pending = make([]*PendingInvite, len(vres.Pending))
		for i, val := range vres.Pending {
			res.Pending[i] = transformProjectviewsPendingInviteViewToPendingInvite(val)
		}
	}
	if vres.Projects != nil {
		res.Projects = newProjectCollection(vres.Projects)
	}
	return res
}

// newPendingInvitesView projects result type PendingInvites to projected type
// PendingInvitesView using the "default" view.
func newPendingInvitesView(res *PendingInvites) *projectviews.PendingInvitesView {
	vres := &projectviews.PendingInvitesView{}
	if res.Pending != nil {
		vres.Pending = make([]*projectviews.PendingInviteView, len(res.Pending))
		for i, val := range res.Pending {
			vres.Pending[i] = transformPendingInviteToProjectviewsPendingInviteView(val)
		}
	}
	if res.Projects != nil {
		vres.Projects = newProjectCollectionView(res.Projects)
	}
	return vres
}

// newProjectCollection converts projected type ProjectCollection to service
// type ProjectCollection.
func newProjectCollection(vres projectviews.ProjectCollectionView) ProjectCollection {
	res := make(ProjectCollection, len(vres))
	for i, n := range vres {
		res[i] = newProject(n)
	}
	return res
}

// newProjectCollectionView projects result type ProjectCollection to projected
// type ProjectCollectionView using the "default" view.
func newProjectCollectionView(res ProjectCollection) projectviews.ProjectCollectionView {
	vres := make(projectviews.ProjectCollectionView, len(res))
	for i, n := range res {
		vres[i] = newProjectView(n)
	}
	return vres
}

// newProject converts projected type Project to service type Project.
func newProject(vres *projectviews.ProjectView) *Project {
	res := &Project{
		StartTime: vres.StartTime,
		EndTime:   vres.EndTime,
		Photo:     vres.Photo,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Description != nil {
		res.Description = *vres.Description
	}
	if vres.Goal != nil {
		res.Goal = *vres.Goal
	}
	if vres.Location != nil {
		res.Location = *vres.Location
	}
	if vres.Tags != nil {
		res.Tags = *vres.Tags
	}
	if vres.Privacy != nil {
		res.Privacy = *vres.Privacy
	}
	if vres.ReadOnly != nil {
		res.ReadOnly = *vres.ReadOnly
	}
	if vres.ShowStations != nil {
		res.ShowStations = *vres.ShowStations
	}
	if vres.Bounds != nil {
		res.Bounds = transformProjectviewsProjectBoundsViewToProjectBounds(vres.Bounds)
	}
	if vres.Following != nil {
		res.Following = transformProjectviewsProjectFollowingViewToProjectFollowing(vres.Following)
	}
	return res
}

// newProjectView projects result type Project to projected type ProjectView
// using the "default" view.
func newProjectView(res *Project) *projectviews.ProjectView {
	vres := &projectviews.ProjectView{
		ID:           &res.ID,
		Name:         &res.Name,
		Description:  &res.Description,
		Goal:         &res.Goal,
		Location:     &res.Location,
		Tags:         &res.Tags,
		Privacy:      &res.Privacy,
		StartTime:    res.StartTime,
		EndTime:      res.EndTime,
		Photo:        res.Photo,
		ReadOnly:     &res.ReadOnly,
		ShowStations: &res.ShowStations,
	}
	if res.Bounds != nil {
		vres.Bounds = transformProjectBoundsToProjectviewsProjectBoundsView(res.Bounds)
	}
	if res.Following != nil {
		vres.Following = transformProjectFollowingToProjectviewsProjectFollowingView(res.Following)
	}
	return vres
}

// newProjects converts projected type Projects to service type Projects.
func newProjects(vres *projectviews.ProjectsView) *Projects {
	res := &Projects{}
	if vres.Projects != nil {
		res.Projects = newProjectCollection(vres.Projects)
	}
	return res
}

// newProjectsView projects result type Projects to projected type ProjectsView
// using the "default" view.
func newProjectsView(res *Projects) *projectviews.ProjectsView {
	vres := &projectviews.ProjectsView{}
	if res.Projects != nil {
		vres.Projects = newProjectCollectionView(res.Projects)
	}
	return vres
}

// newDownloadedPhoto converts projected type DownloadedPhoto to service type
// DownloadedPhoto.
func newDownloadedPhoto(vres *projectviews.DownloadedPhotoView) *DownloadedPhoto {
	res := &DownloadedPhoto{
		Body: vres.Body,
	}
	if vres.Length != nil {
		res.Length = *vres.Length
	}
	if vres.ContentType != nil {
		res.ContentType = *vres.ContentType
	}
	if vres.Etag != nil {
		res.Etag = *vres.Etag
	}
	return res
}

// newDownloadedPhotoView projects result type DownloadedPhoto to projected
// type DownloadedPhotoView using the "default" view.
func newDownloadedPhotoView(res *DownloadedPhoto) *projectviews.DownloadedPhotoView {
	vres := &projectviews.DownloadedPhotoView{
		Length:      &res.Length,
		ContentType: &res.ContentType,
		Etag:        &res.Etag,
		Body:        res.Body,
	}
	return vres
}

// transformProjectviewsPendingInviteViewToPendingInvite builds a value of type
// *PendingInvite from a value of type *projectviews.PendingInviteView.
func transformProjectviewsPendingInviteViewToPendingInvite(v *projectviews.PendingInviteView) *PendingInvite {
	if v == nil {
		return nil
	}
	res := &PendingInvite{
		ID:   *v.ID,
		Time: *v.Time,
		Role: *v.Role,
	}
	if v.Project != nil {
		res.Project = transformProjectviewsProjectSummaryViewToProjectSummary(v.Project)
	}

	return res
}

// transformProjectviewsProjectSummaryViewToProjectSummary builds a value of
// type *ProjectSummary from a value of type *projectviews.ProjectSummaryView.
func transformProjectviewsProjectSummaryViewToProjectSummary(v *projectviews.ProjectSummaryView) *ProjectSummary {
	res := &ProjectSummary{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// transformPendingInviteToProjectviewsPendingInviteView builds a value of type
// *projectviews.PendingInviteView from a value of type *PendingInvite.
func transformPendingInviteToProjectviewsPendingInviteView(v *PendingInvite) *projectviews.PendingInviteView {
	res := &projectviews.PendingInviteView{
		ID:   &v.ID,
		Time: &v.Time,
		Role: &v.Role,
	}
	if v.Project != nil {
		res.Project = transformProjectSummaryToProjectviewsProjectSummaryView(v.Project)
	}

	return res
}

// transformProjectSummaryToProjectviewsProjectSummaryView builds a value of
// type *projectviews.ProjectSummaryView from a value of type *ProjectSummary.
func transformProjectSummaryToProjectviewsProjectSummaryView(v *ProjectSummary) *projectviews.ProjectSummaryView {
	res := &projectviews.ProjectSummaryView{
		ID:   &v.ID,
		Name: &v.Name,
	}

	return res
}

// transformProjectviewsProjectBoundsViewToProjectBounds builds a value of type
// *ProjectBounds from a value of type *projectviews.ProjectBoundsView.
func transformProjectviewsProjectBoundsViewToProjectBounds(v *projectviews.ProjectBoundsView) *ProjectBounds {
	if v == nil {
		return nil
	}
	res := &ProjectBounds{}
	if v.Min != nil {
		res.Min = make([]float64, len(v.Min))
		for i, val := range v.Min {
			res.Min[i] = val
		}
	}
	if v.Max != nil {
		res.Max = make([]float64, len(v.Max))
		for i, val := range v.Max {
			res.Max[i] = val
		}
	}

	return res
}

// transformProjectviewsProjectFollowingViewToProjectFollowing builds a value
// of type *ProjectFollowing from a value of type
// *projectviews.ProjectFollowingView.
func transformProjectviewsProjectFollowingViewToProjectFollowing(v *projectviews.ProjectFollowingView) *ProjectFollowing {
	if v == nil {
		return nil
	}
	res := &ProjectFollowing{
		Total:     *v.Total,
		Following: *v.Following,
	}

	return res
}

// transformProjectBoundsToProjectviewsProjectBoundsView builds a value of type
// *projectviews.ProjectBoundsView from a value of type *ProjectBounds.
func transformProjectBoundsToProjectviewsProjectBoundsView(v *ProjectBounds) *projectviews.ProjectBoundsView {
	res := &projectviews.ProjectBoundsView{}
	if v.Min != nil {
		res.Min = make([]float64, len(v.Min))
		for i, val := range v.Min {
			res.Min[i] = val
		}
	}
	if v.Max != nil {
		res.Max = make([]float64, len(v.Max))
		for i, val := range v.Max {
			res.Max[i] = val
		}
	}

	return res
}

// transformProjectFollowingToProjectviewsProjectFollowingView builds a value
// of type *projectviews.ProjectFollowingView from a value of type
// *ProjectFollowing.
func transformProjectFollowingToProjectviewsProjectFollowingView(v *ProjectFollowing) *projectviews.ProjectFollowingView {
	res := &projectviews.ProjectFollowingView{
		Total:     &v.Total,
		Following: &v.Following,
	}

	return res
}
