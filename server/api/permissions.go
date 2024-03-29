package api

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

type ProjectPermissions interface {
	Permissions
	Project() *data.Project
	CanView() error
	CanModify() error
}

type StationPermissions interface {
	Permissions
	Station() *data.Station
	CanView() error
	CanModify() error
	IsReadOnly() (bool, error)
}

type DiscussionPermissions interface {
	Permissions
	Discussion() *data.DiscussionPost
	CanDelete() error
}

type Permissions interface {
	Unwrap() (permissions Permissions, err error)
	UserID() int32
	Anonymous() bool
	MaybeUserID() *int32
	RefreshToken() string
	RequireAdmin() error
	IsAdmin() bool
	ForProjectByID(id int32) (permissions ProjectPermissions, err error)
	ForProject(project *data.Project) (permissions ProjectPermissions, err error)
	ForStationByID(id int) (permissions StationPermissions, err error)
	ForStationByDeviceID(id []byte) (permissions StationPermissions, err error)
	ForStation(station *data.Station) (permissions StationPermissions, err error)
	ForDiscussions(discussion *data.DiscussionPost) (permissions DiscussionPermissions, err error)
}

type unwrappedPermissions struct {
	userID       *int32
	scopes       []string
	admin        bool
	refreshToken string
}

type defaultPermissions struct {
	context     context.Context
	options     *ControllerOptions
	authAttempt *common.AuthAttempt
	unwrapped   *unwrappedPermissions
	anonymous   bool
}

func NewPermissions(ctx context.Context, options *ControllerOptions) Permissions {
	return &defaultPermissions{
		context:   ctx,
		options:   options,
		unwrapped: nil,
	}
}

func addAuthAttemptToContext(ctx context.Context, aa *common.AuthAttempt) context.Context {
	newCtx := context.WithValue(ctx, "authAttempt", aa)
	return newCtx
}

func getAuthAttempt(ctx context.Context) *common.AuthAttempt {
	if v, ok := ctx.Value("authAttempt").(*common.AuthAttempt); ok {
		return v
	}
	return nil
}

func addClaimsToContext(ctx context.Context, claims jwtgo.MapClaims) context.Context {
	newCtx := context.WithValue(ctx, "claims", claims)
	return newCtx
}

func getClaims(ctx context.Context) (jwtgo.MapClaims, bool) {
	if v, ok := ctx.Value("claims").(jwtgo.MapClaims); ok {
		return v, true
	}
	return nil, false
}

func (p *defaultPermissions) unauthorized(m string) error {
	if p.authAttempt == nil || p.authAttempt.Unauthorized == nil {
		return fmt.Errorf("unable to make unauthorized error (%v)", m)
	}
	return p.authAttempt.Unauthorized(m)
}

func (p *defaultPermissions) forbidden(m string) error {
	if p.authAttempt == nil || p.authAttempt.Forbidden == nil {
		return fmt.Errorf("unable to make forbidden error (%v)", m)
	}
	return p.authAttempt.Forbidden(m)
}

func (p *defaultPermissions) notFound(m string) error {
	if p.authAttempt == nil || p.authAttempt.NotFound == nil {
		return fmt.Errorf("unable to make not-found error (%v)", m)
	}
	return p.authAttempt.NotFound(m)
}

func (p *defaultPermissions) getClaims() (jwtgo.MapClaims, error) {
	if claims, ok := getClaims(p.context); ok {
		return claims, nil
	}

	if token := jwt.ContextJWT(p.context); token != nil {
		if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
			return claims, nil
		}
	}

	authAttempt := getAuthAttempt(p.context)
	if authAttempt == nil {
		return nil, fmt.Errorf("forbidden (no attempt)")
	}

	return nil, authAttempt.Unauthorized("unauthorized")
}

func (p *defaultPermissions) unwrap() error {
	if p.unwrapped != nil || p.anonymous {
		return nil
	}

	authAttempt := getAuthAttempt(p.context)
	p.authAttempt = authAttempt

	claims, err := p.getClaims()
	if err != nil {
		p.anonymous = true
		return nil
	}

	userID := int32(claims["sub"].(float64))
	scopesRaw, ok := claims["scopes"].([]interface{})
	if !ok {
		return authAttempt.Unauthorized("invalid scopes")
	}
	scopesArray := make([]string, len(scopesRaw))
	admin := false
	for _, scope := range scopesRaw {
		scopesArray = append(scopesArray, scope.(string))
		if scope == "api:admin" {
			admin = true
		}
	}

	p.unwrapped = &unwrappedPermissions{
		userID:       &userID,
		scopes:       scopesArray,
		admin:        admin,
		refreshToken: claims["refresh_token"].(string),
	}

	return nil
}

