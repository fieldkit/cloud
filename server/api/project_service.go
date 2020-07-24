package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"goa.design/goa/v3/security"

	project "github.com/fieldkit/cloud/server/api/gen/project"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ProjectService struct {
	projects *repositories.ProjectRepository
	options  *ControllerOptions
}

func NewProjectService(ctx context.Context, options *ControllerOptions) *ProjectService {
	return &ProjectService{
		projects: repositories.NewProjectRepository(options.Database),
		options:  options,
	}
}

func (c *ProjectService) Add(ctx context.Context, payload *project.AddPayload) (*project.Project, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	goal := ""
	if payload.Project.Goal != nil {
		goal = *payload.Project.Goal
	}

	location := ""
	if payload.Project.Location != nil {
		location = *payload.Project.Location
	}

	tags := ""
	if payload.Project.Tags != nil {
		tags = *payload.Project.Tags
	}

	private := true
	if payload.Project.Private != nil {
		private = *payload.Project.Private
	}

	newProject := &data.Project{
		Name:        payload.Project.Name,
		Description: payload.Project.Description,
		Goal:        goal,
		Location:    location,
		Tags:        tags,
		Private:     private,
	}

	if start, err := tryParseDate(payload.Project.StartTime); err == nil {
		newProject.StartTime = &start
	}
	if end, err := tryParseDate(payload.Project.EndTime); err == nil {
		newProject.EndTime = &end
	}

	newProject, err = c.projects.AddProject(ctx, p.UserID(), newProject)
	if err != nil {
		return nil, err
	}

	relationship := &data.UserProjectRelationship{
		MemberRole: data.AdministratorRole.ID,
	}

	return ProjectType(c.options.signer, newProject, 0, relationship)
}

func (c *ProjectService) Update(ctx context.Context, payload *project.UpdatePayload) (*project.Project, error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	goal := ""
	if payload.Project.Goal != nil {
		goal = *payload.Project.Goal
	}

	location := ""
	if payload.Project.Location != nil {
		location = *payload.Project.Location
	}

	tags := ""
	if payload.Project.Location != nil {
		tags = *payload.Project.Tags
	}

	private := true
	if payload.Project.Private != nil {
		private = *payload.Project.Private
	}

	updating := &data.Project{
		ID:          payload.ProjectID,
		Name:        payload.Project.Name,
		Description: payload.Project.Description,
		Goal:        goal,
		Location:    location,
		Tags:        tags,
		Private:     private,
	}

	if start, err := tryParseDate(payload.Project.StartTime); err == nil {
		updating.StartTime = &start
	}
	if end, err := tryParseDate(payload.Project.EndTime); err == nil {
		updating.EndTime = &end
	}

	if err := c.options.Database.NamedGetContext(ctx, updating, `
		UPDATE fieldkit.project SET name = :name, description = :description, goal = :goal, location = :location,
		tags = :tags, private = :private, start_time = :start_time, end_time = :end_time WHERE id = :id RETURNING *`, updating); err != nil {
		return nil, err
	}

	relationship := &data.UserProjectRelationship{
		MemberRole: data.AdministratorRole.ID,
	}

	return ProjectType(c.options.signer, updating, 0, relationship)
}

func (c *ProjectService) AddStation(ctx context.Context, payload *project.AddStationPayload) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if err := c.projects.AddStationToProjectByID(ctx, payload.ProjectID, payload.StationID); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) RemoveStation(ctx context.Context, payload *project.RemoveStationPayload) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_station WHERE project_id = $1 AND station_id = $2
		`, payload.ProjectID, payload.StationID); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) Get(ctx context.Context, payload *project.GetPayload) (*project.Project, error) {
	relationships := make(map[int32]*data.UserProjectRelationship)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		relationships, err = c.projects.QueryUserProjectRelationships(ctx, p.UserID())
		if err != nil {
			return nil, err
		}
	}

	if relationships[payload.ProjectID] == nil {
		relationships[payload.ProjectID] = &data.UserProjectRelationship{}
	}

	getting := &data.Project{}
	if err := c.options.Database.GetContext(ctx, getting, `
		SELECT p.* FROM fieldkit.project AS p WHERE p.id = $1
		`, payload.ProjectID); err != nil {
		if err == sql.ErrNoRows {
			return nil, project.MakeNotFound(errors.New("not found"))
		}
		return nil, err
	}

	followerSummaries := []*data.FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followerSummaries, `
		SELECT f.project_id, COUNT(f.*) AS followers FROM fieldkit.project_follower AS f WHERE f.project_id IN ($1) GROUP BY f.project_id
		`, payload.ProjectID); err != nil {
		return nil, err
	}

	followers := int32(0)
	if len(followerSummaries) > 0 {
		followers = followerSummaries[0].Followers
	}

	return ProjectType(c.options.signer, getting, followers, relationships[getting.ID])
}

func (c *ProjectService) ListCommunity(ctx context.Context, payload *project.ListCommunityPayload) (*project.Projects, error) {
	relationships := make(map[int32]*data.UserProjectRelationship)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		relationships, err = c.projects.QueryUserProjectRelationships(ctx, p.UserID())
		if err != nil {
			return nil, err
		}
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT p.* FROM fieldkit.project AS p WHERE NOT p.private ORDER BY p.name LIMIT 10
		`); err != nil {
		return nil, err
	}

	followers := []*data.FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followers, `
		SELECT f.project_id, COUNT(f.*) AS followers FROM fieldkit.project_follower AS f WHERE f.project_id IN (
			SELECT id FROM fieldkit.project ORDER BY name LIMIT 10
		) GROUP BY f.project_id
		`); err != nil {
		return nil, err
	}

	return ProjectsType(c.options.signer, projects, followers, relationships)
}

