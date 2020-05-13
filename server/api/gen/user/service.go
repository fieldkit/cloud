// Code generated by goa v3.1.2, DO NOT EDIT.
//
// user service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package user

import (
	"context"

	userviews "github.com/fieldkit/cloud/server/api/gen/user/views"
	"goa.design/goa/v3/security"
)

// Service is the user service interface.
type Service interface {
	// Roles implements roles.
	Roles(context.Context, *RolesPayload) (res *AvailableRoles, err error)
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
var MethodNames = [1]string{"roles"}

// RolesPayload is the payload type of the user service roles method.
type RolesPayload struct {
	Auth string
}

// AvailableRoles is the result type of the user service roles method.
type AvailableRoles struct {
	Roles []*AvailableRole
}

type AvailableRole struct {
	ID   int32
	Name string
}

// unauthorized
type Unauthorized string

// forbidden
type Forbidden string

// not-found
type NotFound string

// bad-request
type BadRequest string

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "unauthorized"
}

// ErrorName returns "unauthorized".
func (e Unauthorized) ErrorName() string {
	return "unauthorized"
}

// Error returns an error description.
func (e Forbidden) Error() string {
	return "forbidden"
}

// ErrorName returns "forbidden".
func (e Forbidden) ErrorName() string {
	return "forbidden"
}

// Error returns an error description.
func (e NotFound) Error() string {
	return "not-found"
}

// ErrorName returns "not-found".
func (e NotFound) ErrorName() string {
	return "not-found"
}

// Error returns an error description.
func (e BadRequest) Error() string {
	return "bad-request"
}

// ErrorName returns "bad-request".
func (e BadRequest) ErrorName() string {
	return "bad-request"
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
