package api

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/data"
)

type Permissions struct {
	options *ControllerOptions
	UserID  int32
}

type ProjectPermissions struct {
	Project *data.Project
}

func NewPermissions(ctx context.Context, options *ControllerOptions) (p *Permissions, err error) {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return nil, fmt.Errorf("JWT token is missing from context")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, fmt.Errorf("JWT claims error")
	}

	p = &Permissions{
		options: options,
		UserID:  int32(claims["sub"].(float64)),
	}

	return
}

func (p *Permissions) IsStationReadOnly(station *data.Station) bool {
	return station.OwnerID != p.UserID
}

func (p *Permissions) CanModifyStationByStationID(stationId int) error {
	return nil
}

func (p *Permissions) CanViewStationByStationID(stationId int) error {
	return nil
}

func (p *Permissions) CanModifyStationByDeviceID(deviceId []byte) error {
	return nil
}

func (p *Permissions) CanViewStationByDeviceID(deviceId []byte) error {
	return nil
}

func (p *Permissions) CanModifyProject(projectId int) error {
	return nil
}

func (p *Permissions) CanViewProject(projectId int) error {
	return nil
}

func (p *Permissions) ForProject(id int) (permissions *ProjectPermissions, err error) {
	return
}
