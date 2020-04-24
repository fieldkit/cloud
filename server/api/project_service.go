package api

import (
	"context"
	"time"

	"goa.design/goa/v3/security"

	project "github.com/fieldkit/cloud/server/api/gen/project"

	"github.com/fieldkit/cloud/server/data"
)

type ProjectService struct {
	options *ControllerOptions
}

func NewProjectService(ctx context.Context, options *ControllerOptions) *ProjectService {
	return &ProjectService{
		options: options,
	}
}

func (c *ProjectService) Update(ctx context.Context, payload *project.UpdatePayload) (err error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(int(payload.ID))
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	update := data.ProjectUpdate{
		ProjectActivity: data.ProjectActivity{
			CreatedAt: time.Now(),
			ProjectID: payload.ID,
		},
		AuthorID: int64(p.UserID()),
		Body:     payload.Body,
	}

	if _, err = c.options.Database.NamedExecContext(ctx, `INSERT INTO fieldkit.project_update (created_at, project_id, author_id, body) VALUES (:created_at, :project_id, :author_id, :body)`, &update); err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  nil,
		InvalidScopes: nil,
	})
}
