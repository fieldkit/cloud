package api

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
)

var (
	ErrUnauthorized       error = tasks.Unauthorized("invalid username and password combination")
	ErrInvalidToken       error = tasks.Unauthorized("invalid token")
	ErrInvalidTokenScopes error = tasks.Unauthorized("invalid scopes in token")
)

type TasksService struct {
	options *ControllerOptions
}

func NewTasksService(ctx context.Context, options *ControllerOptions) *TasksService {
	return &TasksService{
		options: options,
	}
}

func (s *TasksService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return s.options.JWTHMACKey, nil
	})
	if err != nil {
		return ctx, ErrInvalidToken
	}
	if claims["scopes"] == nil {
		return ctx, ErrInvalidTokenScopes
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return ctx, ErrInvalidTokenScopes
	}
	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}
	if err := scheme.Validate(scopesInToken); err != nil {
		return ctx, tasks.InvalidScopes(err.Error())
	}
	return ctx, nil
}

func (c *TasksService) Five(ctx context.Context) error {
	return nil
}

func (c *TasksService) RefreshDevice(ctx context.Context, payload *tasks.RefreshDevicePayload) error {
	return nil
}
