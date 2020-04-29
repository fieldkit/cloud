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

	if _, err = c.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.project_update (created_at, project_id, author_id, body) VALUES (:created_at, :project_id, :author_id, :body)
		`, &update); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) Invites(ctx context.Context, payload *project.InvitesPayload) (invites *project.PendingInvites, err error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT * FROM fieldkit.project WHERE id IN (
			SELECT pi.project_id FROM fieldkit.project_invite AS pi JOIN fieldkit.user AS u ON (pi.invited_email = u.email) WHERE u.id = $1
		)
		`, p.UserID()); err != nil {
		return nil, err
	}

	projectsByID := make(map[int32]*data.Project)
	for _, p := range projects {
		projectsByID[p.ID] = p
	}

	all := []*data.ProjectInvite{}
	if err := c.options.Database.SelectContext(ctx, &all, `
		SELECT pi.* FROM fieldkit.project_invite AS pi JOIN fieldkit.user AS u ON (pi.invited_email = u.email) WHERE u.id = $1
		`, p.UserID()); err != nil {
		return nil, err
	}

	pending := make([]*project.PendingInvite, 0)
	for _, i := range all {
		p := projectsByID[i.ProjectID]
		pending = append(pending, &project.PendingInvite{
			ID:   int64(i.ID),
			Time: i.InvitedTime.Unix() * 1000,
			Project: &project.ProjectSummary{
				ID:   int64(p.ID),
				Name: p.Name,
			},
		})
	}

	invites = &project.PendingInvites{
		Pending: pending,
	}

	return
}

func (c *ProjectService) LookupInvite(ctx context.Context, payload *project.LookupInvitePayload) (invites *project.PendingInvites, err error) {
	_, err = NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	inviteToken := &data.ValidationToken{}
	if err := inviteToken.Token.UnmarshalText([]byte(payload.Token)); err != nil {
		return nil, err
	}

	all := []*data.ProjectInvite{}
	if err := c.options.Database.SelectContext(ctx, &all, `
		SELECT pi.* FROM fieldkit.project_invite AS pi WHERE pi.token = $1
		`, inviteToken.Token); err != nil {
		return nil, err
	}

	if len(all) != 1 {
		return nil, project.NotFound("invalid invite token")
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT * FROM fieldkit.project WHERE id = $1
		`, all[0].ProjectID); err != nil {
		return nil, err
	}

	projectsByID := make(map[int32]*data.Project)
	for _, p := range projects {
		projectsByID[p.ID] = p
	}

	pending := make([]*project.PendingInvite, 0)
	for _, i := range all {
		p := projectsByID[i.ProjectID]
		pending = append(pending, &project.PendingInvite{
			ID:   int64(i.ID),
			Time: i.InvitedTime.Unix() * 1000,
			Project: &project.ProjectSummary{
				ID:   int64(p.ID),
				Name: p.Name,
			},
		})
	}

	invites = &project.PendingInvites{
		Pending: pending,
	}

	return
}

func (c *ProjectService) AcceptInvite(ctx context.Context, payload *project.AcceptInvitePayload) (err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1`, p.UserID()); err != nil {
		return err
	}

	invite := &data.ProjectInvite{}
	if err := c.options.Database.GetContext(ctx, invite, `SELECT pi.* FROM fieldkit.project_invite AS pi WHERE pi.id = $1`, payload.ID); err != nil {
		return err
	}

	if user.Email != invite.InvitedEmail {
		if payload.Token == nil {
			return project.Unauthorized("permission denied")
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			return project.Unauthorized("permission denied")
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.project_invite SET accepted_time = NOW() WHERE id = $1`, payload.ID); err != nil {
		return err
	}

	role := data.MemberRole
	if _, err := c.options.Database.ExecContext(ctx, `INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)`, invite.ProjectID, p.UserID(), role.ID); err != nil {
		return err
	}

	log.Infow("accepting", "invite", invite)

	_ = p

	return nil
}

func (c *ProjectService) RejectInvite(ctx context.Context, payload *project.RejectInvitePayload) (err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1`, p.UserID()); err != nil {
		return err
	}

	invite := &data.ProjectInvite{}
	if err := c.options.Database.GetContext(ctx, invite, `SELECT pi.* FROM fieldkit.project_invite AS pi WHERE pi.id = $1`, payload.ID); err != nil {
		return err
	}

	if user.Email != invite.InvitedEmail {
		if payload.Token == nil {
			return project.Unauthorized("permission denied")
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			return project.Unauthorized("permission denied")
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.project_invite SET rejected_time = NOW() WHERE id = $1`, payload.ID); err != nil {
		return err
	}

	log.Infow("rejecting", "invite", invite)

	_ = p

	return nil
}

func (s *ProjectService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  project.Unauthorized("invalid token"),
		InvalidScopes: project.Unauthorized("invalid scopes in token"),
	})
}