func (c *ProjectService) ListMine(ctx context.Context, payload *project.ListMinePayload) (*project.Projects, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	relationships, err := c.projects.QueryUserProjectRelationships(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT * FROM fieldkit.project WHERE id IN (SELECT project_id FROM fieldkit.project_user WHERE user_id = $1) ORDER BY name
		`, p.UserID()); err != nil {
		return nil, err
	}

	followers := []*data.FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followers, `
		SELECT f.project_id, COUNT(f.*) AS followers FROM fieldkit.project_follower AS f WHERE f.project_id IN (
			SELECT project_id FROM fieldkit.project_user WHERE user_id = $1
		) GROUP BY f.project_id
		`, p.UserID()); err != nil {
		return nil, err
	}

	return ProjectsType(c.options.signer, projects, followers, relationships)
}

func (c *ProjectService) Invite(ctx context.Context, payload *project.InvitePayload) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	existing := int64(0)
	if err := c.options.Database.GetContext(ctx, &existing, `
		SELECT COUNT(*) FROM (
			SELECT pi.project_id, pi.invited_email AS email FROM fieldkit.project_invite AS pi
			UNION
			SELECT pu.project_id, u.email AS email FROM fieldkit.user AS u JOIN fieldkit.project_user AS pu ON (u.id = pu.user_id)
		) AS q WHERE q.project_id = $1 AND email = $2
		`, payload.ProjectID, payload.Invite.Email); err != nil {
		return err
	}
	if existing > 0 {
		return project.MakeBadRequest(errors.New("duplicate"))
	}

	token, err := data.NewToken(20)
	if err != nil {
		return err
	}

	invite := &data.ProjectInvite{
		ProjectID:    payload.ProjectID,
		UserID:       p.UserID(),
		InvitedTime:  time.Now(),
		InvitedEmail: payload.Invite.Email,
		RoleID:       payload.Invite.Role,
		Token:        token,
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		INSERT INTO fieldkit.project_invite (project_id, user_id, invited_email, invited_time, token, role_id) VALUES ($1, $2, $3, $4, $5, $6)
		`, invite.ProjectID, invite.UserID, invite.InvitedEmail, invite.InvitedTime, invite.Token, invite.RoleID); err != nil {
		return err
	}

	sender := &data.User{}
	if err := c.options.Database.GetContext(ctx, sender, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, p.UserID()); err != nil {
		return err
	}

	if err := c.options.Emailer.SendProjectInvitation(sender, invite); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) RemoveUser(ctx context.Context, payload *project.RemoveUserPayload) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	removing := &data.User{}
	if err := c.options.Database.GetContext(ctx, removing, `SELECT * FROM fieldkit.user WHERE email = $1`, payload.Remove.Email); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		if removing.ID == p.UserID() {
			return project.MakeBadRequest(errors.New("no removing yourself"))
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_invite WHERE project_id = $1 AND invited_email = $2
		`, payload.ProjectID, payload.Remove.Email); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_user WHERE project_id = $1 AND user_id IN (SELECT u.id FROM fieldkit.user AS u WHERE u.email = $2)
		`, payload.ProjectID, payload.Remove.Email); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) Delete(ctx context.Context, payload *project.DeletePayload) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, `SELECT media_url FROM fieldkit.project WHERE id = $1`, payload.ProjectID); err != nil {
		return err
	}
	if project.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.MediaFiles)

		err := mr.DeleteByURL(ctx, *project.MediaURL)
		if err != nil {
			return err
		}
	}

	if err := c.projects.Delete(ctx, payload.ProjectID); err != nil {
		return err
	}

	return nil
}

