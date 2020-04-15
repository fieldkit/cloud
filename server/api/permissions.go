package api

import (
	"context"
	"database/sql"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa"

	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/data"
)

type ProjectPermissions interface {
	Permissions
	CanView() error
	CanModify() error
}

type StationPermissions interface {
	Permissions
	CanView() error
	CanModify() error
	IsReadOnly() bool
}

type Permissions interface {
	Unwrap() (permissions Permissions, err error)
	UserID() int32
	ForProjectByID(id int) (permissions ProjectPermissions, err error)
	ForStationByID(id int) (permissions StationPermissions, err error)
	ForStationByDeviceID(id []byte) (permissions StationPermissions, err error)
	ForStation(station *data.Station) (permissions StationPermissions, err error)
}

type unwrappedPermissions struct {
	userID int32
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

func (p *defaultPermissions) unwrap() error {
	if p.unwrapped != nil {
		return nil
	}

	token := jwt.ContextJWT(p.context)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error")
	}

	userID := int32(claims["sub"].(float64))

	p.unwrapped = &unwrappedPermissions{
		userID: userID,
	}

	return nil
}

func (p *defaultPermissions) UserID() int32 {
	return p.unwrapped.userID
}

func (p *defaultPermissions) Unwrap() (Permissions, error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *defaultPermissions) ForProjectByID(id int) (permissions ProjectPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	project := &data.Project{}
	if err := p.options.Database.GetContext(p.context, project, "SELECT p.* FROM fieldkit.project AS p WHERE p.id = $1", id); err != nil {
		return nil, err
	}

	projectUser := &data.ProjectUser{}
	if err := p.options.Database.GetContext(p.context, projectUser, "SELECT p.* FROM fieldkit.project_user AS p WHERE p.user_id = $1", p.UserID()); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		projectUser = nil
	}

	permissions = &projectPermissions{
		defaultPermissions: *p,
		Project:            project,
		ProjectUser:        projectUser,
	}

	return
}

func (p *defaultPermissions) ForStationByID(id int) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	station := &data.Station{}
	if err := p.options.Database.GetContext(p.context, station, "SELECT s.* FROM fieldkit.station AS s WHERE s.id = $1", id); err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		Station:            station,
	}

	return
}

func (p *defaultPermissions) ForStationByDeviceID(id []byte) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	station := &data.Station{}
	if err := p.options.Database.GetContext(p.context, station, "SELECT s.* FROM fieldkit.station AS s WHERE s.device_id = $1", id); err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		Station:            station,
	}

	return
}

func (p *defaultPermissions) ForStation(station *data.Station) (permissions StationPermissions, err error) {
	if err := p.unwrap(); err != nil {
		return nil, err
	}

	permissions = &stationPermissions{
		defaultPermissions: *p,
		Station:            station,
	}

	return
}

type stationPermissions struct {
	defaultPermissions
	Station *data.Station
}

func (p *stationPermissions) CanView() error {
	return nil
}

func (p *stationPermissions) CanModify() error {
	if p.Station.OwnerID != p.UserID() {
		return goa.ErrUnauthorized(fmt.Sprintf("unauthorized"))
	}
	return nil
}

func (p *stationPermissions) IsReadOnly() bool {
	return p.Station.OwnerID != p.UserID()
}

type projectPermissions struct {
	defaultPermissions
	Project     *data.Project
	ProjectUser *data.ProjectUser
}

func (p *projectPermissions) CanView() error {
	return nil
}

func (p *projectPermissions) CanModify() error {
	if p.ProjectUser == nil {
		return goa.ErrUnauthorized(fmt.Sprintf("unauthorized"))
	}

	role := p.ProjectUser.LookupRole()
	if role.IsProjectReadOnly() {
		return goa.ErrUnauthorized(fmt.Sprintf("unauthorized"))
	}

	return nil
}
