package api

import (
	"context"

	"goa.design/goa/v3/security"

	user "github.com/fieldkit/cloud/server/api/gen/user"

	"github.com/fieldkit/cloud/server/data"
)

type UserService struct {
	options *ControllerOptions
}

func NewUserService(ctx context.Context, options *ControllerOptions) *UserService {
	return &UserService{
		options: options,
	}
}

func (ms *UserService) Roles(ctx context.Context, payload *user.RolesPayload) (*user.AvailableRoles, error) {
	roles := make([]*user.AvailableRole, 0)

	for _, r := range data.AvailableRoles {
		roles = append(roles, &user.AvailableRole{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &user.AvailableRoles{
		Roles: roles,
	}, nil
}

func (s *UserService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  user.Unauthorized("invalid token"),
		InvalidScopes: user.Unauthorized("invalid scopes in token"),
	})
}