func (c *ProjectService) AddUpdate(ctx context.Context, payload *project.AddUpdatePayload) (pu *project.ProjectUpdate, err error) {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
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
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
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
	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ProjectID)
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
		return nil, project.MakeNotFound(errors.New("invalid invite token"))
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
			return project.MakeUnauthorized(errors.New("permission denied"))
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			log.Infow("accept failed, email mismatch, no token")
			return project.MakeUnauthorized(errors.New("permission denied"))
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
			return project.MakeUnauthorized(errors.New("permission denied"))
		}
		given := data.Token{}
		if err := given.UnmarshalText([]byte(*payload.Token)); err != nil {
			return err
		}
		if given.String() != invite.Token.String() {
			log.Infow("decline failed, token mismatch")
			return project.MakeUnauthorized(errors.New("permission denied"))
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.project_invite SET rejected_time = NOW() WHERE id = $1`, payload.ID); err != nil {
		return err
	}

	log.Infow("rejecting", "invite_id", invite.ID)

	return nil
}

func (s *ProjectService) UploadPhoto(ctx context.Context, payload *project.UploadPhotoPayload, body io.ReadCloser) error {
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

func (s *ProjectService) DownloadPhoto(ctx context.Context, payload *project.DownloadPhotoPayload) (*project.DownloadPhotoResult, io.ReadCloser, error) {
	resource := &data.Project{}
	if err := s.options.Database.GetContext(ctx, resource, `
		SELECT media_url, media_content_type FROM fieldkit.project WHERE id = $1
		`, payload.ProjectID); err != nil {
		return nil, nil, err
	}

	if resource.MediaURL == nil || resource.MediaContentType == nil {
		return nil, nil, project.MakeNotFound(errors.New("not found"))
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, *resource.MediaURL)
	if err != nil {
		return nil, nil, project.MakeNotFound(errors.New("not found"))
	}

	return &project.DownloadPhotoResult{
		Length:      int64(lm.Size),
		ContentType: *resource.MediaContentType,
	}, ioutil.NopCloser(lm.Reader), nil
}

func (s *ProjectService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return project.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return project.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return project.MakeForbidden(errors.New(m)) },
	})
}

func ProjectType(signer *Signer, dm *data.Project, numberOfFollowers int32, userRelationship *data.UserProjectRelationship) (*project.Project, error) {
	role := userRelationship.LookupRole()

	photo, err := signer.SignAndBustURL(fmt.Sprintf("/projects/%d/media", dm.ID), dm.MediaURL)
	if err != nil {
		return nil, err
	}

	wm := &project.Project{
		ID:          dm.ID,
		Name:        dm.Name,
		Description: dm.Description,
		Goal:        dm.Goal,
		Location:    dm.Location,
		Tags:        dm.Tags,
		Private:     dm.Private,
		Photo:       photo,
		ReadOnly:    role.IsProjectReadOnly(),
		Following: &project.ProjectFollowing{
			Total:     numberOfFollowers,
			Following: userRelationship.Following,
		},
	}

	if dm.StartTime != nil {
		startString := (*dm.StartTime).Format(time.RFC3339)
		wm.StartTime = &startString
	}

	if dm.EndTime != nil {
		endString := (*dm.EndTime).Format(time.RFC3339)
		wm.EndTime = &endString
	}

	return wm, nil
}

func findNumberOfFollowers(followers []*data.FollowersSummary, id int32) int32 {
	for _, f := range followers {
		if f.ProjectID == id {
			return f.Followers
		}
	}
	return 0
}

func ProjectsType(signer *Signer, projects []*data.Project, followers []*data.FollowersSummary, relationships map[int32]*data.UserProjectRelationship) (*project.Projects, error) {
	projectsCollection := make([]*project.Project, len(projects))
	for i, project := range projects {
		numberOfFollowers := findNumberOfFollowers(followers, project.ID)
		if rel, ok := relationships[project.ID]; ok {
			project, err := ProjectType(signer, project, numberOfFollowers, rel)
			if err != nil {
				return nil, err
			}
			projectsCollection[i] = project
		} else {
			project, err := ProjectType(signer, project, numberOfFollowers, &data.UserProjectRelationship{})
			if err != nil {
				return nil, err
			}
			projectsCollection[i] = project
		}
	}

	return &project.Projects{
		Projects: projectsCollection,
	}, nil
}
