package api

import (
	"context"
	"io"
	"io/ioutil"

	"goa.design/goa/v3/security"

	user "github.com/fieldkit/cloud/server/api/gen/user"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type UserService struct {
	options *ControllerOptions
}

func NewUserService(ctx context.Context, options *ControllerOptions) *UserService {
	return &UserService{options: options}
}

func (s *UserService) Roles(ctx context.Context, payload *user.RolesPayload) (*user.AvailableRoles, error) {
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

func (s *UserService) Delete(ctx context.Context, payload *user.DeletePayload) error {
	r, err := repositories.NewUserRepository(s.options.Database)
	if err != nil {
		return err
	}
	return r.Delete(ctx, payload.UserID)
}

func (s *UserService) UploadPhoto(ctx context.Context, payload *user.UploadPhotoPayload, body io.ReadCloser) error {
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

	log.Infow("media", "content_type", saved.MimeType, "user_id", p.UserID())

	user := &data.User{}
	if err := s.options.Database.GetContext(ctx, user, `
		UPDATE fieldkit.user SET media_url = $1, media_content_type = $2 WHERE id = $3 RETURNING *
		`, saved.URL, saved.MimeType, p.UserID()); err != nil {
		return err
	}

	return nil
}

func (s *UserService) DownloadPhoto(ctx context.Context, payload *user.DownloadPhotoPayload) (*user.DownloadPhotoResult, io.ReadCloser, error) {
	resource := &data.User{}
	if err := s.options.Database.GetContext(ctx, resource, `SELECT * FROM fieldkit.user WHERE id = $1`, payload.UserID); err != nil {
		return nil, nil, err
	}

	if resource.MediaURL == nil || resource.MediaContentType == nil {
		return nil, nil, user.NotFound("not found")
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, *resource.MediaURL)
	if err != nil {
		return nil, nil, user.NotFound("not found")
	}

	return &user.DownloadPhotoResult{
		Length:      int64(lm.Size),
		ContentType: *resource.MediaContentType,
	}, ioutil.NopCloser(lm.Reader), nil
}

func (s *UserService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return user.NotFound(m) },
		Unauthorized: func(m string) error { return user.Unauthorized(m) },
		Forbidden:    func(m string) error { return user.Forbidden(m) },
	})
}
