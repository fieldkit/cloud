package api

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"goa.design/goa/v3/security"

	project "github.com/fieldkit/cloud/server/api/gen/project"

	"github.com/fieldkit/cloud/server/backend/repositories"
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

func (c *ProjectService) AddUpdate(ctx context.Context, payload *project.AddUpdatePayload) (pu *project.ProjectUpdate, err error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(int(payload.ProjectID))
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	update := data.ProjectUpdate{
		ProjectActivity: data.ProjectActivity{
			CreatedAt: time.Now(),
			ProjectID: payload.ProjectID,
		},
		AuthorID: p.UserID(),
		Body:     payload.Body,
	}

	if err = c.options.Database.NamedGetContext(ctx, &update, `
		INSERT INTO fieldkit.project_update (created_at, project_id, author_id, body) VALUES (:created_at, :project_id, :author_id, :body) RETURNING *
		`, &update); err != nil {
		return nil, err
	}

	pu = &project.ProjectUpdate{
		ID:        update.ID,
		CreatedAt: update.CreatedAt.Unix() * 1000,
		Body:      update.Body,
	}

	return
}

func (c *ProjectService) ModifyUpdate(ctx context.Context, payload *project.ModifyUpdatePayload) (pu *project.ProjectUpdate, err error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(int(payload.ProjectID))
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	update := data.ProjectUpdate{
		ProjectActivity: data.ProjectActivity{
			ID:        payload.UpdateID,
			ProjectID: payload.ProjectID,
		},
		Body: payload.Body,
	}

	if _, err = c.options.Database.NamedExecContext(ctx, `
		UPDATE fieldkit.project_update SET body = :body WHERE id = :id
		`, &update); err != nil {
		return nil, err
	}

	pu = &project.ProjectUpdate{
		ID:        update.ID,
		CreatedAt: update.CreatedAt.Unix() * 1000,
		Body:      update.Body,
	}

	return
}

func (c *ProjectService) DeleteUpdate(ctx context.Context, payload *project.DeleteUpdatePayload) (err error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(int(payload.ProjectID))
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err = c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_update WHERE id = $1
		`, payload.UpdateID); err != nil {
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
		SELECT pi.* FROM fieldkit.project_invite AS pi JOIN fieldkit.user AS u ON (pi.invited_email = u.email) WHERE u.id = $1 AND (pi.accepted_time IS NULL AND pi.rejected_time IS NULL)
		`, p.UserID()); err != nil {
		return nil, err
	}

	pending := make([]*project.PendingInvite, 0)
	for _, i := range all {
		p := projectsByID[i.ProjectID]
		pending = append(pending, &project.PendingInvite{
			ID:   int64(i.ID),
			Time: i.InvitedTime.Unix() * 1000,
			Role: i.RoleID,
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
			log.Infow("accept failed, no token")
			return project.Unauthorized("permission denied")
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			log.Infow("accept failed, email mismatch, no token")
			return project.Unauthorized("permission denied")
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.project_invite SET accepted_time = NOW() WHERE id = $1`, payload.ID); err != nil {
		return err
	}

	role := data.LookupRole(invite.RoleID)
	if role == nil {
		role = data.MemberRole
	}

	if _, err := c.options.Database.ExecContext(ctx, `INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)`, invite.ProjectID, p.UserID(), role.ID); err != nil {
		return err
	}

	log.Infow("accepting", "invite_id", invite.ID)

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
			log.Infow("decline failed, email mismatch, no token")
			return project.Unauthorized("permission denied")
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			log.Infow("decline failed, token mismatch")
			return project.Unauthorized("permission denied")
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.project_invite SET rejected_time = NOW() WHERE id = $1`, payload.ID); err != nil {
		return err
	}

	log.Infow("rejecting", "invite_id", invite.ID)

	return nil
}

func (s *ProjectService) UploadMedia(ctx context.Context, payload *project.UploadMediaPayload, body io.ReadCloser) error {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)
	saved, err := mr.SaveFromService(ctx, body, payload.ContentLength)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	log.Infow("media", "project_id", payload.ProjectID, "content_type", saved.MimeType, "user_id", p.UserID())

	project := &data.Project{}
	if err := s.options.Database.GetContext(ctx, project, `
		UPDATE fieldkit.project SET media_url = $1, media_content_type = $2 WHERE id = $3 RETURNING *
		`, saved.URL, saved.MimeType, payload.ProjectID); err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) DownloadMedia(ctx context.Context, payload *project.DownloadMediaPayload) (*project.DownloadMediaResult, io.ReadCloser, error) {
	resource := &data.Project{}
	if err := s.options.Database.GetContext(ctx, resource, `
		SELECT media_url, media_content_type FROM fieldkit.project WHERE id = $1
		`, payload.ProjectID); err != nil {
		return nil, nil, err
	}

	if resource.MediaURL == nil || resource.MediaContentType == nil {
		return nil, nil, project.NotFound("not found")
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, *resource.MediaURL)
	if err != nil {
		return nil, nil, project.NotFound("not found")
	}

	return &project.DownloadMediaResult{
		Length:      int64(lm.Size),
		ContentType: *resource.MediaContentType,
	}, ioutil.NopCloser(lm.Reader), nil
}

func (s *ProjectService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return project.NotFound(m) },
		Unauthorized: func(m string) error { return project.Unauthorized(m) },
		Forbidden:    func(m string) error { return project.Forbidden(m) },
	})
}