func (p *defaultPermissions) RefreshToken() string {
	return p.unwrapped.refreshToken
}

func (p *defaultPermissions) Anonymous() bool {
	return p.anonymous
}

func (p *defaultPermissions) MaybeUserID() *int32 {
	if p.anonymous {
		return nil
	}
	return p.unwrapped.userID
}

func (p *defaultPermissions) UserID() int32 {
	return *p.unwrapped.userID
}

func (p *defaultPermissions) IsAdmin() bool {
	if p.anonymous {
		return false
	}
	return p.unwrapped.admin
}

func (p *defaultPermissions) RequireAdmin() error {
	if p.IsAdmin() {
		return nil
	}
	return p.forbidden("forbidden")
}

func (p *defaultPermissions) Unwrap() (Permissions, error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *defaultPermissions) ForProject(project *data.Project) (permissions ProjectPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	var projectUser *data.ProjectUser = nil
	if !p.Anonymous() {
		pr := repositories.NewProjectRepository(p.options.Database)
		if found, err := pr.QueryProjectUser(p.context, p.UserID(), project.ID); err != nil {
			return nil, err
		} else if found != nil {
			projectUser = found
		}
	}

	permissions = &projectPermissions{
		defaultPermissions: *p,
		project:            project,
		projectUser:        projectUser,
	}

	return
}

func (p *defaultPermissions) ForProjectByID(id int32) (permissions ProjectPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	pr := repositories.NewProjectRepository(p.options.Database)
	project, err := pr.QueryByID(p.context, id)
	if err != nil {
		return nil, err
	}

	return p.ForProject(project)
}

func (p *defaultPermissions) ForStationByID(id int) (permissions StationPermissions, err error) {
	sr := repositories.NewStationRepository(p.options.Database)
	pr := repositories.NewProjectRepository(p.options.Database)

	if err := p.unwrap(); err != nil {
		return nil, err
	}

	station, err := sr.QueryStationByID(p.context, int32(id))
	if err != nil {
		return nil, p.notFound(fmt.Sprintf("station not found: %v", err))
	}

	stationProjects, err := pr.QueryProjectsByStationIDForPermissions(p.context, station.ID)
	if err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
		stationProjects:    stationProjects,
	}

	return
}

func (p *defaultPermissions) ForStationByDeviceID(id []byte) (permissions StationPermissions, err error) {
	sr := repositories.NewStationRepository(p.options.Database)
	pr := repositories.NewProjectRepository(p.options.Database)

	if err := p.unwrap(); err != nil {
		return nil, err
	}

	station, err := sr.QueryStationByDeviceID(p.context, id)
	if err != nil {
		return nil, p.notFound(fmt.Sprintf("station not found: %v", err))
	}

	stationProjects, err := pr.QueryProjectsByStationIDForPermissions(p.context, station.ID)
	if err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
		stationProjects:    stationProjects,
	}

	return
}

func (p *defaultPermissions) ForStation(station *data.Station) (permissions StationPermissions, err error) {
	pr := repositories.NewProjectRepository(p.options.Database)

	if err := p.unwrap(); err != nil {
		return nil, err
	}

	stationProjects, err := pr.QueryProjectsByStationIDForPermissions(p.context, station.ID)
	if err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
		stationProjects:    stationProjects,
	}

	return
}

