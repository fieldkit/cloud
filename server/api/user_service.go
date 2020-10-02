package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
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

func (s *UserService) Login(ctx context.Context, payload *user.LoginPayload) (*user.LoginResult, error) {
	now := time.Now()

	s.options.Metrics.AuthTry()

	authed, err := s.authenticateOrSpoof(ctx, payload.Login.Email, payload.Login.Password)
	if err == data.IncorrectPasswordError {
		return nil, user.MakeUnauthorized(errors.New("invalid email or password"))
	}
	if err != nil {
		return nil, err
	}
	if authed == nil {
		return nil, user.MakeUnauthorized(errors.New("invalid email or password"))
	}

	refreshToken, err := data.NewRefreshToken(authed.ID, 20, now.Add(time.Duration(72)*time.Hour))
	if err != nil {
		return nil, err
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.refresh_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, refreshToken); err != nil {
		return nil, user.MakeUnauthorized(errors.New("invalid email or password"))
	}

	token := authed.NewToken(now, refreshToken)
	signedToken, err := token.SignedString(s.options.JWTHMACKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	s.options.Metrics.AuthSuccess()

	return &user.LoginResult{
		Authorization: "Bearer " + signedToken,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, payload *user.LogoutPayload) error {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return err
	}

	if _, err := s.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.refresh_token WHERE token = $1
		`, p.RefreshToken()); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Add(ctx context.Context, payload *user.AddPayload) (*user.User, error) {
	log := Logger(ctx).Sugar()

	if yes, err := s.userExists(ctx, payload.User.Email); err != nil {
		return nil, err
	} else if yes {
		return nil, user.MakeBadRequest(errors.New("bad request"))
	}

	user := &data.User{
		Name:     data.Name(payload.User.Name),
		Email:    payload.User.Email,
		Username: payload.User.Email,
		Bio:      "",
	}

	if err := user.SetPassword(payload.User.Password); err != nil {
		return nil, err
	}

	if err := s.options.Database.NamedGetContext(ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio, created_at, updated_at)
		VALUES (:name, :email, :email, :password, :bio, NOW(), NOW()) RETURNING *
		`, user); err != nil {
		return nil, err
	}

	validationToken, err := data.NewValidationToken(user.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
	if err != nil {
		return nil, err
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.validation_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, validationToken); err != nil {
		return nil, err
	}

	if err := s.options.Emailer.SendValidationToken(user, validationToken); err != nil {
		return nil, err
	}

	log.Infow("validation", "token", validationToken.Token)

	s.options.Metrics.UserAdded()

	s.options.Metrics.EmailVerificationSent()

	pr := repositories.NewProjectRepository(s.options.Database)
	if err != nil {
		return nil, err
	}

	if _, err := pr.AddDefaultProject(ctx, user); err != nil {
		return nil, err
	}

	return UserType(s.options.signer, user)
}

func (s *UserService) Update(ctx context.Context, payload *user.UpdatePayload) (*user.User, error) {
	user := &data.User{
		ID:    payload.UserID,
		Name:  data.Name(payload.Update.Name),
		Email: payload.Update.Email,
		Bio:   payload.Update.Bio,
	}

	if err := s.options.Database.NamedGetContext(ctx, user, `
		UPDATE fieldkit.user SET name = :name, username = :email, email = :email, bio = :bio, updated_at = NOW() WHERE id = :id RETURNING *
		`, user); err != nil {
		return nil, err
	}

	return UserType(s.options.signer, user)
}

func (s *UserService) ChangePassword(ctx context.Context, payload *user.ChangePasswordPayload) (*user.User, error) {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	log := Logger(ctx).Sugar()

	log.Infow("password", "authorized_user_id", p.UserID(), "user_id", payload.UserID)

	updating := &data.User{}
	err = s.options.Database.GetContext(ctx, updating, `SELECT * FROM fieldkit.user WHERE id = $1`, p.UserID())
	if err == sql.ErrNoRows {
		return nil, user.MakeBadRequest(errors.New("bad request"))
	}
	if err != nil {
		return nil, err
	}

	if updating.ID != payload.UserID {
		return nil, user.MakeForbidden(errors.New("forbidden"))
	}

	err = updating.CheckPassword(payload.Change.OldPassword)
	if err == data.IncorrectPasswordError {
		return nil, user.MakeBadRequest(errors.New("bad request"))
	}
	if err != nil {
		return nil, err
	}

	if err := updating.SetPassword(payload.Change.NewPassword); err != nil {
		return nil, err
	}

	if err := s.options.Database.NamedGetContext(ctx, updating, `
		UPDATE fieldkit.user SET password = :password WHERE id = :id RETURNING *
		`, updating); err != nil {
		return nil, err
	}

	return UserType(s.options.signer, updating)
}

func (s *UserService) GetCurrent(ctx context.Context, payload *user.GetCurrentPayload) (*user.User, error) {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	currentUser := &data.User{}
	if err := s.options.Database.GetContext(ctx, currentUser, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, p.UserID()); err != nil {
		log := Logger(ctx).Sugar()
		log.Warnw("missing", "user_id", p.UserID())
		return nil, err
	}

	return UserType(s.options.signer, currentUser)
}

func (s *UserService) ListByProject(ctx context.Context, payload *user.ListByProjectPayload) (*user.ProjectUsers, error) {
	users := []*data.ProjectUserAndUser{}
	if err := s.options.Database.SelectContext(ctx, &users, `
		SELECT pu.*, u.* FROM fieldkit.user AS u JOIN fieldkit.project_user AS pu ON pu.user_id = u.id WHERE pu.project_id = $1 ORDER BY pu.role DESC, u.id
		`, payload.ProjectID); err != nil {
		return nil, err
	}

	invites := []*data.ProjectInvite{}
	if err := s.options.Database.SelectContext(ctx, &invites, `
		SELECT * FROM fieldkit.project_invite WHERE project_id = $1 AND accepted_time IS NULL ORDER BY invited_time
		`, payload.ProjectID); err != nil {
		return nil, err
	}

	return ProjectUsersType(s.options.signer, users, invites)
}

func (s *UserService) ProjectRoles(ctx context.Context) (user.ProjectRoleCollection, error) {
	roles := make([]*user.ProjectRole, 0)
	for _, role := range data.Roles {
		roles = append(roles, &user.ProjectRole{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return roles, nil
}

func (s *UserService) Recovery(ctx context.Context, payload *user.RecoveryPayload) error {
	log := Logger(ctx).Sugar()

	token := data.Token{}
	if err := token.UnmarshalText([]byte(payload.Recovery.Token)); err != nil {
		return err
	}

	log.Infow("recovery", "token_raw", payload.Recovery.Token, "token", token)

	recoveryToken := &data.RecoveryToken{}
	err := s.options.Database.GetContext(ctx, recoveryToken, `SELECT * FROM fieldkit.recovery_token WHERE token = $1`, token)
	if err == sql.ErrNoRows {
		log.Infow("recovery, token bad")
		return user.MakeUnauthorized(errors.New("unauthorized"))
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.After(recoveryToken.Expires) {
		log.Infow("recovery, token expired", "token_expires", recoveryToken.Expires, "now", now)
		return user.MakeUnauthorized(errors.New("unauthorized"))
	}

	trying := &data.User{}
	err = s.options.Database.GetContext(ctx, trying, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, recoveryToken.UserID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	if err := trying.SetPassword(payload.Recovery.Password); err != nil {
		return err
	}

	if _, err := s.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.recovery_token WHERE user_id = $1
		`, trying.ID); err != nil {
		return err
	}

	if err := s.options.Database.NamedGetContext(ctx, trying, `
		UPDATE fieldkit.user SET password = :password, valid = true WHERE id = :id RETURNING *
		`, trying); err != nil {
		return err
	}

	return nil
}

func (s *UserService) RecoveryLookup(ctx context.Context, payload *user.RecoveryLookupPayload) error {
	log := Logger(ctx).Sugar()

	trying := &data.User{}
	err := s.options.Database.GetContext(ctx, trying, `SELECT * FROM fieldkit.user WHERE LOWER(email) = LOWER($1)`, payload.Recovery.Email)
	if err == sql.ErrNoRows {
		log.Infow("recovery, no user")
		return nil
	}
	if err != nil {
		log.Errorw("recovery", "error", err)
		return nil
	}

	now := time.Now().UTC()

	recoveryToken, err := data.NewRecoveryToken(trying.ID, 20, now.Add(time.Duration(1)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := s.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.recovery_token WHERE user_id = $1`, trying.ID); err != nil {
		return err
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.recovery_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, recoveryToken); err != nil {
		return err
	}

	if err := s.options.Emailer.SendRecoveryToken(trying, recoveryToken); err != nil {
		return err
	}

	log.Infow("recovery", "token", recoveryToken.Token)

	s.options.Metrics.EmailRecoverPasswordSent()

	return nil
}

func (s *UserService) Refresh(ctx context.Context, payload *user.RefreshPayload) (*user.RefreshResult, error) {
	s.options.Metrics.AuthRefreshTry()

	token := data.Token{}
	if err := token.UnmarshalText([]byte(payload.RefreshToken)); err != nil {
		return nil, err
	}

	refreshToken := &data.RefreshToken{}
	err := s.options.Database.GetContext(ctx, refreshToken, `SELECT * FROM fieldkit.refresh_token WHERE token = $1`, token)
	if err == sql.ErrNoRows {
		return nil, user.MakeUnauthorized(errors.New("unauthorized"))
	}
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if now.After(refreshToken.Expires) {
		return nil, user.MakeUnauthorized(errors.New("unauthorized"))
	}

	trying := &data.User{}
	err = s.options.Database.GetContext(ctx, trying, `SELECT * FROM fieldkit.user WHERE id = $1`, refreshToken.UserID)
	if err == sql.ErrNoRows {
		return nil, user.MakeUnauthorized(errors.New("unauthorized"))
	}
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := data.NewRefreshToken(trying.ID, 20, now.Add(time.Duration(72)*time.Hour))
	if err != nil {
		return nil, err
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		DELETE FROM fieldkit.refresh_token WHERE token = :token
	`, refreshToken); err != nil {
		return nil, err
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.refresh_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, newRefreshToken); err != nil {
		return nil, err
	}

	jwtToken := trying.NewToken(now, newRefreshToken)
	signedToken, err := jwtToken.SignedString(s.options.JWTHMACKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err)
	}

	s.options.Metrics.AuthRefreshSuccess()

	return &user.RefreshResult{
		Authorization: "Bearer " + signedToken,
	}, nil
}

func (s *UserService) SendValidation(ctx context.Context, payload *user.SendValidationPayload) error {
	log := Logger(ctx).Sugar()

	updating := &data.User{}
	if err := s.options.Database.GetContext(ctx, updating, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, payload.UserID); err != nil {
		return err
	}

	// TODO Rate limit the number of these we can send?
	if !updating.Valid {
		validationToken, err := data.NewValidationToken(updating.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
		if err != nil {
			return err
		}

		if _, err := s.options.Database.ExecContext(ctx, `
			DELETE FROM fieldkit.validation_token WHERE user_id = $1
			`, updating.ID); err != nil {
			return err
		}

		if _, err := s.options.Database.NamedExecContext(ctx, `
			INSERT INTO fieldkit.validation_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
			`, validationToken); err != nil {
			return err
		}

		if err := s.options.Emailer.SendValidationToken(updating, validationToken); err != nil {
			return err
		}

		log.Infow("validation", "token", validationToken.Token)

		s.options.Metrics.EmailVerificationSent()
	}

	return nil
}

func (s *UserService) Validate(ctx context.Context, payload *user.ValidatePayload) (*user.ValidateResult, error) {
	log := Logger(ctx).Sugar()

	validationToken := &data.ValidationToken{}
	if err := validationToken.Token.UnmarshalText([]byte(payload.Token)); err != nil {
		return nil, err
	}

	err := s.options.Database.GetContext(ctx, validationToken, `
		SELECT * FROM fieldkit.validation_token WHERE token = $1
		`, validationToken.Token)
	if err == sql.ErrNoRows {
		log.Infow("invalid", "token", payload.Token)
		return &user.ValidateResult{
			Location: "https://" + s.options.PortalDomain + "?bad_token=true",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	if _, err := s.options.Database.ExecContext(ctx, `UPDATE fieldkit.user SET valid = true WHERE id = $1`, validationToken.UserID); err != nil {
		return nil, err
	}

	if _, err := s.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.validation_token WHERE token = $1`, validationToken.Token); err != nil {
		return nil, err
	}

	log.Infow("verified", "token", payload.Token)

	s.options.Metrics.UserValidated()

	return &user.ValidateResult{
		Location: "https://" + s.options.PortalDomain + "/",
	}, nil
}

func NewTransmissionToken(now time.Time, user *data.User) *jwtgo.Token {
	scopes := []string{"api:transmission"}

	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"iat":    now.Unix(),
		"exp":    now.Add(time.Hour * 24 * 365).Unix(),
		"sub":    user.ID,
		"email":  user.Email,
		"scopes": scopes,
	}

	return token
}

func (s *UserService) IssueTransmissionToken(ctx context.Context, payload *user.IssueTransmissionTokenPayload) (*user.TransmissionToken, error) {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	authed := &data.User{}
	if err := s.options.Database.GetContext(ctx, authed, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, p.UserID()); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	transmissionToken := NewTransmissionToken(now, authed)
	signedToken, err := transmissionToken.SignedString(s.options.JWTHMACKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err)
	}

	url := fmt.Sprintf("https://api.%s/ingestion", s.options.Domain)

	return &user.TransmissionToken{
		URL:   url,
		Token: signedToken,
	}, nil
}

func (s *UserService) AdminDelete(ctx context.Context, payload *user.AdminDeletePayload) error {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return err
	}

	admin := &data.User{}
	if err := s.options.Database.GetContext(ctx, admin, `SELECT * FROM fieldkit.user WHERE id = $1`, p.UserID()); err != nil {
		return user.MakeForbidden(errors.New("forbidden"))
	}

	err = admin.CheckPassword(payload.Delete.Password)
	if err != nil {
		return user.MakeForbidden(errors.New("forbidden"))
	}

	deleting := &data.User{}
	if err = s.options.Database.GetContext(ctx, deleting, `SELECT * FROM fieldkit.user WHERE LOWER(email) = LOWER($1)`, payload.Delete.Email); err != nil {
		return user.MakeForbidden(errors.New("forbidden"))
	}

	log.Infow("deleting", "user_id", deleting.ID)

	queries := []string{
		`DELETE FROM fieldkit.project_invite WHERE user_id = $1`,
		`DELETE FROM fieldkit.project_follower WHERE follower_id = $1`,
		`DELETE FROM fieldkit.project_user WHERE user_id = $1`,

		`DELETE FROM fieldkit.notes_media WHERE user_id = $1`,
		`DELETE FROM fieldkit.notes WHERE author_id = $1`,
		`DELETE FROM fieldkit.notes WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.notes_media WHERE id IN (SELECT media_id FROM fieldkit.notes_media_link WHERE note_id IN (SELECT id FROM fieldkit.notes WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)))`,
		`DELETE FROM fieldkit.notes_media_link WHERE note_id IN (SELECT id FROM fieldkit.notes WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1))`,
		`DELETE FROM fieldkit.notes WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_24h WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_12h WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_6h WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_1h WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_30m WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_10m WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.aggregated_1m WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.visible_configuration WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.project_station WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.station_activity WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`,
		`DELETE FROM fieldkit.station_ingestion WHERE uploader_id = $1`,
		`DELETE FROM fieldkit.station WHERE owner_id = $1`,
		`DELETE FROM fieldkit.user WHERE id = $1`,
	}

	for _, query := range queries {
		if _, err := s.options.Database.ExecContext(ctx, query, deleting.ID); err != nil {
			return err
		}
	}

	return nil
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
	saved, err := mr.Save(ctx, body, payload.ContentLength, payload.ContentType)
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

func (s *UserService) DownloadPhoto(ctx context.Context, payload *user.DownloadPhotoPayload) (*user.DownloadedPhoto, error) {
	resource := &data.User{}
	if err := s.options.Database.GetContext(ctx, resource, `SELECT * FROM fieldkit.user WHERE id = $1`, payload.UserID); err != nil {
		return nil, err
	}

	if resource.MediaURL == nil || resource.MediaContentType == nil {
		return nil, user.MakeNotFound(errors.New("not found"))
	}

	etag := quickHash(*resource.MediaURL) + ""
	if payload.Size != nil {
		etag += fmt.Sprintf(":%d", *payload.Size)
	}

	if payload.IfNoneMatch != nil {
		if *payload.IfNoneMatch == fmt.Sprintf(`"%s"`, etag) {
			return &user.DownloadedPhoto{
				Etag: etag,
				Body: []byte{},
			}, nil
		}
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)
	lm, err := mr.LoadByURL(ctx, *resource.MediaURL)
	if err != nil {
		return nil, user.MakeNotFound(errors.New("not found"))
	}

	data, err := ioutil.ReadAll(lm.Reader)
	if err != nil {
		return nil, err
	}

	return &user.DownloadedPhoto{
		Length:      int64(lm.Size),
		ContentType: *resource.MediaContentType,
		Etag:        etag,
		Body:        data,
	}, nil
}

func (s *UserService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return user.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return user.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return user.MakeForbidden(errors.New(m)) },
	})
}

