package api

import (
	"context"
	"database/sql"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/backend/repositories"
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
	IsReadOnly() bool
}

type Permissions interface {
	Unwrap() (permissions Permissions, err error)
	UserID() int32
	RefreshToken() string
	RequireAdmin() error
	IsAdmin() bool
	ForProjectByID(id int32) (permissions ProjectPermissions, err error)
	ForStationByID(id int) (permissions StationPermissions, err error)
	ForStationByDeviceID(id []byte) (permissions StationPermissions, err error)
	ForStation(station *data.Station) (permissions StationPermissions, err error)
}

type unwrappedPermissions struct {
	userID       int32
	authAttempt  *AuthAttempt
	scopes       []string
	admin        bool
	refreshToken string
}

type defaultPermissions struct {
	context   context.Context
	options   *ControllerOptions
	unwrapped *unwrappedPermissions
}

func NewPermissions(ctx context.Context, options *ControllerOptions) Permissions {
	return &defaultPermissions{
		context:   ctx,
		options:   options,
		unwrapped: nil,
	}
}

func addAuthAttemptToContext(ctx context.Context, aa *AuthAttempt) context.Context {
	newCtx := context.WithValue(ctx, "authAttempt", aa)
	return newCtx
}

func getAuthAttempt(ctx context.Context) *AuthAttempt {
	if v, ok := ctx.Value("authAttempt").(*AuthAttempt); ok {
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
	if p.unwrapped == nil || p.unwrapped.authAttempt == nil || p.unwrapped.authAttempt.Unauthorized == nil {
		return fmt.Errorf("unable to make unauthorized error (%v)", m)
	}
	return p.unwrapped.authAttempt.Unauthorized(m)
}

func (p *defaultPermissions) forbidden(m string) error {
	if p.unwrapped == nil || p.unwrapped.authAttempt == nil || p.unwrapped.authAttempt.Forbidden == nil {
		return fmt.Errorf("unable to make forbidden error (%v)", m)
	}
	return p.unwrapped.authAttempt.Forbidden(m)
}

func (p *defaultPermissions) notFound(m string) error {
	if p.unwrapped == nil || p.unwrapped.authAttempt == nil || p.unwrapped.authAttempt.NotFound == nil {
		return fmt.Errorf("unable to make not-found error (%v)", m)
	}
	return p.unwrapped.authAttempt.NotFound(m)
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
		return nil, fmt.Errorf("invalid auth attempted")
	}

	return nil, authAttempt.Unauthorized("unauthorized")
}

func (p *defaultPermissions) unwrap() error {
	if p.unwrapped != nil {
		return nil
	}

	claims, err := p.getClaims()
	if err != nil {
		return err
	}

	authAttempt := getAuthAttempt(p.context)

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
		authAttempt:  authAttempt,
		userID:       userID,
		scopes:       scopesArray,
		admin:        admin,
		refreshToken: claims["refresh_token"].(string),
	}

	return nil
}

func (p *defaultPermissions) RefreshToken() string {
	return p.unwrapped.refreshToken
}

func (p *defaultPermissions) UserID() int32 {
	return p.unwrapped.userID
}

func (p *defaultPermissions) IsAdmin() bool {
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

func (p *defaultPermissions) ForProjectByID(id int32) (permissions ProjectPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	project := &data.Project{}
	if err := p.options.Database.GetContext(p.context, project, "SELECT p.* FROM fieldkit.project AS p WHERE p.id = $1", id); err != nil {
		return nil, p.notFound(fmt.Sprintf("project not found: %v", err))
	}

	projectUser := &data.ProjectUser{}
	if err := p.options.Database.GetContext(p.context, projectUser, "SELECT p.* FROM fieldkit.project_user AS p WHERE p.user_id = $1 AND p.project_id = $2", p.UserID(), id); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		projectUser = nil
	}

	permissions = &projectPermissions{
		defaultPermissions: *p,
		project:            project,
		projectUser:        projectUser,
	}

	return
}

func (p *defaultPermissions) ForStationByID(id int) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	r, err := repositories.NewStationRepository(p.options.Database)
	if err != nil {
		return nil, err
	}

	station, err := r.QueryStationByID(p.context, int32(id))
	if err != nil {
		return nil, p.notFound(fmt.Sprintf("station not found"))
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
	}

	return
}

func (p *defaultPermissions) ForStationByDeviceID(id []byte) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	r, err := repositories.NewStationRepository(p.options.Database)
	if err != nil {
		return nil, err
	}

	station, err := r.QueryStationByDeviceID(p.context, id)
	if err != nil {
		return nil, p.notFound(fmt.Sprintf("station not found"))
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
	}

	return
}

func (p *defaultPermissions) ForStation(station *data.Station) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		station:            station,
	}

	return
}

type stationPermissions struct {
	defaultPermissions
	station *data.Station
}

func (p *stationPermissions) Station() *data.Station {
	return p.station
}

func (p *stationPermissions) CanView() error {
	return nil
}

func (p *stationPermissions) CanModify() error {
	if p.IsAdmin() {
		return nil
	}

	if p.station.OwnerID != p.UserID() {
		return p.forbidden("forbidden")
	}

	return nil
}

func (p *stationPermissions) IsReadOnly() bool {
	return p.station.OwnerID != p.UserID()
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