func (p *defaultPermissions) ForDiscussions(discussion *data.DiscussionPost) (permissions DiscussionPermissions, err error) {
	pr := repositories.NewProjectRepository(p.options.Database)

	if err := p.unwrap(); err != nil {
		return nil, err
	}

	var projects []*data.Project

	if discussion.ProjectID != nil {
		project, err := pr.QueryByID(p.context, *discussion.ProjectID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if len(discussion.StationIDs) > 0 {

		for _, stationID := range discussion.StationIDs {
			project, err := pr.QueryProjectsByStationIDForPermissions(p.context, int32(stationID))
			if err != nil {
				return nil, err
			}
			projects = append(projects, project...)
		}
	}

	permissions = &discussionPermissions{
		defaultPermissions: *p,
		discussion:         discussion,
		discussionProjects: projects,
	}

	return
}

type stationPermissions struct {
	defaultPermissions
	station         *data.Station
	stationProjects []*data.Project
}

func (p *stationPermissions) Station() *data.Station {
	return p.station
}

func (p *stationPermissions) CanView() error {
	if p.IsAdmin() {
		return nil
	}

	// Owners can always see their stations.
	if !p.Anonymous() {
		if p.station.OwnerID == p.UserID() {
			return nil
		}
	}

	// We don't know until we check the projects the station is a part
	// of, one of them has to be public.
	for _, project := range p.stationProjects {
		projectPermissions, err := NewPermissions(p.context, p.options).ForProjectByID(project.ID)
		if err != nil {
			return err
		}

		if err := projectPermissions.CanView(); err == nil {
			return nil
		}
	}

	return p.forbidden("forbidden")
}

func (p *stationPermissions) CanModify() error {
	// Administrators can do anything.
	if p.IsAdmin() {
		return nil
	}

	// Owners can always modify their stations.
	if !p.Anonymous() {
		if p.station.OwnerID == p.UserID() {
			return nil
		}
	} else {
		// Anonymous can't modify anything.
		return p.forbidden("forbidden")
	}

	// Check each of the station's projects to see if we have write permissions.
	for _, project := range p.stationProjects {
		projectPermissions, err := NewPermissions(p.context, p.options).ForProjectByID(project.ID)
		if err != nil {
			return err
		}

		if err := projectPermissions.CanModify(); err == nil {
			return nil
		}
	}

	return p.forbidden("forbidden")
}

func (p *stationPermissions) IsReadOnly() (bool, error) {
	if p.anonymous {
		return true, nil
	}

	if p.station.OwnerID == p.UserID() {
		return false, nil
	}

	pr := repositories.NewProjectRepository(p.options.Database)
	relationships, err := pr.QueryUserProjectRelationships(p.context, p.UserID())
	if err != nil {
		return true, err
	}

	for _, stationProject := range p.stationProjects {
		if rel, ok := relationships[stationProject.ID]; ok {
			if rel.CanModify() {
				return false, nil
			}
		}
	}

	return true, nil
}

type projectPermissions struct {
	defaultPermissions
	project     *data.Project
	projectUser *data.ProjectUser
}

func (p *projectPermissions) Project() *data.Project {
	return p.project
}

func (p *projectPermissions) CanView() error {
	if p.IsAdmin() {
		return nil
	}

	if p.project.Privacy == data.Public {
		return nil
	}

	if p.projectUser == nil {
		return p.forbidden("forbidden")
	}

	role := p.projectUser.LookupRole()
	if role == nil {
		return p.forbidden("forbidden")
	}

	return nil
}

func (p *projectPermissions) CanModify() error {
	if p.IsAdmin() {
		return nil
	}

	if p.projectUser == nil {
		return p.forbidden("forbidden")
	}

	role := p.projectUser.LookupRole()
	if role.IsProjectReadOnly() {
		return p.forbidden("forbidden")
	}

	return nil
}

type discussionPermissions struct {
	defaultPermissions
	discussion         *data.DiscussionPost
	discussionProjects []*data.Project
}

func (p *discussionPermissions) CanDelete() error {
	if p.IsAdmin() {
		return nil
	}

	if p.discussion.UserID == p.UserID() {
		return nil
	}

	for _, project := range p.discussionProjects {
		if !p.Anonymous() {
			pr := repositories.NewProjectRepository(p.options.Database)
			if projectUser, err := pr.QueryProjectUser(p.context, p.UserID(), project.ID); err != nil {
				return p.forbidden("forbidden")
			} else if projectUser != nil {
				if !projectUser.LookupRole().IsProjectAdministrator() {
					return p.forbidden("forbidden")
				}
				return nil
			}
		}
	}

	return p.forbidden("forbidden")
}

func (p *discussionPermissions) Discussion() *data.DiscussionPost {
	return p.discussion
}