func UserType(signer *Signer, dm *data.User) (*user.User, error) {
	userType := &user.User{
		ID:    dm.ID,
		Name:  dm.Name,
		Email: dm.Email,
		Bio:   dm.Bio,
		Admin: dm.Admin,
	}

	if dm.MediaURL != nil {
		url := fmt.Sprintf("/user/%d/media", dm.ID)
		userType.Photo = &user.UserPhoto{
			URL: &url,
		}
	}

	return userType, nil
}

func ProjectUserType(signer *Signer, dm *data.ProjectUserAndUser) (*user.ProjectUser, error) {
	userWm, err := UserType(signer, &dm.User)
	if err != nil {
		return nil, err
	}
	return &user.ProjectUser{
		User:       userWm,
		Role:       dm.RoleName(),
		Membership: data.MembershipAccepted,
	}, nil
}

func ProjectUsersType(signer *Signer, users []*data.ProjectUserAndUser, invites []*data.ProjectInvite) (*user.ProjectUsers, error) {
	usersCollection := make([]*user.ProjectUser, 0, len(users)+len(invites))
	for _, dm := range users {
		pu, err := ProjectUserType(signer, dm)
		if err != nil {
			return nil, err
		}
		usersCollection = append(usersCollection, pu)
	}
	for _, invite := range invites {
		membership := data.MembershipPending
		if invite.AcceptedTime != nil {
			membership = data.MembershipAccepted
		}
		if invite.RejectedTime != nil {
			membership = data.MembershipRejected
		}

		if invite.RejectedTime == nil && invite.AcceptedTime == nil {
			usersCollection = append(usersCollection, &user.ProjectUser{
				User: &user.User{
					Name:  invite.InvitedEmail,
					Email: invite.InvitedEmail,
				},
				Role:       invite.LookupRole().Name,
				Membership: membership,
				Accepted:   invite.AcceptedTime != nil,
				Rejected:   invite.RejectedTime != nil,
				Invited:    true,
			})
		}
	}

	return &user.ProjectUsers{
		Users: usersCollection,
	}, nil
}

