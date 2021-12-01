// Code generated by goa v3.2.4, DO NOT EDIT.
//
// user service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package user

import (
	"context"
	"io"

	userviews "github.com/fieldkit/cloud/server/api/gen/user/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the user service interface.
type Service interface {
	// Roles implements roles.
	Roles(context.Context, *RolesPayload) (res *AvailableRoles, err error)
	// UploadPhoto implements upload photo.
	UploadPhoto(context.Context, *UploadPhotoPayload, io.ReadCloser) (err error)
	// DownloadPhoto implements download photo.
	DownloadPhoto(context.Context, *DownloadPhotoPayload) (res *DownloadedPhoto, err error)
	// Login implements login.
	Login(context.Context, *LoginPayload) (res *LoginResult, err error)
	// RecoveryLookup implements recovery lookup.
	RecoveryLookup(context.Context, *RecoveryLookupPayload) (err error)
	// Recovery implements recovery.
	Recovery(context.Context, *RecoveryPayload) (err error)
	// Resume implements resume.
	Resume(context.Context, *ResumePayload) (res *ResumeResult, err error)
	// Logout implements logout.
	Logout(context.Context, *LogoutPayload) (err error)
	// Refresh implements refresh.
	Refresh(context.Context, *RefreshPayload) (res *RefreshResult, err error)
	// SendValidation implements send validation.
	SendValidation(context.Context, *SendValidationPayload) (err error)
	// Validate implements validate.
	Validate(context.Context, *ValidatePayload) (res *ValidateResult, err error)
	// Add implements add.
	Add(context.Context, *AddPayload) (res *User, err error)
	// Update implements update.
	Update(context.Context, *UpdatePayload) (res *User, err error)
	// ChangePassword implements change password.
	ChangePassword(context.Context, *ChangePasswordPayload) (res *User, err error)
	// AcceptTnc implements accept tnc.
	AcceptTnc(context.Context, *AcceptTncPayload) (res *User, err error)
	// GetCurrent implements get current.
	GetCurrent(context.Context, *GetCurrentPayload) (res *User, err error)
	// ListByProject implements list by project.
	ListByProject(context.Context, *ListByProjectPayload) (res *ProjectUsers, err error)
	// IssueTransmissionToken implements issue transmission token.
	IssueTransmissionToken(context.Context, *IssueTransmissionTokenPayload) (res *TransmissionToken, err error)
	// ProjectRoles implements project roles.
	ProjectRoles(context.Context) (res ProjectRoleCollection, err error)
	// AdminTermsAndConditions implements admin terms and conditions.
	AdminTermsAndConditions(context.Context, *AdminTermsAndConditionsPayload) (err error)
	// AdminDelete implements admin delete.
	AdminDelete(context.Context, *AdminDeletePayload) (err error)
	// AdminSearch implements admin search.
	AdminSearch(context.Context, *AdminSearchPayload) (res *AdminSearchResult, err error)
	// Mentionables implements mentionables.
	Mentionables(context.Context, *MentionablesPayload) (res *MentionableOptions, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "user"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [23]string{"roles", "upload photo", "download photo", "login", "recovery lookup", "recovery", "resume", "logout", "refresh", "send validation", "validate", "add", "update", "change password", "accept tnc", "get current", "list by project", "issue transmission token", "project roles", "admin terms and conditions", "admin delete", "admin search", "mentionables"}

// RolesPayload is the payload type of the user service roles method.
type RolesPayload struct {
	Auth string
}

// AvailableRoles is the result type of the user service roles method.
type AvailableRoles struct {
	Roles []*AvailableRole
}

// UploadPhotoPayload is the payload type of the user service upload photo
// method.
type UploadPhotoPayload struct {
	Auth          string
	ContentLength int64
	ContentType   string
}

// DownloadPhotoPayload is the payload type of the user service download photo
// method.
type DownloadPhotoPayload struct {
	UserID      int32
	Size        *int32
	IfNoneMatch *string
}

// DownloadedPhoto is the result type of the user service download photo method.
type DownloadedPhoto struct {
	Length      int64
	ContentType string
	Etag        string
	Body        []byte
}

// LoginPayload is the payload type of the user service login method.
type LoginPayload struct {
	Login *LoginFields
}

// LoginResult is the result type of the user service login method.
type LoginResult struct {
	Authorization string
}

// RecoveryLookupPayload is the payload type of the user service recovery
// lookup method.
type RecoveryLookupPayload struct {
	Recovery *RecoveryLookupFields
}

// RecoveryPayload is the payload type of the user service recovery method.
type RecoveryPayload struct {
	Recovery *RecoveryFields
}

// ResumePayload is the payload type of the user service resume method.
type ResumePayload struct {
	Token string
}

// ResumeResult is the result type of the user service resume method.
type ResumeResult struct {
	Authorization string
}

// LogoutPayload is the payload type of the user service logout method.
type LogoutPayload struct {
	Auth string
}

// RefreshPayload is the payload type of the user service refresh method.
type RefreshPayload struct {
	RefreshToken string
}

// RefreshResult is the result type of the user service refresh method.
type RefreshResult struct {
	Authorization string
}

// SendValidationPayload is the payload type of the user service send
// validation method.
type SendValidationPayload struct {
	UserID int32
}

// ValidatePayload is the payload type of the user service validate method.
type ValidatePayload struct {
	Token string
}

// ValidateResult is the result type of the user service validate method.
type ValidateResult struct {
	Location string
}

// AddPayload is the payload type of the user service add method.
type AddPayload struct {
	User *AddUserFields
}

// User is the result type of the user service add method.
type User struct {
	ID        int32
	Name      string
	Email     string
	Bio       string
	Photo     *UserPhoto
	Admin     bool
	UpdatedAt int64
	TncDate   int64
}

// UpdatePayload is the payload type of the user service update method.
type UpdatePayload struct {
	Auth   string
	UserID int32
	Update *UpdateUserFields
}

// ChangePasswordPayload is the payload type of the user service change
// password method.
type ChangePasswordPayload struct {
	Auth   string
	UserID int32
	Change *UpdateUserPasswordFields
}

// AcceptTncPayload is the payload type of the user service accept tnc method.
type AcceptTncPayload struct {
	Auth   string
	UserID int32
	Accept *AcceptTncFields
}

// GetCurrentPayload is the payload type of the user service get current method.
type GetCurrentPayload struct {
	Auth string
}

// ListByProjectPayload is the payload type of the user service list by project
// method.
type ListByProjectPayload struct {
	Auth      string
	ProjectID int32
}

// ProjectUsers is the result type of the user service list by project method.
type ProjectUsers struct {
	Users ProjectUserCollection
}

// IssueTransmissionTokenPayload is the payload type of the user service issue
// transmission token method.
type IssueTransmissionTokenPayload struct {
	Auth string
}

// TransmissionToken is the result type of the user service issue transmission
// token method.
type TransmissionToken struct {
	Token string
	URL   string
}

// ProjectRoleCollection is the result type of the user service project roles
// method.
type ProjectRoleCollection []*ProjectRole

// AdminTermsAndConditionsPayload is the payload type of the user service admin
// terms and conditions method.
type AdminTermsAndConditionsPayload struct {
	Auth   string
	Update *AdminTermsAndConditionsFields
}

// AdminDeletePayload is the payload type of the user service admin delete
// method.
type AdminDeletePayload struct {
	Auth   string
	Delete *AdminDeleteFields
}

// AdminSearchPayload is the payload type of the user service admin search
// method.
type AdminSearchPayload struct {
	Auth  string
	Query string
}

// AdminSearchResult is the result type of the user service admin search method.
type AdminSearchResult struct {
	Users UserCollection
}

// MentionablesPayload is the payload type of the user service mentionables
// method.
type MentionablesPayload struct {
	Auth      string
	ProjectID *int32
	Bookmark  *string
	Query     string
}

// MentionableOptions is the result type of the user service mentionables
// method.
type MentionableOptions struct {
	Users MentionableUserCollection
}

type AvailableRole struct {
	ID   int32
	Name string
}

type LoginFields struct {
	Email    string
	Password string
}

type RecoveryLookupFields struct {
	Email string
}

type RecoveryFields struct {
	Token    string
	Password string
}

type AddUserFields struct {
	Name        string
	Email       string
	Password    string
	InviteToken *string
	TncAccept   *bool
}

type UserPhoto struct {
	URL *string
}

type UpdateUserFields struct {
	Name  string
	Email string
	Bio   string
}

type UpdateUserPasswordFields struct {
	OldPassword string
	NewPassword string
}

type AcceptTncFields struct {
	Accept bool
}

type ProjectUserCollection []*ProjectUser

type ProjectUser struct {
	User       *User
	Role       string
	Membership string
	Invited    bool
	Accepted   bool
	Rejected   bool
}

type ProjectRole struct {
	ID   int32
	Name string
}

type AdminTermsAndConditionsFields struct {
	Email string
}

type AdminDeleteFields struct {
	Email    string
	Password string
}

type UserCollection []*User

type MentionableUserCollection []*MentionableUser

type MentionableUser struct {
	ID      int32
	Name    string
	Mention string
	Photo   *UserPhoto
}

// MakeUserUnverified builds a goa.ServiceError from an error.
func MakeUserUnverified(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "user-unverified",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeUserEmailRegistered builds a goa.ServiceError from an error.
func MakeUserEmailRegistered(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "user-email-registered",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
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

// NewAvailableRoles initializes result type AvailableRoles from viewed result
// type AvailableRoles.
func NewAvailableRoles(vres *userviews.AvailableRoles) *AvailableRoles {
	return newAvailableRoles(vres.Projected)
}

// NewViewedAvailableRoles initializes viewed result type AvailableRoles from
// result type AvailableRoles using the given view.
func NewViewedAvailableRoles(res *AvailableRoles, view string) *userviews.AvailableRoles {
	p := newAvailableRolesView(res)
	return &userviews.AvailableRoles{Projected: p, View: "default"}
}

// NewDownloadedPhoto initializes result type DownloadedPhoto from viewed
// result type DownloadedPhoto.
func NewDownloadedPhoto(vres *userviews.DownloadedPhoto) *DownloadedPhoto {
	return newDownloadedPhoto(vres.Projected)
}

// NewViewedDownloadedPhoto initializes viewed result type DownloadedPhoto from
// result type DownloadedPhoto using the given view.
func NewViewedDownloadedPhoto(res *DownloadedPhoto, view string) *userviews.DownloadedPhoto {
	p := newDownloadedPhotoView(res)
	return &userviews.DownloadedPhoto{Projected: p, View: "default"}
}

// NewUser initializes result type User from viewed result type User.
func NewUser(vres *userviews.User) *User {
	return newUser(vres.Projected)
}

// NewViewedUser initializes viewed result type User from result type User
// using the given view.
func NewViewedUser(res *User, view string) *userviews.User {
	p := newUserView(res)
	return &userviews.User{Projected: p, View: "default"}
}

// NewProjectUsers initializes result type ProjectUsers from viewed result type
// ProjectUsers.
func NewProjectUsers(vres *userviews.ProjectUsers) *ProjectUsers {
	return newProjectUsers(vres.Projected)
}

// NewViewedProjectUsers initializes viewed result type ProjectUsers from
// result type ProjectUsers using the given view.
func NewViewedProjectUsers(res *ProjectUsers, view string) *userviews.ProjectUsers {
	p := newProjectUsersView(res)
	return &userviews.ProjectUsers{Projected: p, View: "default"}
}

// NewTransmissionToken initializes result type TransmissionToken from viewed
// result type TransmissionToken.
func NewTransmissionToken(vres *userviews.TransmissionToken) *TransmissionToken {
	return newTransmissionToken(vres.Projected)
}

// NewViewedTransmissionToken initializes viewed result type TransmissionToken
// from result type TransmissionToken using the given view.
func NewViewedTransmissionToken(res *TransmissionToken, view string) *userviews.TransmissionToken {
	p := newTransmissionTokenView(res)
	return &userviews.TransmissionToken{Projected: p, View: "default"}
}

// NewProjectRoleCollection initializes result type ProjectRoleCollection from
// viewed result type ProjectRoleCollection.
func NewProjectRoleCollection(vres userviews.ProjectRoleCollection) ProjectRoleCollection {
	return newProjectRoleCollection(vres.Projected)
}

// NewViewedProjectRoleCollection initializes viewed result type
// ProjectRoleCollection from result type ProjectRoleCollection using the given
// view.
func NewViewedProjectRoleCollection(res ProjectRoleCollection, view string) userviews.ProjectRoleCollection {
	p := newProjectRoleCollectionView(res)
	return userviews.ProjectRoleCollection{Projected: p, View: "default"}
}

// NewMentionableOptions initializes result type MentionableOptions from viewed
// result type MentionableOptions.
func NewMentionableOptions(vres *userviews.MentionableOptions) *MentionableOptions {
	return newMentionableOptions(vres.Projected)
}

// NewViewedMentionableOptions initializes viewed result type
// MentionableOptions from result type MentionableOptions using the given view.
func NewViewedMentionableOptions(res *MentionableOptions, view string) *userviews.MentionableOptions {
	p := newMentionableOptionsView(res)
	return &userviews.MentionableOptions{Projected: p, View: "default"}
}

// newAvailableRoles converts projected type AvailableRoles to service type
// AvailableRoles.
func newAvailableRoles(vres *userviews.AvailableRolesView) *AvailableRoles {
	res := &AvailableRoles{}
	if vres.Roles != nil {
		res.Roles = make([]*AvailableRole, len(vres.Roles))
		for i, val := range vres.Roles {
			res.Roles[i] = transformUserviewsAvailableRoleViewToAvailableRole(val)
		}
	}
	return res
}

// newAvailableRolesView projects result type AvailableRoles to projected type
// AvailableRolesView using the "default" view.
func newAvailableRolesView(res *AvailableRoles) *userviews.AvailableRolesView {
	vres := &userviews.AvailableRolesView{}
	if res.Roles != nil {
		vres.Roles = make([]*userviews.AvailableRoleView, len(res.Roles))
		for i, val := range res.Roles {
			vres.Roles[i] = transformAvailableRoleToUserviewsAvailableRoleView(val)
		}
	}
	return vres
}

// newDownloadedPhoto converts projected type DownloadedPhoto to service type
// DownloadedPhoto.
func newDownloadedPhoto(vres *userviews.DownloadedPhotoView) *DownloadedPhoto {
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
func newDownloadedPhotoView(res *DownloadedPhoto) *userviews.DownloadedPhotoView {
	vres := &userviews.DownloadedPhotoView{
		Length:      &res.Length,
		ContentType: &res.ContentType,
		Etag:        &res.Etag,
		Body:        res.Body,
	}
	return vres
}

// newUser converts projected type User to service type User.
func newUser(vres *userviews.UserView) *User {
	res := &User{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Email != nil {
		res.Email = *vres.Email
	}
	if vres.Bio != nil {
		res.Bio = *vres.Bio
	}
	if vres.Admin != nil {
		res.Admin = *vres.Admin
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.TncDate != nil {
		res.TncDate = *vres.TncDate
	}
	if vres.Photo != nil {
		res.Photo = transformUserviewsUserPhotoViewToUserPhoto(vres.Photo)
	}
	return res
}

// newUserView projects result type User to projected type UserView using the
// "default" view.
func newUserView(res *User) *userviews.UserView {
	vres := &userviews.UserView{
		ID:        &res.ID,
		Name:      &res.Name,
		Email:     &res.Email,
		Bio:       &res.Bio,
		Admin:     &res.Admin,
		UpdatedAt: &res.UpdatedAt,
		TncDate:   &res.TncDate,
	}
	if res.Photo != nil {
		vres.Photo = transformUserPhotoToUserviewsUserPhotoView(res.Photo)
	}
	return vres
}

// newProjectUsers converts projected type ProjectUsers to service type
// ProjectUsers.
func newProjectUsers(vres *userviews.ProjectUsersView) *ProjectUsers {
	res := &ProjectUsers{}
	if vres.Users != nil {
		res.Users = newProjectUserCollection(vres.Users)
	}
	return res
}

// newProjectUsersView projects result type ProjectUsers to projected type
// ProjectUsersView using the "default" view.
func newProjectUsersView(res *ProjectUsers) *userviews.ProjectUsersView {
	vres := &userviews.ProjectUsersView{}
	if res.Users != nil {
		vres.Users = newProjectUserCollectionView(res.Users)
	}
	return vres
}

// newProjectUserCollection converts projected type ProjectUserCollection to
// service type ProjectUserCollection.
func newProjectUserCollection(vres userviews.ProjectUserCollectionView) ProjectUserCollection {
	res := make(ProjectUserCollection, len(vres))
	for i, n := range vres {
		res[i] = newProjectUser(n)
	}
	return res
}

// newProjectUserCollectionView projects result type ProjectUserCollection to
// projected type ProjectUserCollectionView using the "default" view.
func newProjectUserCollectionView(res ProjectUserCollection) userviews.ProjectUserCollectionView {
	vres := make(userviews.ProjectUserCollectionView, len(res))
	for i, n := range res {
		vres[i] = newProjectUserView(n)
	}
	return vres
}

// newProjectUser converts projected type ProjectUser to service type
// ProjectUser.
func newProjectUser(vres *userviews.ProjectUserView) *ProjectUser {
	res := &ProjectUser{}
	if vres.Role != nil {
		res.Role = *vres.Role
	}
	if vres.Membership != nil {
		res.Membership = *vres.Membership
	}
	if vres.Invited != nil {
		res.Invited = *vres.Invited
	}
	if vres.Accepted != nil {
		res.Accepted = *vres.Accepted
	}
	if vres.Rejected != nil {
		res.Rejected = *vres.Rejected
	}
	if vres.User != nil {
		res.User = newUser(vres.User)
	}
	return res
}

// newProjectUserView projects result type ProjectUser to projected type
// ProjectUserView using the "default" view.
func newProjectUserView(res *ProjectUser) *userviews.ProjectUserView {
	vres := &userviews.ProjectUserView{
		Role:       &res.Role,
		Membership: &res.Membership,
		Invited:    &res.Invited,
		Accepted:   &res.Accepted,
		Rejected:   &res.Rejected,
	}
	if res.User != nil {
		vres.User = newUserView(res.User)
	}
	return vres
}

// newTransmissionToken converts projected type TransmissionToken to service
// type TransmissionToken.
func newTransmissionToken(vres *userviews.TransmissionTokenView) *TransmissionToken {
	res := &TransmissionToken{}
	if vres.Token != nil {
		res.Token = *vres.Token
	}
	if vres.URL != nil {
		res.URL = *vres.URL
	}
	return res
}

// newTransmissionTokenView projects result type TransmissionToken to projected
// type TransmissionTokenView using the "default" view.
func newTransmissionTokenView(res *TransmissionToken) *userviews.TransmissionTokenView {
	vres := &userviews.TransmissionTokenView{
		Token: &res.Token,
		URL:   &res.URL,
	}
	return vres
}

// newProjectRoleCollection converts projected type ProjectRoleCollection to
// service type ProjectRoleCollection.
func newProjectRoleCollection(vres userviews.ProjectRoleCollectionView) ProjectRoleCollection {
	res := make(ProjectRoleCollection, len(vres))
	for i, n := range vres {
		res[i] = newProjectRole(n)
	}
	return res
}

// newProjectRoleCollectionView projects result type ProjectRoleCollection to
// projected type ProjectRoleCollectionView using the "default" view.
func newProjectRoleCollectionView(res ProjectRoleCollection) userviews.ProjectRoleCollectionView {
	vres := make(userviews.ProjectRoleCollectionView, len(res))
	for i, n := range res {
		vres[i] = newProjectRoleView(n)
	}
	return vres
}

// newProjectRole converts projected type ProjectRole to service type
// ProjectRole.
func newProjectRole(vres *userviews.ProjectRoleView) *ProjectRole {
	res := &ProjectRole{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	return res
}

// newProjectRoleView projects result type ProjectRole to projected type
// ProjectRoleView using the "default" view.
func newProjectRoleView(res *ProjectRole) *userviews.ProjectRoleView {
	vres := &userviews.ProjectRoleView{
		ID:   &res.ID,
		Name: &res.Name,
	}
	return vres
}

// newMentionableOptions converts projected type MentionableOptions to service
// type MentionableOptions.
func newMentionableOptions(vres *userviews.MentionableOptionsView) *MentionableOptions {
	res := &MentionableOptions{}
	if vres.Users != nil {
		res.Users = newMentionableUserCollection(vres.Users)
	}
	return res
}

// newMentionableOptionsView projects result type MentionableOptions to
// projected type MentionableOptionsView using the "default" view.
func newMentionableOptionsView(res *MentionableOptions) *userviews.MentionableOptionsView {
	vres := &userviews.MentionableOptionsView{}
	if res.Users != nil {
		vres.Users = newMentionableUserCollectionView(res.Users)
	}
	return vres
}

// newMentionableUserCollection converts projected type
// MentionableUserCollection to service type MentionableUserCollection.
func newMentionableUserCollection(vres userviews.MentionableUserCollectionView) MentionableUserCollection {
	res := make(MentionableUserCollection, len(vres))
	for i, n := range vres {
		res[i] = newMentionableUser(n)
	}
	return res
}

// newMentionableUserCollectionView projects result type
// MentionableUserCollection to projected type MentionableUserCollectionView
// using the "default" view.
func newMentionableUserCollectionView(res MentionableUserCollection) userviews.MentionableUserCollectionView {
	vres := make(userviews.MentionableUserCollectionView, len(res))
	for i, n := range res {
		vres[i] = newMentionableUserView(n)
	}
	return vres
}

// newMentionableUser converts projected type MentionableUser to service type
// MentionableUser.
func newMentionableUser(vres *userviews.MentionableUserView) *MentionableUser {
	res := &MentionableUser{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Mention != nil {
		res.Mention = *vres.Mention
	}
	if vres.Photo != nil {
		res.Photo = transformUserviewsUserPhotoViewToUserPhoto(vres.Photo)
	}
	return res
}

// newMentionableUserView projects result type MentionableUser to projected
// type MentionableUserView using the "default" view.
func newMentionableUserView(res *MentionableUser) *userviews.MentionableUserView {
	vres := &userviews.MentionableUserView{
		ID:      &res.ID,
		Name:    &res.Name,
		Mention: &res.Mention,
	}
	if res.Photo != nil {
		vres.Photo = transformUserPhotoToUserviewsUserPhotoView(res.Photo)
	}
	return vres
}

// transformUserviewsAvailableRoleViewToAvailableRole builds a value of type
// *AvailableRole from a value of type *userviews.AvailableRoleView.
func transformUserviewsAvailableRoleViewToAvailableRole(v *userviews.AvailableRoleView) *AvailableRole {
	if v == nil {
		return nil
	}
	res := &AvailableRole{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// transformAvailableRoleToUserviewsAvailableRoleView builds a value of type
// *userviews.AvailableRoleView from a value of type *AvailableRole.
func transformAvailableRoleToUserviewsAvailableRoleView(v *AvailableRole) *userviews.AvailableRoleView {
	res := &userviews.AvailableRoleView{
		ID:   &v.ID,
		Name: &v.Name,
	}

	return res
}

// transformUserviewsUserPhotoViewToUserPhoto builds a value of type *UserPhoto
// from a value of type *userviews.UserPhotoView.
func transformUserviewsUserPhotoViewToUserPhoto(v *userviews.UserPhotoView) *UserPhoto {
	if v == nil {
		return nil
	}
	res := &UserPhoto{
		URL: v.URL,
	}

	return res
}

// transformUserPhotoToUserviewsUserPhotoView builds a value of type
// *userviews.UserPhotoView from a value of type *UserPhoto.
func transformUserPhotoToUserviewsUserPhotoView(v *UserPhoto) *userviews.UserPhotoView {
	if v == nil {
		return nil
	}
	res := &userviews.UserPhotoView{
		URL: v.URL,
	}

	return res
}