func (s *UserService) userExists(ctx context.Context, email string) (bool, error) {
	existing := []*data.User{}
	err := s.options.Database.SelectContext(ctx, &existing, `SELECT * FROM fieldkit.user WHERE LOWER(email) = LOWER($1)`, email)
	if err != nil {
		return false, err
	}
	if len(existing) > 0 {
		return true, nil
	}
	return false, nil
}

func (s *UserService) authenticateOrSpoof(ctx context.Context, email, password string) (*data.User, error) {
	user := &data.User{}
	err := s.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE LOWER(u.email) = LOWER($1)`, email)
	if err == sql.ErrNoRows {
		return nil, data.IncorrectPasswordError
	}
	if err != nil {
		return nil, err
	}

	parts := strings.Split(password, " ")
	if len(parts) == 2 {
		log := Logger(ctx).Sugar()

		// NOTE Not logging password here, may be a real one.
		log.Infow("spoofing", "email", email)

		adminUser := &data.User{}
		err := s.options.Database.GetContext(ctx, adminUser, `SELECT u.* FROM fieldkit.user AS u WHERE LOWER(u.email) = LOWER($1) AND u.admin`, parts[0])
		if err == nil {
			// We can safely log the user doing the spoofing here.
			err = adminUser.CheckPassword(parts[1])
			if err == nil {
				log.Infow("spoofed", "email", email, "spoofer", adminUser.Email)
				return user, nil
			} else {
				log.Infow("denied", "email", email, "spoofer", adminUser.Email, "spoofing", true)
			}
		}
	}

	err = user.CheckPassword(password)
	if err == data.IncorrectPasswordError {
		return nil, data.IncorrectPasswordError
	}
	if err != nil {
		return nil, err
	}

	if !user.Valid {
		return nil, data.IncorrectPasswordError
	}

	return user, nil
}
